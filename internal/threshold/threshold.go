package threshold

import (
	"fmt"

	"github.com/alenon/gokanon/internal/models"
)

// Result represents the result of a threshold check
type Result struct {
	Passed       bool
	Failures     []Failure
	TotalChecked int
}

// Failure represents a benchmark that failed the threshold check
type Failure struct {
	BenchmarkName string
	DeltaPercent  float64
	Threshold     float64
	Message       string
}

// Checker handles threshold checking for benchmarks
type Checker struct {
	maxDegradation float64 // Maximum allowed performance degradation (%)
}

// NewChecker creates a new threshold checker
func NewChecker(maxDegradation float64) *Checker {
	return &Checker{
		maxDegradation: maxDegradation,
	}
}

// Check checks if the comparisons meet the threshold requirements
func (c *Checker) Check(comparisons []models.Comparison) *Result {
	result := &Result{
		Passed:       true,
		TotalChecked: len(comparisons),
	}

	for _, comp := range comparisons {
		// Check if performance degraded beyond threshold
		if comp.DeltaPercent > c.maxDegradation {
			result.Passed = false
			result.Failures = append(result.Failures, Failure{
				BenchmarkName: comp.Name,
				DeltaPercent:  comp.DeltaPercent,
				Threshold:     c.maxDegradation,
				Message: fmt.Sprintf(
					"Performance degraded by %.2f%% (threshold: %.2f%%)",
					comp.DeltaPercent,
					c.maxDegradation,
				),
			})
		}
	}

	return result
}

// FormatResult formats the threshold check result for display
func FormatResult(result *Result) string {
	if result.Passed {
		return fmt.Sprintf("✓ All %d benchmarks passed the threshold check", result.TotalChecked)
	}

	output := fmt.Sprintf("✗ %d/%d benchmarks failed the threshold check:\n\n",
		len(result.Failures), result.TotalChecked)

	for _, failure := range result.Failures {
		output += fmt.Sprintf("  • %s: %s\n", failure.BenchmarkName, failure.Message)
	}

	return output
}

// ExitCode returns the appropriate exit code for CI/CD
func (r *Result) ExitCode() int {
	if r.Passed {
		return 0
	}
	return 1
}
