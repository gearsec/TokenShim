package cli

import (
	"fmt"

	"github.com/gearsec/tokenshim/internal/keyring"
	"github.com/spf13/cobra"
)

var secretProfile string

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage real API credentials securely",
	Long:  `Manage real API credentials securely. Secrets are stored in the OS keyring and never written to disk.`,
}

var secretsSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Store a real API credential",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		err := keyring.Set(secretProfile, key, value)
		if err != nil {
			return fmt.Errorf("failed to save secret to keyring: %w", err)
		}

		fmt.Printf("Secret successfully saved to OS keyring for key: %s (profile: %s)\n", key, secretProfile)
		return nil
	},
}

var secretsGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Retrieve a real API credential (mostly for debugging)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		value, err := keyring.Get(secretProfile, key)
		if err != nil {
			return fmt.Errorf("failed to retrieve secret from keyring: %w", err)
		}

		fmt.Printf("%s\n", value)
		return nil
	},
}

var secretsDeleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete an API credential from the OS keyring",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		err := keyring.Delete(secretProfile, key)
		if err != nil {
			return fmt.Errorf("failed to delete secret from keyring: %w", err)
		}

		fmt.Printf("Secret successfully deleted from OS keyring for key: %s (profile: %s)\n", key, secretProfile)
		return nil
	},
}

func init() {
	secretsCmd.PersistentFlags().StringVarP(&secretProfile, "profile", "p", "default", "Profile to use for the secret")
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsSetCmd)
	secretsCmd.AddCommand(secretsGetCmd)
	secretsCmd.AddCommand(secretsDeleteCmd)
}
