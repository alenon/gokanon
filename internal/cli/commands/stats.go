package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alenon/gokanon/internal/stats"
	"github.com/alenon/gokanon/internal/storage"
)

// Stats handles the 'stats' subcommand
func Stats() error {
	statsFlags := flag.NewFlagSet("stats", flag.ExitOnError)
	storageDir := statsFlags.String("storage", ".gokanon", "Storage directory for results")
	lastN := statsFlags.Int("last", 0, "Analyze last N runs (0 = all)")
	cvThreshold := statsFlags.Float64("cv-threshold", 10.0, "Coefficient of variation threshold for stability (%)")
	statsFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)
	runs, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to list results: %w", err)
	}

	if len(runs) == 0 {
		return fmt.Errorf("no benchmark results found")
	}

	// Limit to last N runs if specified
	if *lastN > 0 && *lastN < len(runs) {
		runs = runs[:*lastN]
	}

	fmt.Printf("Statistical Analysis (%d runs)\n", len(runs))
	fmt.Printf("Runs: %s to %s\n\n",
		runs[len(runs)-1].Timestamp.Format("2006-01-02 15:04:05"),
		runs[0].Timestamp.Format("2006-01-02 15:04:05"),
	)

	// Analyze
	analyzer := stats.NewAnalyzer()
	statistics := analyzer.AnalyzeMultiple(runs)

	// Display
	fmt.Println("Benchmark Statistics:")
	fmt.Println(strings.Repeat("-", 150))

	for _, stat := range statistics {
		fmt.Println(stats.FormatStats(stat))

		// Show stability indicator
		if stat.IsStable(*cvThreshold) {
			fmt.Print(" ✓ Stable")
		} else {
			fmt.Print(" ⚠ Variable")
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 150))
	fmt.Printf("\nNote: Benchmarks with CV (coefficient of variation) <= %.1f%% are considered stable.\n", *cvThreshold)

	return nil
}
