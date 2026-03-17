package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage named profiles for agent workloads",
	Long:  `Create and manage named profiles to specify allowed services for specific agent workloads.`,
}

var services []string

var profileCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a named profile specifying allowed services",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("Created profile '%s' with services: %s\n", name, strings.Join(services, ", "))
		// TODO: Save profile configuration to OS/fs
		return nil
	},
}

func init() {
	profileCreateCmd.Flags().StringSliceVarP(&services, "service", "s", []string{}, "Services to allow in this profile")
	profileCmd.AddCommand(profileCreateCmd)
	rootCmd.AddCommand(profileCmd)
}
