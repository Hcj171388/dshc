package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/deepseek/dsh-android/go-runtime/llm"
	"github.com/deepseek/dsh-android/go-runtime/session"
	"github.com/deepseek/dsh-android/go-runtime/tools"
)

// Config holds agent loop configuration
type Config struct {
	MaxTurns       int
	MaxSteps       int
	ToolTimeoutMs  int
	Model          string
	APIKey         string
	BaseURL        string
	Temperature    float64
	MaxTokens      int
	SystemPrompt   string
	CompactionMode string // "none", "summary", "truncate"
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		MaxTurns:       30,
		MaxSteps:       50,
		ToolTimeoutMs:  30000,
		Temperature:    0.7,
		MaxTokens:      4096,
		SystemPrompt:   "You are DeepSeek Harness, a helpful assistant that can execute commands and manipulate files.",
		CompactionMode: "none",
	}
}

// SetConfig sets the global default config
var globalConfig = DefaultConfig()

func SetConfig(cfg Config) {
	globalConfig = cfg
}

// Loop is the agent loop
type Loop struct {
	sessionID session.SessionID
	registry  *tools.Registry
	store     session.Store
	config    Config
	abort     context.CancelFunc
	mu        sync.Mutex
	llmClient *llm.Client
}

// NewLoop creates a new agent loop
func NewLoop(sid session.SessionID, reg *tools.Registry, store session.Store, cfg Config) *Loop {
	if cfg.MaxTurns == 0 {
		cfg = globalConfig
	}
	
	llmCfg := llm.Config{
		APIKey:    cfg.APIKey,
		BaseURL:   cfg.BaseURL,
		Model:     cfg.Model,
		MaxTokens: cfg.MaxTokens,
	}
	
	return &Loop{
		sessionID: sid,
		registry:  reg,
		store:     store,
		config:    cfg,
		llmClient: llm.NewClient(llmCfg),
	}
}

// Abort cancels the current loop execution
func (l *Loop) Abort() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.abort != nil {
		l.abort()
	}
}

// Run starts the agent loop with the given prompt
func (l *Loop) Run(prompt string) <-chan *AgentEvent {
	evChan := make(chan *AgentEvent, 256)
	ctx, cancel := context.WithCancel(context.Background())
	l.mu.Lock()
	l.abort = cancel
	l.mu.Unlock()

	go func() {
		defer close(evChan)
		defer cancel()

		l.sendEvent(evChan, "turn_start", map[string]string{"prompt": prompt})
		
		// Store user message
		l.store.AddEvent(&session.Event{
			SessionID: string(l.sessionID),
			Type:      "user_message",
			Payload:   prompt,
		})
		l.sendEvent(evChan, "message", map[string]string{
			"role":    "user",
			"content": prompt,
		})

		// Build initial messages
		messages := l.buildMessages(ctx)
		
		// Add system prompt
		if l.config.SystemPrompt != "" {
			messages = append([]llm.Message{{Role: "system", Content: l.config.SystemPrompt}}, messages...)
		}

		// Add user message
		messages = append(messages, llm.Message{Role: "user", Content: prompt})

		turnCount := 0
		stepCount := 0
		
		for turnCount < l.config.MaxTurns {
			if ctx.Err() != nil {
				break
			}
			
			turnCount++
			l.sendEvent(evChan, "turn", map[string]int{"number": turnCount})
			
			// Process this turn
			response, hasToolCalls, err := l.processTurn(ctx, messages)
			if err != nil {
				l.sendEvent(evChan, "error", map[string]string{"message": err.Error()})
				break
			}
			
			// Store assistant response
			l.store.AddEvent(&session.Event{
				SessionID: string(l.sessionID),
				Type:      "assistant_message",
				Payload:   response,
			})
			l.sendEvent(evChan, "message", map[string]string{
				"role":    "assistant",
				"content": response,
			})
			
			if !hasToolCalls {
				break
			}
			
			// Check max steps
			stepCount++
			if stepCount >= l.config.MaxSteps {
				l.sendEvent(evChan, "error", map[string]string{"message": "max steps reached"})
				break
			}
		}
		
		l.sendEvent(evChan, "turn_complete", map[string]int{"turns": turnCount, "steps": stepCount})
	}()
	return evChan
}

func (l *Loop) buildMessages(ctx context.Context) []llm.Message {
	messages := []llm.Message{}
	
	// Load previous events
	events, err := l.store.GetEvents(string(l.sessionID), 0, 200)
	if err != nil {
		return messages
	}
	
	for _, ev := range events {
		switch ev.Type {
		case "user_message":
			messages = append(messages, llm.Message{Role: "user", Content: ev.Payload})
		case "assistant_message":
			messages = append(messages, llm.Message{Role: "assistant", Content: ev.Payload})
		case "tool_call":
			var tc ToolCallData
			if json.Unmarshal([]byte(ev.Payload), &tc) == nil {
				// Add assistant message with tool calls
				messages = append(messages, llm.Message{Role: "assistant", Content: map[string]interface{}{
					"tool_calls": []map[string]interface{}{
						{
							"id": tc.ID,
							"type": "function",
							"function": map[string]interface{}{
								"name":      tc.Name,
								"arguments": tc.Arguments,
							},
						},
					},
				}})
				// Add tool result
				messages = append(messages, llm.Message{Role: "tool", Content: map[string]interface{}{
					"tool_call_id": tc.ID,
					"content":      tc.Result,
				}})
			}
		}
	}
	
	return messages
}

func (l *Loop) processTurn(ctx context.Context, messages []llm.Message) (string, bool, error) {
	// Get tools from registry
	llmTools := l.registry.ToLLMTools()
	
	// Call LLM
	resp, err := l.llmClient.Call(ctx, llm.CallOptions{
		Messages:    messages,
		Tools:       toLLMTools(llmTools),
		Temperature: l.config.Temperature,
		MaxTokens:   l.config.MaxTokens,
	})
	if err != nil {
		return "", false, fmt.Errorf("llm call: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return "", false, fmt.Errorf("empty response")
	}
	
	assistantMsg := resp.Choices[0].Message
	
	// Check if response has tool calls
	content, ok := assistantMsg.Content.(map[string]interface{})
	if !ok {
		text, _ := assistantMsg.Content.(string)
		return text, false, nil
	}
	
	toolCalls, hasToolCalls := content["tool_calls"].([]interface{})
	if !hasToolCalls || len(toolCalls) == 0 {
		text, _ := content["text"].(string)
		return text, false, nil
	}
	
	// Execute tool calls
	var toolResults []string
	for _, tc := range toolCalls {
		call := tc.(map[string]interface{})
		funcCall := call["function"].(map[string]interface{})
		name := funcCall["name"].(string)
		argsStr := funcCall["arguments"].(string)
		callID := call["id"].(string)
		
		args := json.RawMessage(argsStr)
		
		// Send tool call event
		l.sendEvent(nil, "tool_call", map[string]string{
			"id":      callID,
			"name":    name,
			"args":    argsStr,
		})
		
		// Execute tool
		result, err := l.registry.Call(name, args)
		if err != nil {
			result = fmt.Sprintf("Error: %v", err)
		}
		
		toolResults = append(toolResults, result)
		
		// Store tool call event
		l.store.AddEvent(&session.Event{
			SessionID: string(l.sessionID),
			Type:      "tool_call",
			Payload: ToolCallData{
				ID:        callID,
				Name:      name,
				Arguments: argsStr,
				Result:    result,
			}.String(),
		})
		
		// Add tool result to messages
		messages = append(messages, llm.Message{Role: "assistant", Content: map[string]interface{}{
			"tool_calls": []map[string]interface{}{
				{
					"id": callID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      name,
						"arguments": argsStr,
					},
				},
			},
		}})
		messages = append(messages, llm.Message{Role: "tool", Content: map[string]interface{}{
			"tool_call_id": callID,
			"content":      result,
		}})
	}
	
	// Continue with tool results
	nextResponse, _, err := l.processTurn(ctx, messages)
	if err != nil {
		return "", false, err
	}
	
	return nextResponse, true, nil
}

func (l *Loop) sendEvent(ch chan<- *AgentEvent, typ string, payload interface{}) {
	ev := &AgentEvent{Type: typ}
	if payload != nil {
		data, _ := json.Marshal(payload)
		ev.Payload = data
	}
	if ch != nil {
		select {
		case ch <- ev:
		case <-time.After(5 * time.Second):
		}
	}
}

// ToolCallData represents tool call information
type ToolCallData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	Result    string `json:"result"`
}

func (t ToolCallData) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// HasToolCalls checks if a response contains tool calls
func HasToolCalls(content interface{}) bool {
	m, ok := content.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasTools := m["tool_calls"]
	return hasTools
}

// IsTextResponse checks if the response is pure text
func IsTextResponse(content interface{}) bool {
	_, isString := content.(string)
	return isString
}

// FormatToolCalls formats tool calls for display
func FormatToolCalls(content interface{}) string {
	m, ok := content.(map[string]interface{})
	if !ok {
		return ""
	}
	toolCalls, ok := m["tool_calls"].([]interface{})
	if !ok {
		return ""
	}
	var parts []string
	for _, tc := range toolCalls {
		call := tc.(map[string]interface{})
		funcCall := call["function"].(map[string]interface{})
		parts = append(parts, fmt.Sprintf("%s(%s)", funcCall["name"], funcCall["arguments"]))
	}
	return strings.Join(parts, ", ")
}

// toLLMTools converts tools.LLMTool to llm.Tool
func toLLMTools(tools []tools.LLMTool) []llm.Tool {
	result := make([]llm.Tool, len(tools))
	for i, t := range tools {
		result[i] = llm.Tool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
		}
	}
	return result
}
