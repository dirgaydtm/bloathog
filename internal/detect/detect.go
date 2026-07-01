package detect

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Detect parses args for manual mode or auto-detects the project via manifests.
func Detect(dir string, args []string) (DetectResult, error) {
	parsedArgs, err := ParseArgs(args)
	if err != nil {
		return DetectResult{}, err
	}

	// Manual mode
	if len(parsedArgs) > 0 {
		cmd := parsedArgs[0]
		if _, err := exec.LookPath(cmd); err != nil {
			return DetectResult{}, fmt.Errorf("executable '%s' not found in PATH\n\nTip: check your command spelling", cmd)
		}

		// Enforce package.json for known JS package managers
		if cmd == "npm" || cmd == "yarn" || cmd == "pnpm" || cmd == "bun" {
			if _, err := os.Stat(filepath.Join(dir, "package.json")); os.IsNotExist(err) {
				return DetectResult{}, fmt.Errorf("%w\n\nTip: check if you are in the correct directory", ErrNoPackageJSON)
			}
		}

		return DetectResult{Command: cmd, Args: parsedArgs[1:], IsManual: true}, nil
	}

	// Auto-detect mode
	scriptName, err := FindDevScript(dir)
	if err != nil {
		return DetectResult{}, err
	}

	pm := DetectPackageManager(dir)
	if _, err := exec.LookPath(pm); err != nil {
		return DetectResult{}, fmt.Errorf("package manager '%s' not found in PATH\n\nTip: install %s or use manual mode", pm, pm)
	}

	runCommand := "run"
	if pm == "deno" {
		runCommand = "task"
	}

	return DetectResult{
		Command:        pm,
		Args:           []string{runCommand, scriptName},
		PackageManager: pm,
		ScriptName:     scriptName,
	}, nil
}
