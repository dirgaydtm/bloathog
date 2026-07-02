package types

import "time"

// MemorySample represents a single RSS measurement at a point in time.
type MemorySample struct {
	Timestamp time.Time
	RSS       uint64 // bytes
}

// MonitorState holds the current monitoring statistics.
type MonitorState struct {
	CurrentRSS      uint64
	PeakRSS         uint64
	RunningSumRSS   float64
	CurrentCPU      float64
	PeakCPU         float64
	RunningSumCPU   float64
	SampleCount     uint64
	ActiveProcesses int
	GraphBuffer     []float64 // ring buffer values in MB, up to 120 entries
	LogBuffer       []string  // ring buffer, up to 5000 lines
}

// AverageRSSMB returns the average RSS in megabytes.
func (s MonitorState) AverageRSSMB() float64 {
	if s.SampleCount == 0 {
		return 0
	}
	return (s.RunningSumRSS / float64(s.SampleCount)) / (1024 * 1024)
}

// AverageCPU returns the average CPU percentage.
func (s MonitorState) AverageCPU() float64 {
	if s.SampleCount == 0 {
		return 0
	}
	return s.RunningSumCPU / float64(s.SampleCount)
}

// CurrentRSSMB returns the current RSS in megabytes.
func (s MonitorState) CurrentRSSMB() float64 {
	return float64(s.CurrentRSS) / (1024 * 1024)
}

// PeakRSSMB returns the peak RSS in megabytes.
func (s MonitorState) PeakRSSMB() float64 {
	return float64(s.PeakRSS) / (1024 * 1024)
}

// ProjectInfo holds information about the project/command to run.
type ProjectInfo struct {
	Command        string
	Args           []string
	PackageManager string // npm, yarn, pnpm, bun, or "" for manual
	ScriptName     string // "dev", "start", or "" for manual
	IsManual       bool   // true if command was provided via --
}
