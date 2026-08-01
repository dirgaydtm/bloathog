package ui

import (
	"errors"
	"fmt"

	"github.com/dirgaa/bloathog/internal/eror"
)

// RenderError takes any error and returns a styled string.
func RenderError(err error) string {
	msg := err.Error()
	tip := ""

	var appErr *eror.Error
	if errors.As(err, &appErr) {
		msg = appErr.Msg
		if appErr.Tip != "" {
			tip = fmt.Sprintf("\n\033[1;33m💡  TIP:\033[0m %s", appErr.Tip)
		}
	}
	return fmt.Sprintf("\033[31m🚨ERROR:\033[0m %s%s\n", msg, tip)
}
