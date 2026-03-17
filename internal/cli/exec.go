package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var profileName string
var dryRun bool

var execCmd = &cobra.Command{
	Use:   "exec [--profile <name>] [--dry-run] -- <command>",
	Short: "Run a command with injected masked credentials and proxy settings",
	Long:  `Executes the specified command with a masked environment. The agent receives a masked token and a local proxy address.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Executing command: %s\n", strings.Join(args, " "))

		if profileName != "" {
			fmt.Printf("Using profile: %s\n", profileName)
		}

		if dryRun {
			fmt.Println("Dry run mode enabled: Will only print the masked environment")
		} else {
			// TODO: Start local proxy, set env vars, and run the target command
		}

		return nil
	},
}

func init() {
	execCmd.Flags().StringVarP(&profileName, "profile", "p", "", "Named profile to use")
	execCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print the masked environment that will be injected")
	rootCmd.AddCommand(execCmd)
}
