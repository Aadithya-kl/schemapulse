package analyzer

import (
	"fmt"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
	"github.com/phantom-mbx1/schemapulse/internal/migration"
	"github.com/phantom-mbx1/schemapulse/internal/rules"
)

type SeveritySummary struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
}

type Result struct {
	Target          string                    `json:"target"`
	Dialect         string                    `json:"dialect"`
	FilesAnalyzed   int                       `json:"filesAnalyzed"`
	StatementsTotal int                       `json:"statementsTotal"`
	Findings        []rules.Finding           `json:"findings"`
	SequenceIssues  []migration.SequenceIssue `json:"sequenceIssues"`
	Summary         SeveritySummary           `json:"summary"`
	ThresholdBreach bool                      `json:"thresholdBreach"`
	FailOn          string                    `json:"failOn"`
}

type Options struct {
	OnlyRules []string
	SkipRules []string
}

type Analyzer struct {
	adapter dialect.Adapter
	rules   []rules.Rule
	cfg     config.Config
	opts    Options
}

func New(adapter dialect.Adapter, ruleSet []rules.Rule, cfg config.Config) *Analyzer {
	return NewWithOptions(adapter, ruleSet, cfg, Options{})
}

func NewWithOptions(adapter dialect.Adapter, ruleSet []rules.Rule, cfg config.Config, opts Options) *Analyzer {
	return &Analyzer{
		adapter: adapter,
		rules:   ruleSet,
		cfg:     cfg,
		opts:    opts,
	}
}

func (a *Analyzer) shouldRunRule(ruleID string) bool {
	id := strings.ToUpper(strings.TrimSpace(ruleID))
	if a.cfg.IsRuleDisabled(id) {
		return false
	}
	for _, skip := range a.opts.SkipRules {
		if strings.ToUpper(strings.TrimSpace(skip)) == id {
			return false
		}
	}
	if len(a.opts.OnlyRules) > 0 {
		for _, only := range a.opts.OnlyRules {
			if strings.ToUpper(strings.TrimSpace(only)) == id {
				return true
			}
		}
		return false
	}
	return true
}

func (a *Analyzer) Analyze(target string) (*Result, error) {
	scanRes, err := migration.Scan(target)
	if err != nil {
		return nil, err
	}

	res := &Result{
		Target:         target,
		Dialect:        a.adapter.Name(),
		FilesAnalyzed:  len(scanRes.Files),
		Findings:       make([]rules.Finding, 0),
		SequenceIssues: scanRes.Issues,
		FailOn:         a.cfg.FailOn,
	}

	thresholdRank, hasThreshold := config.SeverityRank(a.cfg.FailOn)

	for _, f := range scanRes.Files {
		stmts, err := a.adapter.ParseStatements(f.Path, f.Content)
		if err != nil {
			return nil, fmt.Errorf("parse error in %s: %w", f.Path, err)
		}
		res.StatementsTotal += len(stmts)

		for _, stmt := range stmts {
			for _, r := range a.rules {
				if !a.shouldRunRule(r.ID()) {
					continue
				}
				matches := r.Evaluate(stmt)
				for _, m := range matches {
					m.Severity = a.cfg.EffectiveSeverity(r.ID(), r.DefaultSeverity())
					res.Findings = append(res.Findings, m)

					switch m.Severity {
					case config.SeverityHigh:
						res.Summary.High++
					case config.SeverityMedium:
						res.Summary.Medium++
					case config.SeverityLow:
						res.Summary.Low++
					}

					if hasThreshold {
						fRank, ok := config.SeverityRank(string(m.Severity))
						if ok && fRank >= thresholdRank {
							res.ThresholdBreach = true
						}
					}
				}
			}
		}
	}

	return res, nil
}
