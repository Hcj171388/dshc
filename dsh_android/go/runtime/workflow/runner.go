package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Step represents a single step in a workflow
type Step struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Action      string          `json:"action"`
	Args        json.RawMessage `json:"args"`
	NextStepID  string          `json:"next_step_id,omitempty"`
	OnError     string          `json:"on_error,omitempty"` // "continue", "abort"
	TimeoutMs   int             `json:"timeout_ms,omitempty"`
}

// Workflow represents a sequence of steps
type Workflow struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Steps   []Step                 `json:"steps"`
	Context map[string]interface{} `json:"context"`
}

// RunResult represents the result of running a workflow
type RunResult struct {
	ID        string         `json:"id"`
	Completed bool           `json:"completed"`
	Steps     []StepResult   `json:"steps"`
	Error     string         `json:"error,omitempty"`
	Duration  time.Duration  `json:"duration"`
}

// StepResult represents the result of a single step
type StepResult struct {
	StepID   string      `json:"step_id"`
	Name     string      `json:"name"`
	Success  bool        `json:"success"`
	Output   interface{} `json:"output,omitempty"`
	Error    string      `json:"error,omitempty"`
	Duration time.Duration `json:"duration"`
}

// Runner executes workflows
type Runner struct {
	mu      sync.Mutex
	running map[string]*runState
}

type runState struct {
	wf     *Workflow
	ctx    context.Context
	cancel context.CancelFunc
	result *RunResult
}

// NewRunner creates a new workflow runner
func NewRunner() *Runner {
	return &Runner{
		running: make(map[string]*runState),
	}
}

// Run starts a workflow execution
func (r *Runner) Run(wf *Workflow) (*RunResult, error) {
	r.mu.Lock()
	if _, exists := r.running[wf.ID]; exists {
		r.mu.Unlock()
		return nil, fmt.Errorf("workflow already running: %s", wf.ID)
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	state := &runState{
		wf:     wf,
		ctx:    ctx,
		cancel: cancel,
		result: &RunResult{
			ID:        wf.ID,
			Completed: false,
			Steps:     make([]StepResult, 0),
		},
	}
	r.running[wf.ID] = state
	r.mu.Unlock()
	
	go r.execute(state)
	return state.result, nil
}

func (r *Runner) execute(state *runState) {
	startTime := time.Now()
	defer func() {
		state.result.Duration = time.Since(startTime)
		r.mu.Lock()
		delete(r.running, state.wf.ID)
		r.mu.Unlock()
	}()
	
	for _, step := range state.wf.Steps {
		if state.ctx.Err() != nil {
			return
		}
		
		stepStart := time.Now()
		stepResult := StepResult{
			StepID: step.ID,
			Name:   step.Name,
		}
		
		// Execute step
		output, err := r.executeStep(state.ctx, step)
		stepResult.Duration = time.Since(stepStart)
		
		if err != nil {
			stepResult.Success = false
			stepResult.Error = err.Error()
			state.result.Steps = append(state.result.Steps, stepResult)
			
			if step.OnError == "abort" {
				state.result.Error = err.Error()
				return
			}
			continue
		}
		
		stepResult.Success = true
		stepResult.Output = output
		state.result.Steps = append(state.result.Steps, stepResult)
	}
	
	state.result.Completed = true
}

func (r *Runner) executeStep(ctx context.Context, step Step) (interface{}, error) {
	switch step.Action {
	case "bash":
		return r.execBash(ctx, step)
	case "wait":
		return r.waitForCondition(ctx, step)
	case "branch":
		return r.evaluateBranch(ctx, step)
	default:
		return nil, fmt.Errorf("unknown action: %s", step.Action)
	}
}

func (r *Runner) execBash(ctx context.Context, step Step) (string, error) {
	var params struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(step.Args, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Executed: %s", params.Command), nil
}

func (r *Runner) waitForCondition(ctx context.Context, step Step) (bool, error) {
	var params struct {
		TimeoutMs int `json:"timeout_ms"`
	}
	if err := json.Unmarshal(step.Args, &params); err != nil {
		return false, err
	}
	
	timer := time.NewTimer(time.Duration(params.TimeoutMs) * time.Millisecond)
	select {
	case <-ctx.Done():
		timer.Stop()
		return false, ctx.Err()
	case <-timer.C:
		return true, nil
	}
}

func (r *Runner) evaluateBranch(ctx context.Context, step Step) (string, error) {
	return step.NextStepID, nil
}

// Stop stops a running workflow
func (r *Runner) Stop(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if state, ok := r.running[id]; ok {
		state.cancel()
	}
}

// IsRunning checks if a workflow is running
func (r *Runner) IsRunning(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.running[id]
	return ok
}
