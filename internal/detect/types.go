package detect

import "errors"

var (
	ErrNoPackageJSON = errors.New("package.json or deno.json not found")
	ErrNoDevScript   = errors.New("no 'dev' or 'start' script found")
)

type DetectResult struct {
	Command        string
	Args           []string
	PackageManager string
	ScriptName     string
	IsManual       bool
}
