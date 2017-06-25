package git

import (
	"testing"
	"os/exec"
	"strings"
)

func GetCurrentBranchName(t *testing.T) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	bytes, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to determine current branch name due to error: %s\n", err)
	}
	return strings.TrimSpace(string(bytes))
}
