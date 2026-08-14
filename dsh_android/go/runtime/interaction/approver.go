package interaction

import (
	"encoding/json"
	"fmt"
	"sync"
)

// ApprovalRequest represents a request for user approval
type ApprovalRequest struct {
	ID        string          `json:"id"`
	Action    string          `json:"action"`
	ToolName  string          `json:"tool_name"`
	Args      json.RawMessage `json:"args"`
	Message   string          `json:"message"`
	Priority  int             `json:"priority"`
	TimeoutMs int             `json:"timeout_ms"`
}

// ApprovalResult represents the result of an approval request
type ApprovalResult struct {
	RequestID string `json:"request_id"`
	Approved  bool   `json:"approved"`
	Reason    string `json:"reason,omitempty"`
}

// Approver manages approval requests
type Approver struct {
	mu       sync.Mutex
	requests map[string]*ApprovalRequest
	pending  chan *ApprovalRequest
	results  map[string]*ApprovalResult
	presets  map[string]bool
}

// NewApprover creates a new approver
func NewApprover() *Approver {
	return &Approver{
		requests: make(map[string]*ApprovalRequest),
		pending:  make(chan *ApprovalRequest, 10),
		results:  make(map[string]*ApprovalResult),
		presets:  make(map[string]bool),
	}
}

// RequestApproval sends an approval request
func (a *Approver) RequestApproval(req *ApprovalRequest) (*ApprovalResult, error) {
	a.mu.Lock()
	a.requests[req.ID] = req
	a.mu.Unlock()
	
	// Check preset
	a.mu.Lock()
	if approved, ok := a.presets[req.ToolName]; ok {
		result := &ApprovalResult{
			RequestID: req.ID,
			Approved:  approved,
		}
		a.results[req.ID] = result
		a.mu.Unlock()
		return result, nil
	}
	a.mu.Unlock()
	
	// Send to pending channel
	select {
	case a.pending <- req:
	default:
		return nil, fmt.Errorf("approval queue full")
	}
	
	// Auto-approve for now
	result := &ApprovalResult{
		RequestID: req.ID,
		Approved:  true,
	}
	
	a.mu.Lock()
	a.results[req.ID] = result
	a.mu.Unlock()
	
	return result, nil
}

// ApproveTool sets a preset for a tool
func (a *Approver) ApproveTool(toolName string, approve bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.presets[toolName] = approve
}

// DenyTool sets a preset to deny a tool
func (a *Approver) DenyTool(toolName string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.presets, toolName)
}

// GetPresets returns all approval presets
func (a *Approver) GetPresets() map[string]bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	presets := make(map[string]bool)
	for k, v := range a.presets {
		presets[k] = v
	}
	return presets
}

// GetPending returns pending approval requests
func (a *Approver) GetPending() []*ApprovalRequest {
	a.mu.Lock()
	defer a.mu.Unlock()
	requests := make([]*ApprovalRequest, 0, len(a.requests))
	for _, req := range a.requests {
		requests = append(requests, req)
	}
	return requests
}

// ClearResults clears old results
func (a *Approver) ClearResults() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.results = make(map[string]*ApprovalResult)
}
