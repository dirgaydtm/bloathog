package detect

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
)

// Finds the package manager by walking up the directory tree.
func DetectPackageManager(dir string) string {
	locks := [][2]string{
		{"npm-shrinkwrap.json", "npm"}, {"package-lock.json", "npm"}, {"yarn.lock", "yarn"},
		{"pnpm-lock.yaml", "pnpm"}, {"bun.lockb", "bun"}, {"bun.lock", "bun"},
		{"deno.json", "deno"}, {"deno.jsonc", "deno"}, {"deno.lock", "deno"},
	}

	dir, _ = filepath.Abs(dir)
	for {
		// 1. Check packageManager field (Corepack)
		if d, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
			if pm := gjson.GetBytes(d, "packageManager").String(); pm != "" {
				return strings.Split(pm, "@")[0]
			}
		}

		// 2. Check lockfiles
		for _, lf := range locks {
			if _, err := os.Stat(filepath.Join(dir, lf[0])); err == nil {
				return lf[1]
			}
		}

		// 3. Stop at repository or filesystem root
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil || dir == filepath.Dir(dir) {
			break
		}
		dir = filepath.Dir(dir)
	}
	return "npm"
}
