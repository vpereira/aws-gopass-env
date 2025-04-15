package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var showPassword = true

var showCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show the stored AWS credentials for the given profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := fmt.Sprintf("aws/%s", name)

		out, err := exec.Command("gopass", "show", path).Output()
		if err != nil {
			return fmt.Errorf("failed to fetch profile '%s': %w", name, err)
		}

		fmt.Printf("Profile: %s\n", name)
		fmt.Println(strings.Repeat("-", 30))

		lines := bytes.Split(out, []byte{'\n'})
		for _, line := range lines {
			if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
				continue
			}

			parts := strings.SplitN(string(line), "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			if key == "AWS_SECRET_ACCESS_KEY" && !showPassword {
				val = maskString(val)
			}

			fmt.Printf("%s=%s\n", key, val)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&showPassword, "show-password", false, "Show the AWS secret access key in plain text")
}

func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + strings.Repeat("*", len(s)-4)
}
