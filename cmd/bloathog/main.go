package main

import (
	"github.com/dirgaa/bloathog/internal/ui/types"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/detect"

	"github.com/dirgaa/bloathog/internal/ui"
)

func fatal(err error) {
	msg := err.Error()
	if strings.HasPrefix(msg, "Usage:") {
		fmt.Printf("\n  %s\n\n", msg)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "\n  \033[31mERROR:\033[0m %s\n\n", msg)
	os.Exit(1)
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	result, err := detect.Detect(cwd, os.Args[1:])
	if err != nil {
		fatal(err)
	}

	p := tea.NewProgram(ui.NewModel(types.ProjectInfo{
		Command:        result.Command,
		Args:           result.Args,
		PackageManager: result.PackageManager,
		ScriptName:     result.ScriptName,
		IsManual:       result.IsManual,
	}), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		fatal(fmt.Errorf("bloathog error: %w", err))
	}

	if final, ok := m.(ui.Model); ok {
		if report := final.View(); report != "" {
			fmt.Println(report)
		}
	}
}
