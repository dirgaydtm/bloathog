package detect

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dirgaa/bloathog/internal/eror"
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
			return DetectResult{}, &eror.Error{Msg: fmt.Sprintf("executable '%s' not found in PATH", cmd), Tip: "check your command spelling"}
		}

		// Enforce package.json for known JS package managers
		if cmd == "npm" || cmd == "yarn" || cmd == "pnpm" || cmd == "bun" {
			if _, err := os.Stat(filepath.Join(dir, "package.json")); os.IsNotExist(err) {
				return DetectResult{}, &eror.Error{Msg: "missing package.json", Tip: "check if you are in the correct directory"}
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
		return DetectResult{}, &eror.Error{Msg: fmt.Sprintf("package manager '%s' not found in PATH", pm), Tip: fmt.Sprintf("install %s or use manual mode", pm)}
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
