package config

import (
	"os"
	"testing"
)

type testNestedConfig struct {
	SomeField string `mapstructure:"some_field"`
}

func TestManager_EnvVarLoading(t *testing.T) {
	if err := os.Setenv("TOKENSHIM_TEST_SOME_FIELD", "test_value"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() { _ = os.Unsetenv("TOKENSHIM_TEST_SOME_FIELD") }()

	m := NewManager()

	var cfg struct {
		Test testNestedConfig `mapstructure:"test"`
	}
	m.RegisterStruct(&cfg)

	if err := m.Load(&cfg); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Test.SomeField != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", cfg.Test.SomeField)
	}
}
