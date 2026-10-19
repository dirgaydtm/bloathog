package proc

import (
	"slices"

	"github.com/shirou/gopsutil/v3/process"
)

// TreeStats holds the aggregate result of a process tree scan.
type TreeStats struct {
	TotalRSS     uint64        // bytes, sum of RSS for all processes
	TotalCPU     float64       // sum of CPU percentage
	ProcessCount int           // number of live processes found
	Nodes        []ProcessNode // flat list of all scanned processes
}

// ProcessNode represents a single process entry in the tree.
type ProcessNode struct {
	PID    int32
	Name   string
	RSS    uint64  // bytes
	CPU    float64 // percentage
	Depth  int     // 0 = root, 1 = direct child, etc.
	IsLast []bool  // tracks if ancestors/self are the last child in their level
}

// procs persists between ticks so Percent(0) can measure CPU delta.
var procs = map[int32]*process.Process{}

// CollectTree performs a DFS traversal of the process tree rooted at rootPID,
// summing RSS across all live descendants. Dead/inaccessible processes are silently skipped.
func CollectTree(rootPID int32) TreeStats {
	var stats TreeStats
	seen := map[int32]*process.Process{}
	collectDFS(rootPID, 0, nil, &stats, seen)
	procs = seen
	return stats
}

func collectDFS(pid int32, depth int, isLast []bool, stats *TreeStats, seen map[int32]*process.Process) {
	p, ok := procs[pid]
	if !ok {
		var err error
		if p, err = process.NewProcess(pid); err != nil {
			return // process gone
		}
	}
	seen[pid] = p

	if mem, err := p.MemoryInfo(); err == nil {
		name, _ := p.Name() // best-effort
		cpu, _ := p.Percent(0)
		stats.TotalRSS += mem.RSS
		stats.TotalCPU += cpu
		stats.ProcessCount++
		stats.Nodes = append(stats.Nodes, ProcessNode{
			PID:    pid,
			Name:   name,
			RSS:    mem.RSS,
			CPU:    cpu,
			Depth:  depth,
			IsLast: isLast,
		})
	}

	children, _ := p.Children() // safely ignores errors
	for i, child := range children {
		collectDFS(child.Pid, depth+1, append(slices.Clip(isLast), i == len(children)-1), stats, seen)
	}
}
