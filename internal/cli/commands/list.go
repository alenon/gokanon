package commands

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/alenon/gokanon/internal/storage"
)

// List handles the 'list' subcommand
func List() error {
	listFlags := flag.NewFlagSet("list", flag.ExitOnError)
	storageDir := listFlags.String("storage", ".gokanon", "Storage directory for results")
	listFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)
	runs, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to list results: %w", err)
	}

	if len(runs) == 0 {
		fmt.Println("No benchmark results found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTimestamp\tBenchmarks\tDuration\tPackage")
	fmt.Fprintln(w, "--\t---------\t----------\t--------\t-------")

	for _, run := range runs {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
			run.ID,
			run.Timestamp.Format("2006-01-02 15:04:05"),
			len(run.Results),
			run.Duration,
			run.Package,
		)
	}
	w.Flush()

	return nil
}
