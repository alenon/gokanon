package ui

import (
	"fmt"
	"os"
	"strings"
)

// ErrorWithSuggestion represents an error with helpful suggestions
type ErrorWithSuggestion struct {
	Message     string
	Suggestions []string
	Err         error
}

// Error implements the error interface
func (e *ErrorWithSuggestion) Error() string {
	var b strings.Builder
	b.WriteString(Error(ErrorIcon + " " + e.Message))

	if e.Err != nil {
		b.WriteString("\n")
		b.WriteString(Dim("  Cause: " + e.Err.Error()))
	}

	if len(e.Suggestions) > 0 {
		b.WriteString("\n\n")
		b.WriteString(Info("💡 Suggestions:"))
		for _, suggestion := range e.Suggestions {
			b.WriteString("\n  " + ArrowIcon + " " + suggestion)
		}
	}

	return b.String()
}

// Unwrap returns the underlying error
func (e *ErrorWithSuggestion) Unwrap() error {
	return e.Err
}

// NewError creates a new error with suggestions
func NewError(message string, err error, suggestions ...string) *ErrorWithSuggestion {
	return &ErrorWithSuggestion{
		Message:     message,
		Suggestions: suggestions,
		Err:         err,
	}
}

// Common error scenarios with suggestions

// ErrNoResults returns an error when no benchmark results are found
func ErrNoResults() error {
	return NewError(
		"No benchmark results found",
		nil,
		"Run 'gokanon run' to create your first benchmark result",
		"Check if .gokanon directory exists in your project",
		"Use 'gokanon doctor' to diagnose any issues",
	)
}

// ErrInvalidRunID returns an error for invalid run IDs
func ErrInvalidRunID(id string, availableIDs []string) error {
	suggestions := []string{
		"Use 'gokanon list' to see all available run IDs",
	}

	if len(availableIDs) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Available runs: %s", strings.Join(availableIDs, ", ")))
	}

	return NewError(
		fmt.Sprintf("Invalid run ID: %s", id),
		nil,
		suggestions...,
	)
}

// ErrBenchmarkFailed returns an error when benchmark execution fails
func ErrBenchmarkFailed(err error) error {
	return NewError(
		"Benchmark execution failed",
		err,
		"Ensure your test files contain valid benchmark functions (func BenchmarkXxx(b *testing.B))",
		"Check that your package compiles: 'go test -c'",
		"Try running with verbose output: 'gokanon run -bench=. -v'",
		"Verify Go toolchain is installed: 'go version'",
	)
}

// ErrInvalidThreshold returns an error for invalid threshold values
func ErrInvalidThreshold(value string) error {
	return NewError(
		fmt.Sprintf("Invalid threshold value: %s", value),
		nil,
		"Threshold must be a positive number (e.g., 10 for 10%)",
		"Use decimal values for fractional percentages (e.g., 2.5 for 2.5%)",
		"Example: 'gokanon check --latest -threshold=10'",
	)
}

// ErrStorageCorrupted returns an error when storage is corrupted
func ErrStorageCorrupted(err error) error {
	return NewError(
		"Benchmark storage appears to be corrupted",
		err,
		"Try running 'gokanon doctor' to diagnose the issue",
		"Backup and remove the .gokanon directory to start fresh",
		"Check file permissions on .gokanon directory",
	)
}

// ErrProfileNotFound returns an error when profile data is missing
func ErrProfileNotFound(runID string) error {
	return NewError(
		fmt.Sprintf("Profile data not found for run: %s", runID),
		nil,
		"Run benchmarks with profiling enabled: 'gokanon run -profile=cpu,mem'",
		"Profiles are only available for runs with -profile flag",
		"Use 'gokanon list' to see which runs have profiles",
	)
}

// ErrInvalidFormat returns an error for unsupported export formats
func ErrInvalidFormat(format string) error {
	return NewError(
		fmt.Sprintf("Unsupported export format: %s", format),
		nil,
		"Supported formats: html, csv, markdown, json",
		"Example: 'gokanon export --latest -format=html'",
	)
}

// ErrPortInUse returns an error when a port is already in use
func ErrPortInUse(port int, err error) error {
	return NewError(
		fmt.Sprintf("Port %d is already in use", port),
		err,
		fmt.Sprintf("Try a different port: 'gokanon serve -port=%d'", port+1),
		"Check running processes: 'lsof -i :%d'",
		"Kill the process using the port or choose a different one",
	)
}

// PrintErrorAndExit prints an error with suggestions and exits
func PrintErrorAndExit(err error, exitCode int) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err.Error())
	fmt.Fprintln(os.Stderr)
	os.Exit(exitCode)
}
