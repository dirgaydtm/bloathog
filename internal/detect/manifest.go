package detect

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
)

// FindDevScript checks package.json or deno.json for a "dev" or "start" script.
func FindDevScript(dir string) (string, error) {
	// Try package.json
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err == nil {
		if script := gjson.GetBytes(data, "scripts.dev").String(); script != "" {
			return "dev", nil
		}
		if script := gjson.GetBytes(data, "scripts.start").String(); script != "" {
			return "start", nil
		}
		return "", fmt.Errorf("%w\n\nTip: run with a manual command:\n  bloathog <cmd>", ErrNoDevScript)
	}

	// Try deno.json
	data, err = os.ReadFile(filepath.Join(dir, "deno.json"))
	if err == nil {
		if script := gjson.GetBytes(data, "tasks.dev").String(); script != "" {
			return "dev", nil
		}
		if script := gjson.GetBytes(data, "tasks.start").String(); script != "" {
			return "start", nil
		}
		return "", fmt.Errorf("%w\n\nTip: run with a manual command:\n  bloathog <cmd>", ErrNoDevScript)
	}

	return "", fmt.Errorf("%w\n\nTip: run with a manual command:\n  bloathog <cmd>", ErrNoPackageJSON)
}
