package sheller

import "time"

// CommandResult holds the result of a command execution.
type CommandResult struct {
	Succeeded      bool
	ExitCode       int
	StandardOutput string
	StandardError  string
	StartTime      time.Time
	EndTime        time.Time
	RunTime        time.Duration
}
