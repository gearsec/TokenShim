package doctor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gearsec/tokenshim/internal/config"
)

// Config holds the doctor scan configuration.
type Config struct {
	ScanPaths []string `mapstructure:"scan_paths"`
}

// defaultScanPaths is the built-in list of files to check when no config file exists.
var defaultScanPaths = []string{
	"~/.env",
	"~/.bashrc",
	"~/.zshrc",
	"~/.bash_profile",
	"~/.profile",
	"~/.config/fish/config.fish",
	".env",
	".env.local",
	".env.development",
	".env.production",
	".env.staging",
}

// DefaultConfig returns the built-in configuration.
func DefaultConfig() Config {
	return Config{ScanPaths: defaultScanPaths}
}

// DefaultConfigPath returns ~/.config/tokenshim/doctor.yaml.
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "doctor.yaml"
	}
	return filepath.Join(home, ".config", "tokenshim", "doctor.yaml")
}

// LoadConfig reads the doctor config at path using the shared config Manager.
// If the file does not exist, DefaultConfig is returned without error.
func LoadConfig(path string) (Config, error) {
	mgr := config.NewManager()
	mgr.Register("", config.Item{
		Key:          "scan_paths",
		DefaultValue: defaultScanPaths,
	})

	var cfg Config
	if err := mgr.LoadFile(path, &cfg); err != nil {
		return Config{}, err
	}

	// Viper may return an empty slice if the key was unset despite the default;
	// fall back to the built-in list in that case.
	if len(cfg.ScanPaths) == 0 {
		cfg.ScanPaths = defaultScanPaths
	}

	return cfg, nil
}

// ExpandPath resolves a leading ~ to the user's home directory.
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path[1:]), nil
}
