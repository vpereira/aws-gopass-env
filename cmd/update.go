package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vpereira/aws-gopass-env/internal/gopassutils"
)

var (
	updateAccessKey string
	updateSecretKey string
	updateRegion    string
)

var updateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update fields in an existing AWS profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := fmt.Sprintf("aws/%s", name)

		// Step 1: Fetch existing content
		out, err := exec.Command("gopass", "show", path).Output()
		if err != nil {
			return fmt.Errorf("failed to fetch existing profile '%s': %w", name, err)
		}

		// Step 2: Parse current key/values
		kv := make(map[string]string)
		lines := bytes.Split(out, []byte{'\n'})
		for _, line := range lines {
			parts := strings.SplitN(string(line), "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			kv[key] = val
		}

		// Step 3: Apply updates
		if updateAccessKey != "" {
			kv["AWS_ACCESS_KEY_ID"] = updateAccessKey
		}
		if updateSecretKey != "" {
			kv["AWS_SECRET_ACCESS_KEY"] = updateSecretKey
		}
		if updateRegion != "" {
			kv["AWS_DEFAULT_REGION"] = updateRegion
		}

		// Step 4: Reconstruct .env content
		var buffer bytes.Buffer
		for _, key := range []string{
			"AWS_ACCESS_KEY_ID",
			"AWS_SECRET_ACCESS_KEY",
			"AWS_DEFAULT_REGION",
		} {
			if val, ok := kv[key]; ok {
				buffer.WriteString(fmt.Sprintf("%s=%s\n", key, val))
			}
		}

		exists, _ := gopassutils.EntryExists(path)

		if exists {
			return fmt.Errorf("profile '%s' already exists in gopass", name)
		}

		// Step 5: Overwrite the entry
		insertCmd := exec.Command("gopass", "insert", "-m", "-f", path)
		stdin, err := insertCmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to open stdin pipe: %w", err)
		}

		if err := insertCmd.Start(); err != nil {
			return fmt.Errorf("failed to start gopass insert: %w", err)
		}

		_, err = stdin.Write(buffer.Bytes())
		if err != nil {
			return fmt.Errorf("failed to write new content to gopass: %w", err)
		}
		stdin.Close()

		if err := insertCmd.Wait(); err != nil {
			return fmt.Errorf("gopass insert failed: %w", err)
		}

		fmt.Printf("✅ Updated AWS profile '%s'\n", name)
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateAccessKey, "access-key", "", "New AWS Access Key ID")
	updateCmd.Flags().StringVar(&updateSecretKey, "secret-access-key", "", "New AWS Secret Access Key")
	updateCmd.Flags().StringVar(&updateRegion, "region", "", "New AWS Region")
	rootCmd.AddCommand(updateCmd)
}
