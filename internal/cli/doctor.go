package cli

import (
	"fmt"
	"os"

	"github.com/gearsec/tokenshim/internal/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Scan for exposed secrets in files and environment variables",
	Long:  `Doctor scans your system for exposed API keys and credentials in files and environment variables, and produces a report.`,
}

var (
	doctorCheckFiles   bool
	doctorCheckEnv     bool
	doctorOutputFormat string
	doctorExportPath   string
	doctorConfigPath   string
)

var doctorCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Run secret detection scans",
	Long: `Scan files and/or environment variables for exposed secrets.

By default both --files and --env are run. Use the flags to restrict the scan.

Supported output formats: json (default), yaml, xml, csv, html.`,
	RunE: runDoctorCheck,
}

func runDoctorCheck(cmd *cobra.Command, _ []string) (err error) {
	if err := validateDoctorOutputFormat(doctorOutputFormat); err != nil {
		return err
	}

	// Resolve config path.
	configPath := doctorConfigPath
	if configPath == "" {
		configPath = doctor.DefaultConfigPath()
	}

	cfg, err := doctor.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load doctor config: %w", err)
	}

	// If neither flag is set, run both.
	runFiles := doctorCheckFiles
	runEnv := doctorCheckEnv
	if !runFiles && !runEnv {
		runFiles = true
		runEnv = true
	}

	var (
		fileFindings      []doctor.FileFinding
		filesRead         []string
		filesWithFindings int
		envFindings       []doctor.EnvFinding
		varsScanned       int
		varsWithFindings  int
	)

	if runFiles {
		fileFindings, filesRead, filesWithFindings, err = doctor.ScanFiles(cfg)
		if err != nil {
			return fmt.Errorf("failed to scan files: %w", err)
		}
	}

	if runEnv {
		envFindings, varsScanned, varsWithFindings = doctor.ScanEnv()
	}

	report := doctor.NewReport(
		fileFindings, filesRead, filesWithFindings,
		envFindings, varsScanned, varsWithFindings,
	)

	// Determine output destination.
	out := cmd.OutOrStdout()
	if doctorExportPath != "" {
		f, createErr := os.Create(doctorExportPath)
		if createErr != nil {
			return fmt.Errorf("failed to create export file: %w", createErr)
		}
		defer func() {
			if cerr := f.Close(); cerr != nil && err == nil {
				err = cerr
			}
		}()
		out = f
	}

	if err = doctor.WriteReport(out, report, doctorOutputFormat); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	if doctorExportPath != "" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Report written to %s\n", doctorExportPath)
	}

	return nil
}

func validateDoctorOutputFormat(format string) error {
	switch format {
	case "json", "yaml", "xml", "csv", "html":
		return nil
	}
	return fmt.Errorf("invalid output format %q: must be one of json, yaml, xml, csv, html", format)
}

func init() {
	doctorCheckCmd.Flags().BoolVar(&doctorCheckFiles, "files", false, "Scan files listed in doctor config")
	doctorCheckCmd.Flags().BoolVar(&doctorCheckEnv, "env", false, "Scan environment variables for exposed secrets")
	doctorCheckCmd.Flags().StringVarP(&doctorOutputFormat, "output", "o", "json", "Export format: json|yaml|xml|csv|html")
	doctorCheckCmd.Flags().StringVarP(&doctorExportPath, "export", "e", "", "Write report to this file path (default: stdout)")
	doctorCheckCmd.Flags().StringVarP(&doctorConfigPath, "config", "c", "", "Path to doctor config (default: ~/.config/tokenshim/doctor.yaml)")

	doctorCmd.AddCommand(doctorCheckCmd)
	rootCmd.AddCommand(doctorCmd)
}
