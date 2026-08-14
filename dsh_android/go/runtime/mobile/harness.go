package mobile

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/deepseek/dsh-android/go-runtime/agent"
	"github.com/deepseek/dsh-android/go-runtime/compaction"
	"github.com/deepseek/dsh-android/go-runtime/config"
	"github.com/deepseek/dsh-android/go-runtime/credentials"
	"github.com/deepseek/dsh-android/go-runtime/guard"
	"github.com/deepseek/dsh-android/go-runtime/interaction"
	"github.com/deepseek/dsh-android/go-runtime/mcp"
	"github.com/deepseek/dsh-android/go-runtime/session"
	"github.com/deepseek/dsh-android/go-runtime/subagent"
	"github.com/deepseek/dsh-android/go-runtime/tools"
	"github.com/deepseek/dsh-android/go-runtime/workflow"
)

type Harness struct {
	mu            sync.Mutex
	store         *session.SQLiteStore
	configStore   *config.Store
	registry      *tools.Registry
	currentLoop   *agent.Loop
	currentID     session.SessionID
	eventListeners map[string]chan *agent.AgentEvent
	
	// New components
	mcpClients    map[string]*mcp.Client
	subagentMgr   *subagent.Manager
	workflowMgr   *workflow.Runner
	guardian      *guard.Guardian
	approver      *interaction.Approver
	compactor     *compaction.Compressor
	credStore     *credentials.Store
}

func NewHarness(dataDir string) (*Harness, error) {
	store, err := session.NewSessionStore(fmt.Sprintf("%s/sessions.db", dataDir))
	if err != nil {
		return nil, err
	}
	cfgStore, err := config.NewStore(dataDir)
	if err != nil {
		store.Close()
		return nil, err
	}
	
	reg := tools.NewRegistry()
	
	// Register bash tool
	bashTool := &tools.BashTool{DefaultTimeoutMs: cfgStore.Get().Tools.Bash.TimeoutMs}
	reg.Register(bashTool)
	
	// Register file system tools
	fsTools := tools.NewFsTools(
		cfgStore.Get().Tools.Fs.ReadLimit,
		cfgStore.Get().Tools.Fs.ReadMaxBytes,
		cfgStore.Get().Tools.Fs.ReadMaxLineLen,
	)
	fsTools.RegisterAll(reg)
	
	// Register web tools
	webTools := tools.NewWebTools(30 * time.Second)
	webTools.RegisterAll(reg)
	
	// Register terminal tool
	reg.Register(&tools.TerminalTool{})
		reg.Register(&tools.ACPTool{})
		reg.Register(&tools.LSPTool{})
		reg.Register(&tools.SkillTool{})
		reg.Register(&tools.GoalTool{})
		reg.Register(&tools.ContextTool{})
		reg.Register(&tools.E2BTool{})
		reg.Register(&tools.GitTool{})
		reg.Register(&tools.WebFetchTool{})
		reg.Register(&tools.CurlTool{})
		reg.Register(&tools.SkillTool{})
		reg.Register(&tools.LSPTool{})
		reg.Register(&tools.SkillTool{})
		reg.Register(&tools.GoalTool{})
		reg.Register(&tools.ContextTool{})
		reg.Register(&tools.E2BTool{})
		reg.Register(&tools.GitTool{})
		reg.Register(&tools.WebFetchTool{})
		reg.Register(&tools.CurlTool{})
		reg.Register(&tools.SkillTool{})
	
	// Register todo tool
	reg.Register(&tools.TodoTool{})
	
	// Initialize components
	g := guard.NewGuardian(guard.DefaultConfig())
	app := interaction.NewApprover()
	cmp := compaction.NewCompressor("summary", 100, 500)
	cred, err := credentials.NewStore(dataDir)
	if err != nil {
		return nil, err
	}
	
	h := &Harness{
		store:          store,
		configStore:    cfgStore,
		registry:       reg,
		eventListeners: make(map[string]chan *agent.AgentEvent),
		mcpClients:     make(map[string]*mcp.Client),
		subagentMgr:    subagent.NewManager(store, reg, agent.Config{}),
		workflowMgr:    workflow.NewRunner(),
		guardian:       g,
		approver:       app,
		compactor:      cmp,
		credStore:      cred,
	}
	
	return h, nil
}

func (h *Harness) CreateSession() (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	id, err := h.store.CreateSession()
	if err != nil {
		return "", err
	}
	h.currentID = id
	return string(id), nil
}

func (h *Harness) ListSessions() (string, error) {
	sessions, err := h.store.ListSessions()
	if err != nil {
		return "[]", err
	}
	data, _ := json.Marshal(sessions)
	return string(data), nil
}

func (h *Harness) GetSession(id string) (string, error) {
	meta, err := h.store.GetSession(id)
	if err != nil {
		return "null", err
	}
	data, _ := json.Marshal(meta)
	return string(data), nil
}

func (h *Harness) DeleteSession(id string) error {
	return h.store.DeleteSession(id)
}

func (h *Harness) ArchiveSession(id string) error {
	return h.store.ArchiveSession(id)
}

func (h *Harness) UpdateSessionTitle(id, title string) error {
	return h.store.UpdateSessionTitle(id, title)
}

func (h *Harness) GetEvents(sessionID string, afterID int64, limit int) (string, error) {
	events, err := h.store.GetEvents(sessionID, afterID, limit)
	if err != nil {
		return "[]", err
	}
	data, _ := json.Marshal(events)
	return string(data), nil
}

func (h *Harness) RunAgent(sessionID, prompt string) (string, error) {
	h.mu.Lock()
	if h.currentLoop != nil {
		h.currentLoop.Abort()
		h.currentLoop = nil
	}
	
	cfg := h.configStore.Get()
	loopCfg := agent.Config{
		MaxTurns:       cfg.Agent.MaxTurns,
		ToolTimeoutMs:  cfg.Agent.ToolTimeoutMs,
		Model:          cfg.Agent.Model,
		APIKey:         cfg.Agent.APIKey,
		BaseURL:        cfg.Agent.BaseURL,
		Temperature:    cfg.Agent.Temperature,
		SystemPrompt:   cfg.Agent.SystemPrompt,
		CompactionMode: cfg.Agent.CompactionMode,
	}
	
	loop := agent.NewLoop(session.SessionID(sessionID), h.registry, h.store, loopCfg)
	h.currentLoop = loop
	h.mu.Unlock()
	
	evChan := make(chan *agent.AgentEvent, 128)
	go func() {
		events := loop.Run(prompt)
		for ev := range events {
			select {
			case evChan <- ev:
			default:
			}
		}
	}()
	
	listenerID := fmt.Sprintf("listener_%d", time.Now().UnixNano())
	h.mu.Lock()
	h.eventListeners[listenerID] = evChan
	h.mu.Unlock()
	return listenerID, nil
}

func (h *Harness) ConsumeEvents(listenerID string, cb chan<- string) error {
	h.mu.Lock()
	ch, ok := h.eventListeners[listenerID]
	h.mu.Unlock()
	if !ok {
		return fmt.Errorf("listener not found: %s", listenerID)
	}
	defer func() {
		h.mu.Lock()
		delete(h.eventListeners, listenerID)
		h.mu.Unlock()
	}()
	for ev := range ch {
		data, _ := json.Marshal(ev)
		cb <- string(data)
	}
	return nil
}

func (h *Harness) AbortAgent() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.currentLoop != nil {
		h.currentLoop.Abort()
		h.currentLoop = nil
	}
	h.guardian.Abort()
}

func (h *Harness) GetConfig() (string, error) {
	cfg := h.configStore.Get()
	data, _ := json.Marshal(cfg)
	return string(data), nil
}

func (h *Harness) SaveConfig(raw string) error {
	var cfg config.AppConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return err
	}
	h.configStore.Get().Agent = cfg.Agent
	h.configStore.Get().Tools = cfg.Tools
	return h.configStore.Save()
}

func (h *Harness) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.currentLoop != nil {
		h.currentLoop.Abort()
		h.currentLoop = nil
	}
	for _, client := range h.mcpClients {
		client.Close()
	}
	return h.store.Close()
}

// MCP Methods
func (h *Harness) ConnectMCP(serverName, command string, args []string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, exists := h.mcpClients[serverName]; exists {
		return "", fmt.Errorf("MCP server already connected: %s", serverName)
	}
	
	client, err := mcp.NewClient(serverName, command, args)
	if err != nil {
		return "", err
	}
	
	h.mcpClients[serverName] = client
	return fmt.Sprintf("Connected to MCP server: %s", serverName), nil
}

func (h *Harness) DisconnectMCP(serverName string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if client, ok := h.mcpClients[serverName]; ok {
		client.Close()
		delete(h.mcpClients, serverName)
	}
	return nil
}

// Subagent Methods
func (h *Harness) StartSubagent(sessionID, prompt string) (string, error) {
	result, err := h.subagentMgr.Start(sessionID, prompt)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func (h *Harness) GetSubagentResult(id string) (string, error) {
	result, ok := h.subagentMgr.GetResult(id)
	if !ok {
		return "null", nil
	}
	data, _ := json.Marshal(result)
	return string(data), nil
}

func (h *Harness) AbortSubagent(id string) {
	h.subagentMgr.Abort(id)
}

// Workflow Methods
func (h *Harness) RunWorkflow(wfJSON string) (string, error) {
	var wf workflow.Workflow
	if err := json.Unmarshal([]byte(wfJSON), &wf); err != nil {
		return "", err
	}
	
	result, err := h.workflowMgr.Run(&wf)
	if err != nil {
		return "", err
	}
	
	data, _ := json.Marshal(result)
	return string(data), nil
}

func (h *Harness) StopWorkflow(id string) {
	h.workflowMgr.Stop(id)
}

// Interaction Methods
func (h *Harness) ApproveTool(toolName string, approve bool) {
	h.approver.ApproveTool(toolName, approve)
}

func (h *Harness) GetApprovalPresets() (string, error) {
	presets := h.approver.GetPresets()
	data, _ := json.Marshal(presets)
	return string(data), nil
}

// Credential Methods
func (h *Harness) SetCredential(ref, value string) error {
	return h.credStore.Set(ref, value)
}

func (h *Harness) GetCredential(ref string) (string, error) {
	val, ok := h.credStore.Resolve(ref)
	if !ok {
		return "", fmt.Errorf("credential not found: %s", ref)
	}
	return val, nil
}

// Compaction Methods
func (h *Harness) CompactSession(sessionID string) (string, error) {
	events, err := h.store.GetEvents(sessionID, 0, 1000)
	if err != nil {
		return "", err
	}
	
	evList := make([]compaction.Event, len(events))
	for i, ev := range events {
		evList[i] = compaction.Event{Type: ev.Type, Payload: ev.Payload}
	}
	
	result, err := h.compactor.Compact(evList)
	if err != nil {
		return "", err
	}
	
	data, _ := json.Marshal(result)
	return string(data), nil
}

// Guardian Methods
func (h *Harness) GetGuardianStats() (string, error) {
	stats := h.guardian.Stats()
	data, _ := json.Marshal(stats)
	return string(data), nil
}
