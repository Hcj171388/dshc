package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/deepseek/dsh-android/go-runtime/llm"
	"github.com/deepseek/dsh-android/go-runtime/session"
	"github.com/deepseek/dsh-android/go-runtime/tools"
)

// Config holds agent loop configuration
type Config struct {
	MaxTurns       int
	ToolTimeoutMs  int
	Model          string
	APIKey         string
	BaseURL        string
	Temperature    float64
	SystemPrompt   string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		MaxTurns:      20,
		ToolTimeoutMs: 30000,
		Temperature:   0.7,
		SystemPrompt:  "You are a helpful assistant that can execute commands and manipulate files.",
	}
}

// SetConfig sets the global default config
var globalConfig = DefaultConfig()

func SetConfig(cfg Config) {
	globalConfig = cfg
}

// Loop is the agent loop
type Loop struct {
	sessionID  session.SessionID
	registry   *tools.Registry
	store      session.Store
	config     Config
	abort      context.CancelFunc
	mu         sync.Mutex
	llmClient  *llm.Client
}

// NewLoop creates a new agent loop
func NewLoop(sid session.SessionID, reg *tools.Registry, store session.Store, cfg Config) *Loop {
	if cfg.MaxTurns == 0 {
		cfg = globalConfig
	}
	
	llmCfg := llm.Config{
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
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
	evChan := make(chan *AgentEvent, 128)
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

		// Build messages history
		messages := l.buildMessages(ctx)
		messages = append(messages, llm.Message{Role: "user", Content: prompt})

		for turn := 0; turn < l.config.MaxTurns; turn++ {
			if ctx.Err() != nil {
				break
			}

			response, err := l.processTurn(ctx, messages)
			if err != nil {
				l.sendEvent(evChan, "error", map[string]string{"message": err.Error()})
				break
			}

			// Send response event
			l.sendEvent(evChan, "response", map[string]string{"text": response})
			
			// Store assistant response
			l.store.AddEvent(&session.Event{
				SessionID: string(l.sessionID),
				Type:      "assistant_message",
				Payload:   response,
			})

			// Check if we need another turn (tool calls)
			if !responseHasToolCalls(response) {
				break
			}
		}

		l.sendEvent(evChan, "turn_complete", nil)
	}()
	return evChan
}

func (l *Loop) buildMessages(ctx context.Context) []llm.Message {
	messages := []llm.Message{
		{Role: "system", Content: l.config.SystemPrompt},
	}
	
	// Load previous events
	events, err := l.store.GetEvents(string(l.sessionID), 0, 100)
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
			var tc ToolCallEvent
			json.Unmarshal([]byte(ev.Payload), &tc)
			messages = append(messages, llm.Message{Role: "assistant", Content: map[string]interface{}{
				"tool_calls": []map[string]interface{}{
					{"id": tc.ID, "type": "function", "function": map[string]interface{}{"name": tc.Name, "arguments": tc.Arguments}},
				},
			}})
			messages = append(messages, llm.Message{Role: "tool", Content: map[string]interface{}{
				"tool_call_id": tc.ID,
				"content":      tc.Result,
			}})
		}
	}
	
	return messages
}

func (l *Loop) processTurn(ctx context.Context, messages []llm.Message) (string, error) {
	// Get tools from registry
	tools := l.registry.ToLLMTools()
	
	// Call LLM
	resp, err := l.llmClient.Call(ctx, llm.CallOptions{
		Messages:    messages,
		Tools:       toLLMTools(tools),
		Temperature: l.config.Temperature,
	})
	if err != nil {
		return "", fmt.Errorf("llm call: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}
	
	assistantMsg := resp.Choices[0].Message
	
	// Handle tool calls
	if toolCalls, ok := assistantMsg.Content.(map[string]interface{})["tool_calls"].([]interface{}); ok && len(toolCalls) > 0 {
		var texts []string
		
		for _, tc := range toolCalls {
			call := tc.(map[string]interface{})
			funcCall := call["function"].(map[string]interface{})
			name := funcCall["name"].(string)
			args := json.RawMessage(funcCall["arguments"].(string))
			callID := call["id"].(string)
			
			// Send tool call event
			l.sendEvent(nil, "tool_call", map[string]string{
				"id":      callID,
				"name":    name,
				"args":    string(args),
			})
			
			// Execute tool
			result, err := l.registry.Call(name, args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}
			
			// Store tool call event
			l.store.AddEvent(&session.Event{
				SessionID: string(l.sessionID),
				Type:      "tool_call",
				Payload: ToolCallEvent{
					ID:      callID,
					Name:    name,
					Arguments: string(args),
					Result:  result,
				}.String(),
			})
			
			// Add tool result to messages
			messages = append(messages, llm.Message{Role: "assistant", Content: map[string]interface{}{
				"tool_calls": []map[string]interface{}{
					{"id": callID, "type": "function", "function": map[string]interface{}{"name": name, "arguments": string(args)}},
				},
			}})
			messages = append(messages, llm.Message{Role: "tool", Content: map[string]interface{}{
				"tool_call_id": callID,
				"content":      result,
			}})
			
			texts = append(texts, result)
		}
		
		// Continue with tool results
		return l.processTurn(ctx, messages)
	}
	
	// Return text response
	text, _ := assistantMsg.Content.(string)
	return text, nil
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

// ToolCallEvent represents a tool call event
type ToolCallEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	Result    string `json:"result"`
}

func (t ToolCallEvent) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func responseHasToolCalls(response string) bool {
	// This is a simplified check - in reality we'd parse the JSON response
	return false
}

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
