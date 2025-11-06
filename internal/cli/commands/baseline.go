package commands

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/alenon/gokanon/internal/storage"
	"github.com/alenon/gokanon/internal/ui"
)

// Baseline handles the 'baseline' subcommand
func Baseline() error {
	if len(os.Args) < 3 {
		fmt.Println("Baseline management commands:")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  gokanon baseline <subcommand> [options]")
		fmt.Println()
		fmt.Println("Subcommands:")
		fmt.Println("  save     Save a benchmark run as a baseline")
		fmt.Println("  list     List all saved baselines")
		fmt.Println("  show     Show details of a specific baseline")
		fmt.Println("  delete   Delete a baseline")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  gokanon baseline save -name=v1.0")
		fmt.Println("  gokanon baseline save -name=main -run=run-123 -desc='Main branch baseline'")
		fmt.Println("  gokanon baseline list")
		fmt.Println("  gokanon baseline show -name=v1.0")
		fmt.Println("  gokanon baseline delete -name=v1.0")
		fmt.Println()
		return nil
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "save":
		return baselineSave()
	case "list":
		return baselineList()
	case "show":
		return baselineShow()
	case "delete":
		return baselineDelete()
	default:
		return ui.NewError(
			fmt.Sprintf("Unknown baseline subcommand: %s", subcommand),
			nil,
			"Valid subcommands: save, list, show, delete",
			"Run 'gokanon baseline' to see usage",
		)
	}
}

// baselineSave saves a benchmark run as a baseline
func baselineSave() error {
	saveFlags := flag.NewFlagSet("baseline-save", flag.ExitOnError)
	name := saveFlags.String("name", "", "Baseline name (required)")
	runID := saveFlags.String("run", "", "Run ID to save as baseline (default: latest run)")
	description := saveFlags.String("desc", "", "Baseline description")
	storageDir := saveFlags.String("storage", ".gokanon", "Storage directory for results")
	saveFlags.Parse(os.Args[3:])

	if *name == "" {
		return ui.NewError(
			"Baseline name is required",
			nil,
			"Use -name flag to specify baseline name",
			"Example: gokanon baseline save -name=v1.0",
		)
	}

	store := storage.NewStorage(*storageDir)

	// Determine which run to use
	var targetRunID string
	if *runID != "" {
		targetRunID = *runID
	} else {
		// Use latest run
		run, err := store.GetLatest()
		if err != nil {
			return ui.NewError(
				"Failed to get latest run",
				err,
				"No benchmark runs found",
				"Run 'gokanon run' first to create a benchmark run",
			)
		}
		targetRunID = run.ID
	}

	// Save baseline
	ui.PrintInfo("Saving baseline '%s' from run %s...", *name, targetRunID)
	baseline, err := store.SaveBaseline(*name, targetRunID, *description, nil)
	if err != nil {
		return ui.NewError(
			"Failed to save baseline",
			err,
			"Check that the run ID exists and storage directory is writable",
			"Try: gokanon list",
		)
	}

	fmt.Println()
	ui.PrintSuccess("Baseline saved successfully!")
	fmt.Printf("Name:        %s\n", ui.Bold(baseline.Name))
	fmt.Printf("Run ID:      %s\n", baseline.RunID)
	fmt.Printf("Created:     %s\n", baseline.CreatedAt.Format(time.RFC3339))
	if baseline.Description != "" {
		fmt.Printf("Description: %s\n", baseline.Description)
	}
	fmt.Printf("Benchmarks:  %d\n", len(baseline.Run.Results))
	fmt.Println()
	fmt.Printf("Baseline saved to: %s/baselines/%s.json\n", *storageDir, baseline.Name)
	return nil
}

// baselineList lists all saved baselines
func baselineList() error {
	listFlags := flag.NewFlagSet("baseline-list", flag.ExitOnError)
	storageDir := listFlags.String("storage", ".gokanon", "Storage directory for results")
	listFlags.Parse(os.Args[3:])

	store := storage.NewStorage(*storageDir)
	baselines, err := store.ListBaselines()
	if err != nil {
		return ui.NewError(
			"Failed to list baselines",
			err,
			"Check storage directory permissions",
		)
	}

	if len(baselines) == 0 {
		fmt.Println("No baselines found.")
		fmt.Println()
		fmt.Println("Create a baseline with: gokanon baseline save -name=<name>")
		return nil
	}

	ui.PrintHeader("Saved Baselines")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tCreated\tBenchmarks\tDescription")
	fmt.Fprintln(w, "----\t-------\t----------\t-----------")

	for _, baseline := range baselines {
		desc := baseline.Description
		if desc == "" {
			desc = "-"
		}
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			baseline.Name,
			baseline.CreatedAt.Format("2006-01-02 15:04"),
			len(baseline.Run.Results),
			desc,
		)
	}
	w.Flush()
	fmt.Println()

	return nil
}

// baselineShow shows details of a specific baseline
func baselineShow() error {
	showFlags := flag.NewFlagSet("baseline-show", flag.ExitOnError)
	name := showFlags.String("name", "", "Baseline name (required)")
	storageDir := showFlags.String("storage", ".gokanon", "Storage directory for results")
	showFlags.Parse(os.Args[3:])

	if *name == "" {
		return ui.NewError(
			"Baseline name is required",
			nil,
			"Use -name flag to specify baseline name",
			"Example: gokanon baseline show -name=v1.0",
		)
	}

	store := storage.NewStorage(*storageDir)
	baseline, err := store.LoadBaseline(*name)
	if err != nil {
		return ui.NewError(
			fmt.Sprintf("Failed to load baseline '%s'", *name),
			err,
			"Check that the baseline exists",
			"Try: gokanon baseline list",
		)
	}

	ui.PrintHeader(fmt.Sprintf("Baseline: %s", baseline.Name))
	fmt.Println()

	fmt.Printf("Name:        %s\n", ui.Bold(baseline.Name))
	fmt.Printf("Run ID:      %s\n", baseline.RunID)
	fmt.Printf("Created:     %s\n", baseline.CreatedAt.Format(time.RFC3339))
	if baseline.Description != "" {
		fmt.Printf("Description: %s\n", baseline.Description)
	}
	fmt.Println()

	ui.PrintSection(ui.ChartEmoji, "Run Information")
	fmt.Printf("  Timestamp:  %s\n", baseline.Run.Timestamp.Format(time.RFC3339))
	fmt.Printf("  Duration:   %s\n", baseline.Run.Duration.String())
	fmt.Printf("  Go Version: %s\n", baseline.Run.GoVersion)
	fmt.Printf("  Package:    %s\n", baseline.Run.Package)
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Benchmark\tIterations\tns/op\tB/op\tallocs/op")
	fmt.Fprintln(w, "---------\t----------\t-----\t----\t---------")
	for _, result := range baseline.Run.Results {
		fmt.Fprintf(w, "%s\t%d\t%.2f\t%d\t%d\n",
			result.Name,
			result.Iterations,
			result.NsPerOp,
			result.BytesPerOp,
			result.AllocsPerOp,
		)
	}
	w.Flush()

	return nil
}

// baselineDelete deletes a baseline
func baselineDelete() error {
	deleteFlags := flag.NewFlagSet("baseline-delete", flag.ExitOnError)
	name := deleteFlags.String("name", "", "Baseline name (required)")
	storageDir := deleteFlags.String("storage", ".gokanon", "Storage directory for results")
	deleteFlags.Parse(os.Args[3:])

	if *name == "" {
		return ui.NewError(
			"Baseline name is required",
			nil,
			"Use -name flag to specify baseline name",
			"Example: gokanon baseline delete -name=v1.0",
		)
	}

	store := storage.NewStorage(*storageDir)

	// Check if baseline exists
	if !store.HasBaseline(*name) {
		return ui.NewError(
			fmt.Sprintf("Baseline '%s' not found", *name),
			nil,
			"Check that the baseline exists",
			"Try: gokanon baseline list",
		)
	}

	// Delete baseline
	if err := store.DeleteBaseline(*name); err != nil {
		return ui.NewError(
			"Failed to delete baseline",
			err,
			"Check storage directory permissions",
		)
	}

	ui.PrintSuccess("Baseline '%s' deleted successfully", *name)
	return nil
}
