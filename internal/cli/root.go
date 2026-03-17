package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tokenshim",
	Short: "A local proxy that injects API credentials seamlessly",
	Long: `tokenshim is a local proxy that keeps real API credentials out of AI agent 
environments by injecting them in-flight during API calls.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
