package importrules

import (
	"context"
	"os/exec"

	"github.com/payfazz/go-errors/v2"
)

type commandError struct {
	msg    string
	detail string
}

func (err *commandError) Error() string       { return err.msg }
func (err *commandError) ErrorDetail() string { return err.detail }

func command(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	stdout, err := cmd.Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return "", errors.Trace(&commandError{
				msg:    "failed to execute: " + cmd.String(),
				detail: string(exitError.Stderr),
			})
		}
		return "", errors.Trace(err)
	}
	return string(stdout), nil
}
