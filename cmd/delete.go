package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an AWS profile from the gopass aws store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := fmt.Sprintf("aws/%s", name)

		cmdExec := exec.Command("gopass", "rm", "-f", path)
		output, err := cmdExec.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to delete '%s': %v\n%s", name, err, string(output))
		}

		fmt.Printf("✅ Deleted AWS profile '%s' from gopass store.\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
