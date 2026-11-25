package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

type Config struct {
	Dialect           string            `json:"dialect"`
	FailOn            string            `json:"failOn"`
	DisabledRules     []string          `json:"disabledRules"`
	SeverityOverrides map[string]string `json:"severityOverrides"`
}

func Default() Config {
	return Config{
		Dialect:           "postgres",
		FailOn:            "high",
		DisabledRules:     []string{},
		SeverityOverrides: map[string]string{},
	}
}

func Load(explicitPath string) (Config, string, error) {
	cfg := Default()

	if explicitPath != "" {
		data, err := os.ReadFile(explicitPath)
		if err != nil {
			return cfg, "", fmt.Errorf("failed to read config file %q: %w", explicitPath, err)
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return cfg, "", fmt.Errorf("invalid JSON in config file %q: %w", explicitPath, err)
		}
		if err := cfg.Validate(); err != nil {
			return cfg, "", err
		}
		return cfg, explicitPath, nil
	}

	candidates := []string{
		"schemapulse.config.json",
		".schemapulserc.json",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			data, err := os.ReadFile(candidate)
			if err != nil {
				return cfg, "", fmt.Errorf("failed to read config file %q: %w", candidate, err)
			}
			if err := json.Unmarshal(data, &cfg); err != nil {
				return cfg, "", fmt.Errorf("invalid JSON in config file %q: %w", candidate, err)
			}
			if err := cfg.Validate(); err != nil {
				return cfg, "", err
			}
			return cfg, candidate, nil
		}
	}

	if err := cfg.Validate(); err != nil {
		return cfg, "", err
	}
	return cfg, "", nil
}

func (c Config) Validate() error {
	switch strings.ToLower(strings.TrimSpace(c.FailOn)) {
	case "low", "medium", "high", "none":
	default:
		return fmt.Errorf("invalid failOn value %q (expected low, medium, high, or none)", c.FailOn)
	}

	for ruleID, sev := range c.SeverityOverrides {
		switch strings.ToLower(strings.TrimSpace(sev)) {
		case "low", "medium", "high":
		default:
			return fmt.Errorf("invalid severity override %q for rule %s (expected low, medium, or high)", sev, ruleID)
		}
	}

	return nil
}

func (c Config) IsRuleDisabled(ruleID string) bool {
	target := strings.ToUpper(strings.TrimSpace(ruleID))
	for _, id := range c.DisabledRules {
		if strings.ToUpper(strings.TrimSpace(id)) == target {
			return true
		}
	}
	return false
}

func (c Config) EffectiveSeverity(ruleID string, defaultSev Severity) Severity {
	for k, v := range c.SeverityOverrides {
		if strings.EqualFold(strings.TrimSpace(k), strings.TrimSpace(ruleID)) {
			return Severity(strings.ToLower(strings.TrimSpace(v)))
		}
	}
	return defaultSev
}

func SeverityRank(sev string) (int, bool) {
	switch strings.ToLower(strings.TrimSpace(sev)) {
	case "low":
		return 1, true
	case "medium":
		return 2, true
	case "high":
		return 3, true
	default:
		return 0, false
	}
}

func WriteStarter(targetPath string, force bool) error {
	if !force {
		if _, err := os.Stat(targetPath); err == nil {
			return fmt.Errorf("config file %q already exists (use --force to overwrite)", targetPath)
		}
	}

	if dir := filepath.Dir(targetPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory for config %q: %w", targetPath, err)
		}
	}

	starter := Config{
		Dialect:       "postgres",
		FailOn:        "high",
		DisabledRules: []string{},
		SeverityOverrides: map[string]string{
			"PG003": "high",
		},
	}

	payload, err := json.MarshalIndent(starter, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal starter config: %w", err)
	}
	payload = append(payload, '\n')

	if err := os.WriteFile(targetPath, payload, 0o644); err != nil {
		return fmt.Errorf("failed to write config file %q: %w", targetPath, err)
	}
	return nil
}
