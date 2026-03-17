package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View audit logs of proxied requests",
	Long:  `View a complete record of what the agent called and when.`,
}

var sessionName string

var logShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show logs for a specific session",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Showing logs for session: %s\n", sessionName)
		// TODO: Fetch and format logs
		return nil
	},
}

func init() {
	logShowCmd.Flags().StringVar(&sessionName, "session", "latest", "Session name or 'latest'")
	logCmd.AddCommand(logShowCmd)
	rootCmd.AddCommand(logCmd)
}
