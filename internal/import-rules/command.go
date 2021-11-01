package importrules

import (
	"context"
	"os/exec"
	"strings"

	"github.com/payfazz/go-errors/v2"
)

func command(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	stdout, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			stderr := string(exitErr.Stderr)
			var sb strings.Builder
			sb.WriteString("cannot execute \"")
			sb.WriteString(cmd.String())
			sb.WriteString("\"\n")
			sb.WriteString(stderr)
			if !strings.HasSuffix(stderr, "\n") {
				sb.WriteString("\n")
			}
			return "", errors.New(sb.String())
		}
		return "", errors.Trace(err)
	}
	return string(stdout), nil
}
