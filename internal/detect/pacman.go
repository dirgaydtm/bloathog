package detect

import (
	"os"
	"path/filepath"
)

// DetectPackageManager finds the project's package manager via lockfiles (defaults to npm).
func DetectPackageManager(dir string) string {
	lockfiles := []struct{ name, pm string }{
		{"deno.json", "deno"},
		{"deno.jsonc", "deno"},
		{"deno.lock", "deno"},
		{"bun.lockb", "bun"},
		{"bun.lock", "bun"},
		{"pnpm-lock.yaml", "pnpm"},
		{"yarn.lock", "yarn"},
	}

	for _, lf := range lockfiles {
		if _, err := os.Stat(filepath.Join(dir, lf.name)); err == nil {
			return lf.pm
		}
	}
	
	return "npm"
}