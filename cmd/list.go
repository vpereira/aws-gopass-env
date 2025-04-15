package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all AWS credential profiles from gopass aws store",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentAccessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
		currentProfile := strings.TrimSpace(os.Getenv("AWS_PROFILE"))

		out, err := exec.Command("gopass", "list", "--flat", "aws/").Output()
		if err != nil {
			return fmt.Errorf("failed to list secrets from gopass aws store: %w", err)
		}

		lines := bytes.Split(out, []byte{'\n'})
		for _, line := range lines {
			entry := strings.TrimSpace(string(line))
			if entry == "" || !strings.HasPrefix(entry, "aws/") {
				continue
			}
			name := strings.TrimPrefix(entry, "aws/")

			isActive := false

			// ✅ Authoritative check by AWS_PROFILE
			if name == currentProfile && currentProfile != "" {
				isActive = true
			} else if currentProfile == "" {
				// Fallback: Match by AWS_ACCESS_KEY_ID (only if AWS_PROFILE not set)
				profileOut, err := exec.Command("gopass", "show", fmt.Sprintf("aws/%s", name)).Output()
				if err == nil {
					lines := strings.Split(string(profileOut), "\n")
					for _, l := range lines {
						if strings.HasPrefix(l, "AWS_ACCESS_KEY_ID=") {
							ak := strings.TrimSpace(strings.SplitN(l, "=", 2)[1])
							if ak == currentAccessKey && ak != "" {
								isActive = true
							}
							break
						}
					}
				}
			}

			if isActive {
				fmt.Printf("* %s\n", name)
			} else {
				fmt.Printf("  %s\n", name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
