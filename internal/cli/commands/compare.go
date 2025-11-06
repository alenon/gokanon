package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/aianalyzer"
	"github.com/alenon/gokanon/internal/compare"
	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/storage"
	"github.com/alenon/gokanon/internal/ui"
)

// Compare handles the 'compare' subcommand
func Compare() error {
	compareFlags := flag.NewFlagSet("compare", flag.ExitOnError)
	storageDir := compareFlags.String("storage", ".gokanon", "Storage directory for results")
	latest := compareFlags.Bool("latest", false, "Compare the last two runs")
	baseline := compareFlags.String("baseline", "", "Compare latest run against a baseline")
	compareFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)

	var oldID, newID string
	var oldRun, newRun *models.BenchmarkRun

	if *baseline != "" {
		// Compare latest run against baseline
		baselineData, err := store.LoadBaseline(*baseline)
		if err != nil {
			return ui.NewError(
				fmt.Sprintf("Failed to load baseline '%s'", *baseline),
				err,
				"Check that the baseline exists",
				"Try: gokanon baseline list",
			)
		}

		latestRun, err := store.GetLatest()
		if err != nil {
			return ui.NewError(
				"Failed to get latest run",
				err,
				"No benchmark runs found",
				"Run 'gokanon run' first",
			)
		}

		oldRun = baselineData.Run
		newRun = latestRun
		oldID = baselineData.Name + " (baseline)"
		newID = latestRun.ID
	} else if *latest {
		// Get the two most recent runs
		runs, err := store.List()
		if err != nil {
			return fmt.Errorf("failed to list results: %w", err)
		}
		if len(runs) < 2 {
			return fmt.Errorf("need at least 2 benchmark runs to compare")
		}
		newID = runs[0].ID
		oldID = runs[1].ID
	} else {
		// Get IDs from arguments
		args := compareFlags.Args()
		if len(args) != 2 {
			return fmt.Errorf("usage: gokanon compare <old-id> <new-id> OR gokanon compare --latest OR gokanon compare --baseline=<name>")
		}
		oldID = args[0]
		newID = args[1]
	}

	// Load benchmark runs if not already loaded
	if oldRun == nil {
		var err error
		oldRun, err = store.Load(oldID)
		if err != nil {
			return fmt.Errorf("failed to load old run: %w", err)
		}
	}

	if newRun == nil {
		var err error
		newRun, err = store.Load(newID)
		if err != nil {
			return fmt.Errorf("failed to load new run: %w", err)
		}
	}

	// Compare
	comparer := compare.NewComparer()
	comparisons := comparer.Compare(oldRun, newRun)

	if len(comparisons) == 0 {
		fmt.Println("No matching benchmarks found between the two runs.")
		return nil
	}

	// Display comparison
	fmt.Printf("Comparing: %s (%s) vs %s (%s)\n\n",
		oldID, oldRun.Timestamp.Format("2006-01-02 15:04:05"),
		newID, newRun.Timestamp.Format("2006-01-02 15:04:05"),
	)

	for _, comp := range comparisons {
		fmt.Println(compare.FormatComparison(comp))
	}

	fmt.Printf("\n%s\n", compare.Summary(comparisons))

	// Add AI analysis if enabled
	aiAnalyzer, err := aianalyzer.NewFromEnv()
	if err == nil {
		analysis, err := aiAnalyzer.AnalyzeComparison(oldRun, newRun, comparisons)
		if err == nil && analysis != "" {
			fmt.Printf("\n--- AI Analysis ---\n%s\n", analysis)
		}
	}

	return nil
}
