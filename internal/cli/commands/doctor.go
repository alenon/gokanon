package commands

import (
	"fmt"

	"github.com/alenon/gokanon/internal/doctor"
)

// Doctor runs diagnostics to check the setup
func Doctor() error {
	results := doctor.RunDiagnostics()
	doctor.PrintResults(results)

	// Return error if any critical checks failed
	for _, result := range results {
		if !result.Passed && (result.Name == "Go Installation" || result.Name == "Go Test Command") {
			return fmt.Errorf("critical check failed: %s", result.Name)
		}
	}

	return nil
}
