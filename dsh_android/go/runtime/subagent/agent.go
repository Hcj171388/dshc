package subagent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/deepseek/dsh-android/go-runtime/agent"
	"github.com/deepseek/dsh-android/go-runtime/session"
	"github.com/deepseek/dsh-android/go-runtime/tools"
)

// RunResult represents the result of a subagent run
type RunResult struct {
	ID      string `json:"id"`
	Output  string `json:"output"`
	Turns   int    `json:"turns"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Manager manages subagent instances
type Manager struct {
	mu       sync.Mutex
	agents   map[string]*subagentInstance
	store    session.Store
	registry *tools.Registry
	config   agent.Config
	nextID   int
}

type subagentInstance struct {
	id        string
	loop      *agent.Loop
	startTime time.Time
	endTime   time.Time
	result    *RunResult
}

// NewManager creates a new subagent manager
func NewManager(store session.Store, reg *tools.Registry, cfg agent.Config) *Manager {
	return &Manager{
		agents:   make(map[string]*subagentInstance),
		store:    store,
		registry: reg,
		config:   cfg,
	}
}

// Start launches a new subagent
func (m *Manager) Start(parentSessionID string, prompt string) (*RunResult, error) {
	m.mu.Lock()
	m.nextID++
	id := fmt.Sprintf("sub_%d_%d", time.Now().UnixNano(), m.nextID)
	m.mu.Unlock()

	loop := agent.NewLoop(session.SessionID(id), m.registry, m.store, m.config)
	
	resultCh := make(chan *RunResult, 1)
	go func() {
		events := loop.Run(prompt)
		output := ""
		turns := 0
		for ev := range events {
			turns++
			if ev.Type == "message" && ev.Payload != nil {
				var msg map[string]string
				if json.Unmarshal(ev.Payload, &msg) == nil {
					if msg["role"] == "assistant" {
						output = msg["content"]
					}
				}
			}
		}
		resultCh <- &RunResult{
			ID:      id,
			Output:  output,
			Turns:   turns,
			Success: true,
		}
	}()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	select {
	case result := <-resultCh:
		m.mu.Lock()
		m.agents[id] = &subagentInstance{
			id:        id,
			loop:      loop,
			startTime: time.Now(),
			endTime:   time.Now(),
			result:    result,
		}
		m.mu.Unlock()
		return result, nil
	case <-ctx.Done():
		return &RunResult{
			ID:      id,
			Success: false,
			Error:   "timeout",
		}, nil
	}
}

// GetResult returns the result of a subagent
func (m *Manager) GetResult(id string) (*RunResult, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	inst, ok := m.agents[id]
	if !ok {
		return nil, false
	}
	return inst.result, true
}

// List returns all subagent IDs
func (m *Manager) List() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := make([]string, 0, len(m.agents))
	for id := range m.agents {
		ids = append(ids, id)
	}
	return ids
}

// Abort cancels a running subagent
func (m *Manager) Abort(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if inst, ok := m.agents[id]; ok {
		inst.loop.Abort()
	}
}
