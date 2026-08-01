package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dirgaa/bloathog/internal/detect"
	"github.com/dirgaa/bloathog/internal/ui"
	"github.com/dirgaa/bloathog/internal/ui/types"
)

func main() {
	os.Exit(run())
}

func run() int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, ui.RenderError(err))
		return 1
	}

	result, err := detect.Detect(cwd, os.Args[1:])
	if err != nil {
		if strings.HasPrefix(err.Error(), "Usage:") {
			fmt.Printf("\n  %s\n\n", err.Error())
			return 0
		}
		fmt.Fprint(os.Stderr, ui.RenderError(err))
		return 1
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
		fmt.Fprint(os.Stderr, ui.RenderError(fmt.Errorf("bloathog error: %w", err)))
		return 1
	}

	if final, ok := m.(ui.Model); ok {
		if final.FatalErr != nil {
			fmt.Fprint(os.Stderr, ui.RenderError(final.FatalErr))
			code := final.ExitCode
			if code == 0 {
				code = 1
			}
			return code
		}
		if report := final.View(); report != "" {
			fmt.Println(report)
		}
		return final.ExitCode
	}

	return 0
}
