package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/dashboard"
	"github.com/alenon/gokanon/internal/storage"
)

// Serve starts the interactive web dashboard
func Serve() error {
	serveFlags := flag.NewFlagSet("serve", flag.ExitOnError)
	storageDir := serveFlags.String("storage", ".gokanon", "Storage directory for results")
	port := serveFlags.Int("port", 8080, "Port for web server")
	addr := serveFlags.String("addr", "localhost", "Address to bind to (use 0.0.0.0 for all interfaces)")
	serveFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)

	// Check if storage directory exists
	if _, err := os.Stat(*storageDir); os.IsNotExist(err) {
		fmt.Printf("Warning: Storage directory '%s' does not exist.\n", *storageDir)
		fmt.Println("Run some benchmarks first with: gokanon run")
		fmt.Println("\nStarting dashboard anyway...")
	}

	// Create and start the dashboard server
	server := dashboard.NewServer(store, *addr, *port)

	fmt.Println("Starting interactive web dashboard...")
	fmt.Printf("Dashboard will be available at: http://%s:%d\n", *addr, *port)
	fmt.Println("\nPress Ctrl+C to stop the server")

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start dashboard server: %w", err)
	}

	return nil
}
