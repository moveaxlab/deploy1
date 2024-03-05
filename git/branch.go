package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func CurrentBranch() (string, error) {
	res, err := exec.Command(
		"git",
		"rev-parse",
		"--abbrev-ref", "HEAD",
	).Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.Trim(string(res), "\n "), nil
}
