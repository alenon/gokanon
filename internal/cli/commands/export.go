package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/compare"
	"github.com/alenon/gokanon/internal/export"
	"github.com/alenon/gokanon/internal/storage"
)

// Export handles the 'export' subcommand
func Export() error {
	exportFlags := flag.NewFlagSet("export", flag.ExitOnError)
	storageDir := exportFlags.String("storage", ".gokanon", "Storage directory for results")
	latest := exportFlags.Bool("latest", false, "Export comparison of last two runs")
	format := exportFlags.String("format", "html", "Export format: html, csv, markdown")
	output := exportFlags.String("output", "", "Output file (default: comparison.<format>)")
	exportFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)

	var oldID, newID string

	if *latest {
		runs, err := store.List()
		if err != nil {
			return fmt.Errorf("failed to list results: %w", err)
		}
		if len(runs) < 2 {
			return fmt.Errorf("need at least 2 benchmark runs to export")
		}
		newID = runs[0].ID
		oldID = runs[1].ID
	} else {
		args := exportFlags.Args()
		if len(args) != 2 {
			return fmt.Errorf("usage: gokanon export <old-id> <new-id> OR gokanon export --latest")
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

	// Determine output filename
	outputFile := *output
	if outputFile == "" {
		outputFile = fmt.Sprintf("comparison.%s", *format)
	}

	// Export
	exporter := export.NewExporter()
	switch *format {
	case "html":
		err = exporter.ToHTML(
			comparisons,
			oldID, newID,
			oldRun.Timestamp.Format("2006-01-02 15:04:05"),
			newRun.Timestamp.Format("2006-01-02 15:04:05"),
			outputFile,
		)
	case "csv":
		err = exporter.ToCSV(comparisons, outputFile)
	case "markdown", "md":
		err = exporter.ToMarkdown(comparisons, oldID, newID, outputFile)
	default:
		return fmt.Errorf("unsupported format: %s (supported: html, csv, markdown)", *format)
	}

	if err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	fmt.Printf("Comparison exported to: %s\n", outputFile)
	return nil
}
