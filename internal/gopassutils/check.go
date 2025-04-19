package gopassutils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// gopassEntryExists checks whether a given path exists in gopass
func EntryExists(path string) (bool, error) {
	cmd := exec.Command("gopass", "ls", "--flat")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list gopass entries: %w", err)
	}

	lines := bytes.Split(output, []byte{'\n'})
	for _, line := range lines {
		if strings.TrimSpace(string(line)) == path {
			return true, nil
		}
	}

	return false, nil
}
