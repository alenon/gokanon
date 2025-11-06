package cli

import (
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/cli/commands"
	"github.com/alenon/gokanon/internal/ui"
)

const (
	usageText = `gokanon - A CLI tool for running and comparing Go benchmarks

Usage:
  gokanon <command> [options]

Commands:
  run          Run benchmarks and save results
  list         List all saved benchmark results
  compare      Compare two benchmark results
  export       Export comparison results to various formats
  stats        Show statistical analysis of multiple runs
  trend        Analyze performance trends over time
  check        Check performance against thresholds (for CI/CD)
  flamegraph   View CPU/memory flame graphs for a run
  serve        Start interactive web dashboard
  delete       Delete a benchmark result
  baseline     Manage baseline benchmarks (save, load, list, show, delete)
  doctor       Run diagnostics to check your setup
  interactive  Start interactive mode with auto-completion
  completion   Install shell completion scripts
  help         Show this help message

Examples:
  gokanon run                            # Run all benchmarks in current package
  gokanon run -bench=. -pkg=./...        # Run all benchmarks in all packages
  gokanon run -bench=BenchmarkFoo        # Run specific benchmark
  gokanon run -profile=cpu,mem           # Run with CPU and memory profiling
  gokanon list                           # List all saved results
  gokanon compare run-123 run-456        # Compare two specific runs
  gokanon compare --latest               # Compare last two runs
  gokanon compare --baseline=v1.0        # Compare latest run with baseline
  gokanon export --latest -format=html   # Export comparison to HTML
  gokanon stats -last=5                  # Show stats for last 5 runs
  gokanon trend -last=10                 # Show performance trends
  gokanon check --latest -threshold=10   # Check if degradation > 10%
  gokanon flamegraph run-123             # View flame graphs in browser
  gokanon serve                          # Start interactive web dashboard
  gokanon serve -port=9000               # Start dashboard on custom port
  gokanon delete run-123                 # Delete a specific run
  gokanon baseline save -name=v1.0       # Save latest run as baseline
  gokanon baseline save -name=v1.0 -run=run-123  # Save specific run as baseline
  gokanon baseline list                  # List all saved baselines
  gokanon baseline show -name=v1.0       # Show baseline details
  gokanon baseline delete -name=v1.0     # Delete a baseline
  gokanon doctor                         # Check your setup
  gokanon interactive                    # Start interactive mode
  gokanon completion bash                # Install bash completion

For more information about a command, use:
  gokanon <command> -h
`
)

// Execute is the main entry point for the CLI
func Execute() error {
	if len(os.Args) < 2 {
		fmt.Print(usageText)
		return nil
	}

	command := os.Args[1]

	switch command {
	case "run":
		return commands.Run()
	case "list":
		return commands.List()
	case "compare":
		return commands.Compare()
	case "export":
		return commands.Export()
	case "stats":
		return commands.Stats()
	case "trend":
		return commands.Trend()
	case "check":
		return commands.Check()
	case "flamegraph":
		return commands.Flamegraph()
	case "serve":
		return commands.Serve()
	case "delete":
		return commands.Delete()
	case "baseline":
		return commands.Baseline()
	case "doctor":
		return commands.Doctor()
	case "interactive", "i":
		return commands.Interactive()
	case "completion":
		return commands.Completion()
	case "help", "-h", "--help":
		fmt.Print(usageText)
		return nil
	default:
		return ui.NewError(
			fmt.Sprintf("Unknown command: %s", command),
			nil,
			"Run 'gokanon help' to see available commands",
			"Use 'gokanon <command> -h' for command-specific help",
			"Try 'gokanon interactive' for an interactive experience",
		)
	}
}
