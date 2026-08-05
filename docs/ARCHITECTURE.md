# Bloathog System Architecture

This document provides a comprehensive technical overview of the Bloathog architecture. It is intended for core contributors, maintainers, and developers seeking a deep understanding of the system's internal subsystems, concurrency model, and rendering pipeline.

## 1. System Overview

Bloathog is a Terminal User Interface (TUI) application designed for real-time aggregation and visualization of system resource metrics (Memory/RSS and CPU) across an entire process tree.

Unlike standard utilities (e.g., `top` or `htop`) that report metrics on a per-PID basis, Bloathog acts as a parent wrapper. It executes a target command within a new process group, recursively tracks all descendant child processes, and aggregates their resource footprints into a unified metric stream.

## 2. Technology Stack

- **Language:** [Go (Golang)](https://go.dev/) (1.26+)
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) (State management and event loop based on The Elm Architecture)
- **UI Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss) (CSS-like styling definitions for the terminal)
- **Graphing Engine:** [Asciigraph](https://github.com/guptarohit/asciigraph) (Terminal-based data plotting)
- **OS/Metrics Integration:** [gopsutil/v3](https://github.com/shirou/gopsutil) (Cross-platform system and process utility abstractions)
- **JSON Parsing:** [gjson](https://github.com/tidwall/gjson) (High-performance JSON extraction for manifest detection)

## 3. Directory Structure

The repository is strictly organized into targeted, internal packages:

```text
bloathog/
├── cmd/
│   └── bloathog/
│       └── main.go           # Application entry point
├── docs/                     # Architecture and Contribution guidelines
├── internal/
│   ├── detect/               # Auto-detection for JS package managers and script manifests
│   ├── eror/                 # Core error definitions
│   ├── monitor/              # Subprocess lifecycle, event streaming, and ticking engine
│   ├── proc/                 # OS-level process tree traversal & raw metrics aggregation
│   ├── ring/                 # Generic O(1) circular buffers for metrics
│   └── ui/                   # TUI Presentation Layer
│       ├── components/       # UI widgets (graph, header, help, log, process, report)
│       ├── theme/            # Centralized Lip Gloss styles and color palettes
│       ├── types/            # Shared interfaces and centralized keybindings
│       ├── error.go          # UI-level error screen rendering
│       ├── io.go             # Subprocess stdout/stderr asynchronous scanners
│       ├── layout.go         # Responsive geometry engine
│       ├── model.go          # Root Bubble Tea model and state container
│       ├── quit.go           # Graceful shutdown sequencer
│       └── render.go         # Master view compositor
├── go.mod                    # Module definition and dependencies
└── go.sum                    # Dependency checksums
```

## 4. Subsystem Detailed Design

### 4.1. Project Detection (`internal/detect`)
Before initialization, Bloathog analyzes the user's workspace. If no explicit command is provided, it recursively traverses parent directories up to the system root to support monorepos. It parses project manifests (e.g., `package.json`), prioritizing explicit package manager declarations before falling back to lockfiles to accurately infer the target command.

### 4.2. Presentation Layer (`internal/ui`)
The UI layer operates on a single-threaded event loop provided by Bubble Tea. State mutations occur exclusively within the `Update` function to guarantee thread safety.

- **`model.go`:** The central state container holding metrics buffers, bounded logs, and interactive input states.
- **`layout.go` & `render.go`:** Calculates dynamic viewport constraints to implement responsive scaling, dynamically resizing panels based on terminal dimensions.
- **`components/`:** Modular struct definitions encapsulating their own `View()` logic. CPU graphs utilize time-series smoothing to prevent jagged visual artifacts and >100% spikes.
- **`io.go`:** Handles asynchronous scanning of standard streams, batching incoming log streams before feeding them to the main event loop to maintain UI responsiveness under heavy I/O.

### 4.3. Monitoring Engine (`internal/monitor`)
The monitoring engine operates asynchronously alongside the main UI thread.

- **Process Spawning (`spawn.go`):** Executes target commands within a dedicated Process Group ID (PGID), ensuring orphaned child processes can be reliably tracked and terminated.
- **Ticker Loop (`tick.go`):** A background loop that continuously issues a tick command to the main application to trigger a fresh metrics pull at fixed intervals.

### 4.4. Metrics Collection (`internal/proc`)
The metrics package (`proc.go`) integrates with `shirou/gopsutil` to abstract operating system complexities.

- **Process Tree Traversal:** Performs a cross-platform Depth-First Search (DFS) over the process tree using `process.Children()` to build an accurate representation of all active descendants.
- **Metric Aggregation:**
  - **Resident Set Size (RSS):** Iteratively polls and sums the memory usage across all active nodes.
  - **CPU Utilization:** Iteratively polls the CPU utilization across all active nodes. Dead or inaccessible processes are silently skipped to maintain loop stability.

## 5. Concurrency Model and Event Loop

Bloathog adheres to a strict message-passing concurrency model to avoid mutex contention.

1. **Asynchronous Producers:**
   - **Ticker Goroutine:** Emits tick messages at a fixed interval.
   - **I/O Goroutines:** Continuously scan `stdout`/`stderr` streams and emit string messages.
   - **Wait Goroutine:** Blocks until the target process completes and emits an exit event.

2. **Synchronous Consumer:**
   The `Update(msg tea.Msg)` function receives all messages sequentially, applies data transformations (appending to arrays, calculating peaks), and triggers view re-renders.

## 6. Data Structures and Memory Management

To maintain low overhead, Bloathog strictly bounds internal memory:

- **Generic Ring Buffers:** Graph data is stored in generic circular buffers to ensure O(1) time complexity for insertions and a stable memory footprint.
- **Bounded Log Buffers:** Log histories are strictly capped to prevent Out-Of-Memory (OOM) crashes during verbose sessions.
- **Style Caching:** `lipgloss.Style` instances are defined as package-level variables rather than being re-allocated on every frame, drastically reducing Garbage Collection (GC) pressure.

## 7. Shutdown Sequence

Clean termination prevents zombie processes on the host machine.

1. Upon receiving an exit signal (e.g., `ctrl+c`), a termination signal is dispatched to the entire Process Group (PGID).
2. The UI enters a terminating state and blocks until background workers confirm the exit sequence.
3. The final Session Summary report (`report.go`) is rendered to `stdout`.
