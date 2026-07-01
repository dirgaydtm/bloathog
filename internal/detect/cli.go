package detect

import (
	"flag"
	"fmt"
	"io"
)

const usageText = "Usage:\n  bloathog              # auto-detect\n  bloathog <cmd>        # manual mode"

// ParseArgs parses the command line arguments and extracts positional arguments.
// It returns an error if unknown flags are passed or if help is requested.
func ParseArgs(args []string) ([]string, error) {
	fs := flag.NewFlagSet("bloathog", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // Mute default stderr output

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil, fmt.Errorf(usageText)
		}
		return nil, fmt.Errorf("unknown flag provided\n\n%s", usageText)
	}

	return fs.Args(), nil
}
