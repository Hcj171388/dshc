package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/deepseek/dsh-android/go-runtime/session"
	"github.com/deepseek/dsh-android/go-runtime/tools"
)

var globalConfig = DefaultConfig()

func SetConfig(cfg Config) {
	globalConfig = cfg
}

type Loop struct {
	sessionID session.SessionID
	registry  *tools.Registry
	store     session.Store
	config    Config
	abort     context.CancelFunc
	mu        sync.Mutex
}

func NewLoop(sid session.SessionID, reg *tools.Registry, store session.Store, cfg Config) *Loop {
	if cfg.MaxTurns == 0 {
		cfg = globalConfig
	}
	return &Loop{sessionID: sid, registry: reg, store: store, config: cfg}
}

func (l *Loop) Abort() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.abort != nil {
		l.abort()
	}
}

func (l *Loop) Run(prompt string) <-chan *AgentEvent {
	evChan := make(chan *AgentEvent, 64)
	ctx, cancel := context.WithCancel(context.Background())
	l.mu.Lock()
	l.abort = cancel
	l.mu.Unlock()

	go func() {
		defer close(evChan)
		defer cancel()

		l.sendEvent(evChan, "turn_start", map[string]string{"prompt": prompt})
		l.store.AddEvent(&session.Event{
			SessionID: string(l.sessionID),
			Type:      "user_message",
			Payload:   prompt,
		})

		for turn := 0; turn < l.config.MaxTurns; turn++ {
			if ctx.Err() != nil {
				break
			}
			response, err := l.processTurn(ctx, prompt)
			if err != nil {
				l.sendEvent(evChan, "error", map[string]string{"message": err.Error()})
				break
			}
			l.sendEvent(evChan, "response", map[string]string{"text": response})
			prompt = response
		}
		l.sendEvent(evChan, "turn_complete", nil)
	}()
	return evChan
}

func (l *Loop) processTurn(ctx context.Context, prompt string) (string, error) {
	result, err := l.registry.Call("bash", json.RawMessage(fmt.Sprintf(`{"command":"echo \"%s\""}`, prompt)))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (l *Loop) sendEvent(ch chan<- *AgentEvent, typ string, payload interface{}) {
	ev := &AgentEvent{Type: typ}
	if payload != nil {
		data, _ := json.Marshal(payload)
		ev.Payload = data
	}
	select {
	case ch <- ev:
	case <-time.After(5 * time.Second):
	}
}
