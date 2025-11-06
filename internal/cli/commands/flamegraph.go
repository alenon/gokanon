package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/storage"
	"github.com/alenon/gokanon/internal/webserver"
)

// Flamegraph handles the 'flamegraph' subcommand
func Flamegraph() error {
	flamegraphFlags := flag.NewFlagSet("flamegraph", flag.ExitOnError)
	storageDir := flamegraphFlags.String("storage", ".gokanon", "Storage directory for results")
	port := flamegraphFlags.String("port", "8080", "Port for web server")
	latest := flamegraphFlags.Bool("latest", false, "View profiles for latest run")
	flamegraphFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)

	var runID string

	if *latest {
		// Get the most recent run
		latestRun, err := store.GetLatest()
		if err != nil {
			return fmt.Errorf("failed to get latest run: %w", err)
		}
		runID = latestRun.ID
	} else {
		// Get run ID from arguments
		args := flamegraphFlags.Args()
		if len(args) != 1 {
			return fmt.Errorf("usage: gokanon flamegraph <run-id> OR gokanon flamegraph --latest")
		}
		runID = args[0]
	}

	// Load the run to verify it has profiles
	run, err := store.Load(runID)
	if err != nil {
		return fmt.Errorf("failed to load run: %w", err)
	}

	if run.CPUProfile == "" && run.MemoryProfile == "" {
		return fmt.Errorf("no profiles found for run %s\n\nRun benchmarks with profiling enabled:\n  gokanon run --profile=cpu,mem", runID)
	}

	fmt.Printf("Starting flame graph viewer for run: %s\n", runID)
	if run.CPUProfile != "" {
		fmt.Println("  ✓ CPU profile available")
	}
	if run.MemoryProfile != "" {
		fmt.Println("  ✓ Memory profile available")
	}
	fmt.Println()

	// Start web server
	server := webserver.NewServer(store, *port)
	return server.Start(runID)
}
