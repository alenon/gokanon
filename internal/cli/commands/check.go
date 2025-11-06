package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/compare"
	"github.com/alenon/gokanon/internal/storage"
	"github.com/alenon/gokanon/internal/threshold"
)

// Check handles the 'check' subcommand for CI/CD
func Check() error {
	checkFlags := flag.NewFlagSet("check", flag.ExitOnError)
	storageDir := checkFlags.String("storage", ".gokanon", "Storage directory for results")
	latest := checkFlags.Bool("latest", false, "Check last two runs")
	thresholdPercent := checkFlags.Float64("threshold", 5.0, "Maximum allowed performance degradation (%)")
	checkFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)

	var oldID, newID string

	if *latest {
		runs, err := store.List()
		if err != nil {
			return fmt.Errorf("failed to list results: %w", err)
		}
		if len(runs) < 2 {
			return fmt.Errorf("need at least 2 benchmark runs to check")
		}
		newID = runs[0].ID
		oldID = runs[1].ID
	} else {
		args := checkFlags.Args()
		if len(args) != 2 {
			return fmt.Errorf("usage: gokanon check <old-id> <new-id> OR gokanon check --latest")
		}
		oldID = args[0]
		newID = args[1]
	}

	// Load benchmark runs
	oldRun, err := store.Load(oldID)
	if err != nil {
		return fmt.Errorf("failed to load old run: %w", err)
	}

	newRun, err := store.Load(newID)
	if err != nil {
		return fmt.Errorf("failed to load new run: %w", err)
	}

	// Compare
	comparer := compare.NewComparer()
	comparisons := comparer.Compare(oldRun, newRun)

	if len(comparisons) == 0 {
		return fmt.Errorf("no matching benchmarks found between the two runs")
	}

	// Check thresholds
	checker := threshold.NewChecker(*thresholdPercent)
	result := checker.Check(comparisons)

	// Display result
	fmt.Printf("Threshold Check (max degradation: %.1f%%)\n", *thresholdPercent)
	fmt.Printf("Comparing: %s vs %s\n\n", oldID, newID)
	fmt.Println(threshold.FormatResult(result))

	// Exit with appropriate code for CI/CD
	if !result.Passed {
		os.Exit(1)
	}

	return nil
}
