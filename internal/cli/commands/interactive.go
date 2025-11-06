package commands

import (
	"os"

	"github.com/alenon/gokanon/internal/interactive"
	"github.com/alenon/gokanon/internal/ui"
)

// Interactive starts the interactive mode
func Interactive() error {
	session, err := interactive.New()
	if err != nil {
		return ui.NewError(
			"Failed to start interactive mode",
			err,
			"Ensure your terminal supports interactive input",
			"Try running individual commands instead: 'gokanon help'",
		)
	}

	// Register all command handlers
	session.RegisterCommand("run", func(args []string) error {
		os.Args = append([]string{"gokanon", "run"}, args...)
		return Run()
	})

	session.RegisterCommand("list", func(args []string) error {
		os.Args = append([]string{"gokanon", "list"}, args...)
		return List()
	})

	session.RegisterCommand("compare", func(args []string) error {
		os.Args = append([]string{"gokanon", "compare"}, args...)
		return Compare()
	})

	session.RegisterCommand("export", func(args []string) error {
		os.Args = append([]string{"gokanon", "export"}, args...)
		return Export()
	})

	session.RegisterCommand("stats", func(args []string) error {
		os.Args = append([]string{"gokanon", "stats"}, args...)
		return Stats()
	})

	session.RegisterCommand("trend", func(args []string) error {
		os.Args = append([]string{"gokanon", "trend"}, args...)
		return Trend()
	})

	session.RegisterCommand("check", func(args []string) error {
		os.Args = append([]string{"gokanon", "check"}, args...)
		return Check()
	})

	session.RegisterCommand("flamegraph", func(args []string) error {
		os.Args = append([]string{"gokanon", "flamegraph"}, args...)
		return Flamegraph()
	})

	session.RegisterCommand("serve", func(args []string) error {
		os.Args = append([]string{"gokanon", "serve"}, args...)
		return Serve()
	})

	session.RegisterCommand("delete", func(args []string) error {
		os.Args = append([]string{"gokanon", "delete"}, args...)
		return Delete()
	})

	session.RegisterCommand("doctor", func(args []string) error {
		return Doctor()
	})

	return session.Run()
}
