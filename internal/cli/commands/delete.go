package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/storage"
)

// Delete handles the 'delete' subcommand
func Delete() error {
	deleteFlags := flag.NewFlagSet("delete", flag.ExitOnError)
	storageDir := deleteFlags.String("storage", ".gokanon", "Storage directory for results")
	deleteFlags.Parse(os.Args[2:])

	args := deleteFlags.Args()
	if len(args) != 1 {
		return fmt.Errorf("usage: gokanon delete <id>")
	}

	id := args[0]
	store := storage.NewStorage(*storageDir)

	if err := store.Delete(id); err != nil {
		return fmt.Errorf("failed to delete run: %w", err)
	}

	fmt.Printf("Deleted benchmark run: %s\n", id)
	return nil
}
