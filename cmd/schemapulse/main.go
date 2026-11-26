package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect/postgres"
	"github.com/phantom-mbx1/schemapulse/internal/migration"
	"github.com/phantom-mbx1/schemapulse/internal/reporter"
	"github.com/phantom-mbx1/schemapulse/internal/rules"
)

const version = "0.2.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	var exitCode int
	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Printf("schemapulse %s\n", version)
		exitCode = 0
	case "help", "--help", "-h":
		printUsage()
		exitCode = 0
	case "rules":
		exitCode = runRules()
	case "init":
		exitCode = runInit(os.Args[2:])
	case "inspect":
		exitCode = runInspect(os.Args[2:])
	case "check":
		exitCode = runCheck(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		exitCode = 2
	}

	os.Exit(exitCode)
}

func runRules() int {
	ruleSet := rules.DefaultRules()
	fmt.Printf("Built-in PostgreSQL Rules (%d):\n", len(ruleSet))
	for _, r := range ruleSet {
		fmt.Printf("  - %-6s [%-6s] %s\n", r.ID(), strings.ToUpper(string(r.DefaultSeverity())), r.Title())
	}
	return 0
}

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	outputPath := fs.String("output", "schemapulse.config.json", "Output path for starter config")
	force := fs.Bool("force", false, "Overwrite existing config file if present")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if err := config.WriteStarter(*outputPath, *force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	fmt.Printf("Created starter config at %s\n", *outputPath)
	return 0
}

func runInspect(args []string) int {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	target := "."
	if fs.NArg() > 0 {
		target = fs.Arg(0)
	}

	res, err := migration.Scan(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	fmt.Printf("Migration Inspection — Target: %s (%d file(s))\n", target, len(res.Files))
	for _, f := range res.Files {
		seqLabel := "unnumbered"
		if f.HasSeq {
			seqLabel = fmt.Sprintf("seq=%d", f.Sequence)
		}
		fmt.Printf("  - %s (%s)\n", f.Name, seqLabel)
	}

	if len(res.Issues) > 0 {
		fmt.Printf("\nSequence / File Issues (%d):\n", len(res.Issues))
		for _, issue := range res.Issues {
			fmt.Printf("  [WARN] %s: %s\n", issue.File, issue.Message)
		}
		return 1
	}

	fmt.Println("\nSequence check passed with no issues.")
	return 0
}

func splitCSV(val string) []string {
	if strings.TrimSpace(val) == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	configPath := fs.String("config", "", "Path to SchemaPulse config JSON file")
	format := fs.String("format", "text", "Output format: text, json, markdown")
	failOnFlag := fs.String("fail-on", "", "Override fail threshold: low, medium, high, none")
	onlyFlag := fs.String("only", "", "Comma-separated rule IDs to run (e.g., PG001,PG002)")
	skipFlag := fs.String("skip", "", "Comma-separated rule IDs to skip (e.g., PG003)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	target := "."
	if fs.NArg() > 0 {
		target = fs.Arg(0)
	}

	cfg, _, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	if *failOnFlag != "" {
		cfg.FailOn = *failOnFlag
		if err := cfg.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 2
		}
	}

	rep, err := reporter.New(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	opts := analyzer.Options{
		OnlyRules: splitCSV(*onlyFlag),
		SkipRules: splitCSV(*skipFlag),
	}

	pgAdapter := postgres.New()
	eng := analyzer.NewWithOptions(pgAdapter, rules.DefaultRules(), cfg, opts)
	res, err := eng.Analyze(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	if err := rep.Render(os.Stdout, res); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	if res.ThresholdBreach {
		return 1
	}
	return 0
}

func printUsage() {
	fmt.Print("SchemaPulse - PostgreSQL migration safety & impact analyzer\n\n" +
		"Usage:\n" +
		"  schemapulse <command> [flags] [path]\n\n" +
		"Commands:\n" +
		"  check     Analyze SQL migration files or directories for risks\n" +
		"  inspect   Validate migration file ordering and sequence integrity\n" +
		"  rules     List all built-in migration safety rules\n" +
		"  init      Generate a starter schemapulse.config.json file\n" +
		"  version   Print the SchemaPulse version\n" +
		"  help      Show this help message\n\n" +
		"Flags for 'check':\n" +
		"  --config <path>   Path to config file (auto-discovers schemapulse.config.json)\n" +
		"  --format <fmt>    Output format: text (default), json, markdown\n" +
		"  --fail-on <sev>   Severity threshold to fail with exit code 1: low, medium, high, none\n" +
		"  --only <ids>      Run only specified rule IDs (comma-separated)\n" +
		"  --skip <ids>      Skip specified rule IDs (comma-separated)\n")
}
