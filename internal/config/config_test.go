package config_test

import (
	"path/filepath"
	"testing"

	"github.com/phantom-mbx1/schemapulse/internal/config"
)

func TestWriteAndLoadStarterConfig(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "schemapulse.config.json")
	if err := config.WriteStarter(cfgPath, false); err != nil {
		t.Fatalf("failed to write starter config: %v", err)
	}

	cfg, loadedFrom, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load starter config: %v", err)
	}
	if loadedFrom != cfgPath {
		t.Fatalf("expected loadedFrom=%q, got %q", cfgPath, loadedFrom)
	}
	if !cfg.IsRuleDisabled("") && cfg.EffectiveSeverity("PG003", config.SeverityMedium) != config.SeverityHigh {
		t.Fatalf("expected PG003 severity override to high, got %s", cfg.EffectiveSeverity("PG003", config.SeverityMedium))
	}
}
