package core

// UsageEvent represents a single usage log entry
// Stored at ~/.config/linesense/usage.log
type UsageEvent struct {
	Timestamp string `json:"timestamp"` // ISO 8601
	CWD       string `json:"cwd"`
	Command   string `json:"command"`
	Accepted  bool   `json:"accepted"` // whether the user executed it as suggested
	Source    string `json:"source"`   // "preset" | "llm"
}

// LogUsage appends a usage event to the usage log
func LogUsage(_ UsageEvent) error {
	// TODO: Resolve usage log path (~/.config/linesense/usage.log)
	// TODO: Ensure directory exists
	// TODO: Append JSON line to usage log
	panic("not implemented")
}

// BuildUsageSummary creates a summary of usage patterns for a given cwd
func BuildUsageSummary(_ string) (*UsageSummary, error) {
	// TODO: Read usage log
	// TODO: Filter events by cwd
	// TODO: Count command frequency
	// TODO: Count preset usage
	// TODO: Return top N frequently used commands
	// TODO: Return nil if usage log doesn't exist
	panic("not implemented")
}

