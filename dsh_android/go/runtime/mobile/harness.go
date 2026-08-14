package mobile

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/deepseek/dsh-android/go-runtime/agent"
	"github.com/deepseek/dsh-android/go-runtime/config"
	"github.com/deepseek/dsh-android/go-runtime/session"
	"github.com/deepseek/dsh-android/go-runtime/tools"
)

type Harness struct {
	mu           sync.Mutex
	store        *session.SQLiteStore
	configStore  *config.Store
	registry     *tools.Registry
	currentLoop  *agent.Loop
	currentID    session.SessionID
	eventListeners map[string]chan *agent.AgentEvent
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
		
		// Register todo tool
		reg.Register(&tools.TodoTool{})
	
	return &Harness{
		store:          store,
		configStore:    cfgStore,
		registry:       reg,
		eventListeners: make(map[string]chan *agent.AgentEvent),
	}, nil
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
		MaxTurns:      cfg.Agent.MaxTurns,
		ToolTimeoutMs: cfg.Agent.ToolTimeoutMs,
		Model:         cfg.Agent.Model,
		APIKey:        cfg.Agent.APIKey,
		BaseURL:       cfg.Agent.BaseURL,
		Temperature:   cfg.Agent.Temperature,
		SystemPrompt:  cfg.Agent.SystemPrompt,
	}
	
	loop := agent.NewLoop(session.SessionID(sessionID), h.registry, h.store, loopCfg)
	h.currentLoop = loop
	h.mu.Unlock()
	
	evChan := make(chan *agent.AgentEvent, 64)
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
	return h.store.Close()
}
