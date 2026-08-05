package types

// MonitorState holds current stats.
type MonitorState struct {
	CurrentRSS      uint64
	PeakRSS         uint64
	RunningSumRSS   float64
	CurrentCPU      float64
	PeakCPU         float64
	RunningSumCPU   float64
	SmoothedCPU     float64
	SampleCount     uint64
	ActiveProcesses int
	PeakProcesses   int
}

// Returns the average RSS in megabytes.
func (s MonitorState) AverageRSSMB() float64 {
	if s.SampleCount == 0 {
		return 0
	}
	return (s.RunningSumRSS / float64(s.SampleCount)) / (1024 * 1024)
}

// Returns the average CPU percentage.
func (s MonitorState) AverageCPU() float64 {
	if s.SampleCount == 0 {
		return 0
	}
	return s.RunningSumCPU / float64(s.SampleCount)
}

// Returns the current RSS in megabytes.
func (s MonitorState) CurrentRSSMB() float64 {
	return float64(s.CurrentRSS) / (1024 * 1024)
}

// Returns the peak RSS in megabytes.
func (s MonitorState) PeakRSSMB() float64 {
	return float64(s.PeakRSS) / (1024 * 1024)
}

// ProjectInfo holds project information.
type ProjectInfo struct {
	Command        string
	Args           []string
	PackageManager string // npm, yarn, pnpm, bun, or "" for manual
	ScriptName     string // "dev", "start", or "" for manual
	IsManual       bool   // true if command was provided via --
}
