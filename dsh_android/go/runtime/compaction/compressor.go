package compaction

import (
	"fmt"
	"strings"
)

// Compressor handles session compaction
type Compressor struct {
	Mode       string // "summary", "truncate", "none"
	Threshold  int    // number of events before compaction
	SummaryLen int    // max length of summary
}

// CompactResult holds the result of compaction
type CompactResult struct {
	OriginalCount int      `json:"original_count"`
	CompactCount  int      `json:"compact_count"`
	Summary       string   `json:"summary,omitempty"`
	Events        []Event  `json:"events"`
}

// Event represents a session event
type Event struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// NewCompressor creates a new compressor
func NewCompressor(mode string, threshold, summaryLen int) *Compressor {
	if mode == "" {
		mode = "summary"
	}
	if threshold == 0 {
		threshold = 100
	}
	if summaryLen == 0 {
		summaryLen = 500
	}
	return &Compressor{
		Mode:       mode,
		Threshold:  threshold,
		SummaryLen: summaryLen,
	}
}

// Compact performs compaction on events
func (c *Compressor) Compact(events []Event) (*CompactResult, error) {
	if len(events) <= c.Threshold {
		return &CompactResult{
			OriginalCount: len(events),
			CompactCount:  len(events),
			Events:        events,
		}, nil
	}
	
	var compacted []Event
	var summary strings.Builder
	
	switch c.Mode {
	case "summary":
		compacted = c.summarize(events, &summary)
	case "truncate":
		compacted = c.truncate(events)
	default:
		compacted = events
	}
	
	return &CompactResult{
		OriginalCount: len(events),
		CompactCount:  len(compacted),
		Summary:       summary.String(),
		Events:        compacted,
	}, nil
}

func (c *Compressor) summarize(events []Event, summary *strings.Builder) []Event {
	turns := 0
	userMsgs := 0
	toolCalls := 0
	
	for _, ev := range events {
		switch ev.Type {
		case "turn/start":
			turns++
		case "user_message":
			userMsgs++
		case "tool_call":
			toolCalls++
		}
	}
	
	summary.WriteString(fmt.Sprintf("Session Summary:\n"))
	summary.WriteString(fmt.Sprintf("- Total turns: %d\n", turns))
	summary.WriteString(fmt.Sprintf("- User messages: %d\n", userMsgs))
	summary.WriteString(fmt.Sprintf("- Tool calls: %d\n", toolCalls))
	
	keep := min(c.Threshold, len(events))
	start := len(events) - keep
	return events[start:]
}

func (c *Compressor) truncate(events []Event) []Event {
	keep := min(c.Threshold, len(events))
	start := len(events) - keep
	return events[start:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CheckSize estimates the size of events
func CheckSize(events []Event) int {
	total := 0
	for _, ev := range events {
		total += len(ev.Type) + len(ev.Payload)
	}
	return total
}
