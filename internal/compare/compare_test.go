package compare

import (
	"strings"
	"testing"
	"time"

	"github.com/alenon/gokanon/internal/models"
)

func TestNewComparer(t *testing.T) {
	c := NewComparer()
	if c == nil {
		t.Fatal("Expected non-nil comparer")
	}
	if c.threshold != 5.0 {
		t.Errorf("Expected default threshold 5.0, got %f", c.threshold)
	}
}

func TestCompare(t *testing.T) {
	c := NewComparer()

	oldRun := &models.BenchmarkRun{
		ID:        "old",
		Timestamp: time.Now().Add(-time.Hour),
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkA", NsPerOp: 100.0},
			{Name: "BenchmarkB", NsPerOp: 200.0},
			{Name: "BenchmarkC", NsPerOp: 300.0},
		},
	}

	newRun := &models.BenchmarkRun{
		ID:        "new",
		Timestamp: time.Now(),
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkA", NsPerOp: 90.0},  // 10% improvement
			{Name: "BenchmarkB", NsPerOp: 220.0}, // 10% degradation
			{Name: "BenchmarkC", NsPerOp: 305.0}, // ~1.67% degradation (within threshold)
		},
	}

	comparisons := c.Compare(oldRun, newRun)

	if len(comparisons) != 3 {
		t.Fatalf("Expected 3 comparisons, got %d", len(comparisons))
	}

	// Check BenchmarkA (improved)
	if comparisons[0].Name != "BenchmarkA" {
		t.Errorf("Expected first comparison to be BenchmarkA, got %s", comparisons[0].Name)
	}
	if comparisons[0].Status != "improved" {
		t.Errorf("Expected status improved, got %s", comparisons[0].Status)
	}
	if comparisons[0].DeltaPercent != -10.0 {
		t.Errorf("Expected DeltaPercent -10.0, got %f", comparisons[0].DeltaPercent)
	}

	// Check BenchmarkB (degraded)
	if comparisons[1].Name != "BenchmarkB" {
		t.Errorf("Expected second comparison to be BenchmarkB, got %s", comparisons[1].Name)
	}
	if comparisons[1].Status != "degraded" {
		t.Errorf("Expected status degraded, got %s", comparisons[1].Status)
	}
	if comparisons[1].DeltaPercent != 10.0 {
		t.Errorf("Expected DeltaPercent 10.0, got %f", comparisons[1].DeltaPercent)
	}

	// Check BenchmarkC (same - within threshold)
	if comparisons[2].Name != "BenchmarkC" {
		t.Errorf("Expected third comparison to be BenchmarkC, got %s", comparisons[2].Name)
	}
	if comparisons[2].Status != "same" {
		t.Errorf("Expected status same, got %s", comparisons[2].Status)
	}
}

func TestCompareNoMatches(t *testing.T) {
	c := NewComparer()

	oldRun := &models.BenchmarkRun{
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkA", NsPerOp: 100.0},
		},
	}

	newRun := &models.BenchmarkRun{
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkB", NsPerOp: 200.0},
		},
	}

	comparisons := c.Compare(oldRun, newRun)

	if len(comparisons) != 0 {
		t.Errorf("Expected 0 comparisons, got %d", len(comparisons))
	}
}

func TestCompareResults(t *testing.T) {
	c := NewComparer()

	tests := []struct {
		name           string
		oldNsPerOp     float64
		newNsPerOp     float64
		expectedStatus string
		expectedDelta  float64
	}{
		{"significant improvement", 100.0, 80.0, "improved", -20.0},
		{"significant degradation", 100.0, 120.0, "degraded", 20.0},
		{"minor improvement", 100.0, 97.0, "same", -3.0},
		{"minor degradation", 100.0, 103.0, "same", 3.0},
		{"no change", 100.0, 100.0, "same", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := models.BenchmarkResult{Name: "Test", NsPerOp: tt.oldNsPerOp}
			new := models.BenchmarkResult{Name: "Test", NsPerOp: tt.newNsPerOp}

			comp := c.compareResults(old, new)

			if comp.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, comp.Status)
			}
			if comp.Delta != tt.expectedDelta {
				t.Errorf("Expected delta %f, got %f", tt.expectedDelta, comp.Delta)
			}
			if comp.OldNsPerOp != tt.oldNsPerOp {
				t.Errorf("Expected OldNsPerOp %f, got %f", tt.oldNsPerOp, comp.OldNsPerOp)
			}
			if comp.NewNsPerOp != tt.newNsPerOp {
				t.Errorf("Expected NewNsPerOp %f, got %f", tt.newNsPerOp, comp.NewNsPerOp)
			}
		})
	}
}

func TestFormatComparison(t *testing.T) {
	tests := []struct {
		name             string
		comparison       models.Comparison
		expectedContains []string
	}{
		{
			name: "improved",
			comparison: models.Comparison{
				Name:         "BenchmarkTest",
				OldNsPerOp:   100.0,
				NewNsPerOp:   90.0,
				Delta:        -10.0,
				DeltaPercent: -10.0,
				Status:       "improved",
			},
			expectedContains: []string{"✓", "BenchmarkTest", "100.00", "90.00", "-10.00%"},
		},
		{
			name: "degraded",
			comparison: models.Comparison{
				Name:         "BenchmarkTest",
				OldNsPerOp:   100.0,
				NewNsPerOp:   120.0,
				Delta:        20.0,
				DeltaPercent: 20.0,
				Status:       "degraded",
			},
			expectedContains: []string{"✗", "BenchmarkTest", "100.00", "120.00", "+20.00%"},
		},
		{
			name: "same",
			comparison: models.Comparison{
				Name:         "BenchmarkTest",
				OldNsPerOp:   100.0,
				NewNsPerOp:   102.0,
				Delta:        2.0,
				DeltaPercent: 2.0,
				Status:       "same",
			},
			expectedContains: []string{"~", "BenchmarkTest", "100.00", "102.00", "+2.00%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatComparison(tt.comparison)

			for _, expected := range tt.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, got: %s", expected, result)
				}
			}
		})
	}
}

func TestSummary(t *testing.T) {
	comparisons := []models.Comparison{
		{Status: "improved"},
		{Status: "improved"},
		{Status: "degraded"},
		{Status: "same"},
		{Status: "same"},
		{Status: "same"},
	}

	summary := Summary(comparisons)

	expected := "Summary: 2 improved, 1 degraded, 3 unchanged"
	if summary != expected {
		t.Errorf("Expected summary %q, got %q", expected, summary)
	}
}

func TestSummaryEmpty(t *testing.T) {
	comparisons := []models.Comparison{}

	summary := Summary(comparisons)

	expected := "Summary: 0 improved, 0 degraded, 0 unchanged"
	if summary != expected {
		t.Errorf("Expected summary %q, got %q", expected, summary)
	}
}

func TestSummaryAllImproved(t *testing.T) {
	comparisons := []models.Comparison{
		{Status: "improved"},
		{Status: "improved"},
		{Status: "improved"},
	}

	summary := Summary(comparisons)

	expected := "Summary: 3 improved, 0 degraded, 0 unchanged"
	if summary != expected {
		t.Errorf("Expected summary %q, got %q", expected, summary)
	}
}

func TestSummaryAllDegraded(t *testing.T) {
	comparisons := []models.Comparison{
		{Status: "degraded"},
		{Status: "degraded"},
	}

	summary := Summary(comparisons)

	expected := "Summary: 0 improved, 2 degraded, 0 unchanged"
	if summary != expected {
		t.Errorf("Expected summary %q, got %q", expected, summary)
	}
}
