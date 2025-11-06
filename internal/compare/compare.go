package compare

import (
	"fmt"
	"math"

	"github.com/alenon/gokanon/internal/models"
)

// Comparer handles benchmark comparison
type Comparer struct {
	threshold float64 // Threshold percentage to consider "same"
}

// NewComparer creates a new comparer with default threshold
func NewComparer() *Comparer {
	return &Comparer{
		threshold: 5.0, // 5% threshold
	}
}

// Compare compares two benchmark runs and returns comparisons for matching benchmarks
func (c *Comparer) Compare(oldRun, newRun *models.BenchmarkRun) []models.Comparison {
	// Create a map of old results for quick lookup
	oldResults := make(map[string]models.BenchmarkResult)
	for _, result := range oldRun.Results {
		oldResults[result.Name] = result
	}

	var comparisons []models.Comparison

	// Compare each new result with corresponding old result
	for _, newResult := range newRun.Results {
		oldResult, exists := oldResults[newResult.Name]
		if !exists {
			continue // Skip benchmarks that don't exist in old run
		}

		comparison := c.compareResults(oldResult, newResult)
		comparisons = append(comparisons, comparison)
	}

	return comparisons
}

// compareResults compares two individual benchmark results
func (c *Comparer) compareResults(old, new models.BenchmarkResult) models.Comparison {
	delta := new.NsPerOp - old.NsPerOp
	deltaPercent := (delta / old.NsPerOp) * 100

	status := "same"
	if math.Abs(deltaPercent) > c.threshold {
		if deltaPercent < 0 {
			status = "improved" // Lower is better
		} else {
			status = "degraded"
		}
	}

	return models.Comparison{
		Name:         new.Name,
		OldNsPerOp:   old.NsPerOp,
		NewNsPerOp:   new.NsPerOp,
		Delta:        delta,
		DeltaPercent: deltaPercent,
		Status:       status,
	}
}

// FormatComparison formats a comparison for display
func FormatComparison(comp models.Comparison) string {
	statusSymbol := "~"
	switch comp.Status {
	case "improved":
		statusSymbol = "✓"
	case "degraded":
		statusSymbol = "✗"
	}

	return fmt.Sprintf("%s %-40s %12.2f ns/op → %12.2f ns/op (%+.2f%%)",
		statusSymbol,
		comp.Name,
		comp.OldNsPerOp,
		comp.NewNsPerOp,
		comp.DeltaPercent,
	)
}

// Summary provides a summary of the comparison
func Summary(comparisons []models.Comparison) string {
	improved := 0
	degraded := 0
	same := 0

	for _, comp := range comparisons {
		switch comp.Status {
		case "improved":
			improved++
		case "degraded":
			degraded++
		case "same":
			same++
		}
	}

	return fmt.Sprintf("Summary: %d improved, %d degraded, %d unchanged",
		improved, degraded, same)
}
