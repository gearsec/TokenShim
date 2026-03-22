package config

// Item represents a single configuration item.
type Item struct {
	// Key is the configuration key (e.g., "secrets.vault.address").
	Key string
	// DefaultValue is the default value for the configuration item.
	DefaultValue interface{}
	// Required indicates if the configuration item is mandatory.
	Required bool
}
