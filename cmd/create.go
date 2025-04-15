/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	accessKey string
	secretKey string
	region    string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new AWS profile in the gopass store",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("you must specify a profile name")
		}
		name := args[0]
		fullPath := fmt.Sprintf("aws/%s", name)

		content := fmt.Sprintf(
			"AWS_ACCESS_KEY_ID=%s\nAWS_SECRET_ACCESS_KEY=%s\nAWS_DEFAULT_REGION=%s\n",
			accessKey,
			secretKey,
			region,
		)

		gopass := exec.Command("gopass", "insert", "-m", fullPath)
		stdin, err := gopass.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to open stdin pipe: %w", err)
		}

		if err := gopass.Start(); err != nil {
			return fmt.Errorf("failed to start gopass: %w", err)
		}

		_, err = stdin.Write([]byte(content))
		if err != nil {
			return fmt.Errorf("failed to write to gopass: %w", err)
		}
		stdin.Close()

		if err := gopass.Wait(); err != nil {
			return fmt.Errorf("gopass insert failed: %w", err)
		}

		fmt.Printf("✅ Created AWS profile '%s' in gopass store.\n", name)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&accessKey, "access-key", "", "AWS Access Key ID")
	createCmd.Flags().StringVar(&secretKey, "secret-access-key", "", "AWS Secret Access Key")
	createCmd.Flags().StringVar(&region, "region", "", "AWS Region")
	createCmd.MarkFlagRequired("access-key")
	createCmd.MarkFlagRequired("secret-access-key")
	createCmd.MarkFlagRequired("region")
	rootCmd.AddCommand(createCmd)
}
