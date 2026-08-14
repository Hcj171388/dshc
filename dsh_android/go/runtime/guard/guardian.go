package guard

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Guardian monitors and limits agent execution
type Guardian struct {
	mu              sync.Mutex
	turnCount       int
	maxTurns        int
	stepCount       int
	maxSteps        int
	toolCallCount   int
	maxToolCalls    int
	lastActivity    time.Time
	activityTimeout time.Duration
	aborted         bool
}

// Config holds guardian configuration
type Config struct {
	MaxTurns        int
	MaxSteps        int
	MaxToolCalls    int
	ActivityTimeout time.Duration
}

// DefaultConfig returns default guardian configuration
func DefaultConfig() Config {
	return Config{
		MaxTurns:        30,
		MaxSteps:        50,
		MaxToolCalls:    100,
		ActivityTimeout: 5 * time.Minute,
	}
}

// NewGuardian creates a new guardian
func NewGuardian(cfg Config) *Guardian {
	if cfg.MaxTurns == 0 {
		cfg = DefaultConfig()
	}
	return &Guardian{
		maxTurns:        cfg.MaxTurns,
		maxSteps:        cfg.MaxSteps,
		maxToolCalls:    cfg.MaxToolCalls,
		activityTimeout: cfg.ActivityTimeout,
		lastActivity:    time.Now(),
	}
}

// CheckTurn checks if we can start a new turn
func (g *Guardian) CheckTurn() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.turnCount++
	if g.turnCount > g.maxTurns {
		return &GuardError{Code: "max_turns", Message: fmt.Sprintf("maximum turns exceeded: %d", g.maxTurns)}
	}
	
	g.lastActivity = time.Now()
	return nil
}

// CheckStep checks if we can execute another step
func (g *Guardian) CheckStep() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.stepCount++
	if g.stepCount > g.maxSteps {
		return &GuardError{Code: "max_steps", Message: fmt.Sprintf("maximum steps exceeded: %d", g.maxSteps)}
	}
	
	g.lastActivity = time.Now()
	return nil
}

// CheckToolCall checks if we can execute another tool call
func (g *Guardian) CheckToolCall() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.toolCallCount++
	if g.toolCallCount > g.maxToolCalls {
		return &GuardError{Code: "max_tool_calls", Message: fmt.Sprintf("maximum tool calls exceeded: %d", g.maxToolCalls)}
	}
	
	g.lastActivity = time.Now()
	return nil
}

// CheckActivity checks if the agent is still active
func (g *Guardian) CheckActivity() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.activityTimeout > 0 && time.Since(g.lastActivity) > g.activityTimeout {
		return &GuardError{Code: "activity_timeout", Message: "agent activity timeout exceeded"}
	}
	
	return nil
}

// IsAborted checks if the guardian has been aborted
func (g *Guardian) IsAborted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.aborted
}

// Abort signals the guardian to stop
func (g *Guardian) Abort() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.aborted = true
}

// Reset resets all counters
func (g *Guardian) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.turnCount = 0
	g.stepCount = 0
	g.toolCallCount = 0
	g.lastActivity = time.Now()
	g.aborted = false
}

// Stats returns current guardian statistics
func (g *Guardian) Stats() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return map[string]int{
		"turn_count":     g.turnCount,
		"step_count":     g.stepCount,
		"tool_call_count": g.toolCallCount,
	}
}

// GuardianError represents an error from the guardian
type GuardError struct {
	Code    string
	Message string
}

func (e *GuardError) Error() string {
	return e.Message
}

// IsMaxTurnsError checks if error is max turns exceeded
func IsMaxTurnsError(err error) bool {
	if e, ok := err.(*GuardError); ok {
		return e.Code == "max_turns"
	}
	return false
}

// IsAborted checks if context is aborted
func IsAborted(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
