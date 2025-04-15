package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var setEnvCmd = &cobra.Command{
	Use:   "set-env [profile]",
	Short: "Set AWS environment variables from gopass aws/$PROFILE",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := args[0]
		path := fmt.Sprintf("aws/%s", profile)

		out, err := exec.Command("gopass", "show", path).Output()
		if err != nil {
			return fmt.Errorf("failed to fetch secret from gopass: %w", err)
		}

		// Clean existing AWS_ environment variables
		fmt.Println("unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN AWS_DEFAULT_REGION AWS_PROFILE")

		// Parse each line and export
		lines := bytes.Split(out, []byte{'\n'})
		for _, line := range lines {
			if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
				continue
			}

			parts := bytes.SplitN(line, []byte{'='}, 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(string(parts[0]))
			val := strings.TrimSpace(string(parts[1]))

			// Only export AWS_* keys
			if strings.HasPrefix(key, "AWS_") {
				fmt.Printf("export %s=\"%s\"\n", key, val)
			}
		}

		fmt.Printf("export AWS_PROFILE=\"%s\"\n", profile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setEnvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setEnvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setEnvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
