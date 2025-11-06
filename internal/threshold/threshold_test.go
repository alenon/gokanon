package threshold

import (
	"strings"
	"testing"

	"github.com/alenon/gokanon/internal/models"
)

func TestNewChecker(t *testing.T) {
	tests := []struct {
		name              string
		maxDegradation    float64
		expectedThreshold float64
	}{
		{"default", 5.0, 5.0},
		{"strict", 1.0, 1.0},
		{"lenient", 20.0, 20.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChecker(tt.maxDegradation)
			if c == nil {
				t.Fatal("Expected non-nil checker")
			}
			if c.maxDegradation != tt.expectedThreshold {
				t.Errorf("Expected threshold %f, got %f", tt.expectedThreshold, c.maxDegradation)
			}
		})
	}
}

func TestCheckPassed(t *testing.T) {
	c := NewChecker(10.0)

	comparisons := []models.Comparison{
		{Name: "BenchmarkA", DeltaPercent: -5.0, Status: "improved"}, // Improvement always passes
		{Name: "BenchmarkB", DeltaPercent: 5.0, Status: "degraded"},  // Within threshold
		{Name: "BenchmarkC", DeltaPercent: 2.0, Status: "same"},      // Within threshold
	}

	result := c.Check(comparisons)

	if !result.Passed {
		t.Error("Expected check to pass")
	}
	if len(result.Failures) != 0 {
		t.Errorf("Expected 0 failures, got %d", len(result.Failures))
	}
	if result.TotalChecked != 3 {
		t.Errorf("Expected TotalChecked 3, got %d", result.TotalChecked)
	}
}

func TestCheckFailed(t *testing.T) {
	c := NewChecker(5.0)

	comparisons := []models.Comparison{
		{Name: "BenchmarkA", DeltaPercent: -2.0, Status: "improved"}, // OK
		{Name: "BenchmarkB", DeltaPercent: 10.0, Status: "degraded"}, // FAIL - exceeds threshold
		{Name: "BenchmarkC", DeltaPercent: 3.0, Status: "same"},      // OK
		{Name: "BenchmarkD", DeltaPercent: 7.5, Status: "degraded"},  // FAIL - exceeds threshold
	}

	result := c.Check(comparisons)

	if result.Passed {
		t.Error("Expected check to fail")
	}
	if len(result.Failures) != 2 {
		t.Errorf("Expected 2 failures, got %d", len(result.Failures))
	}
	if result.TotalChecked != 4 {
		t.Errorf("Expected TotalChecked 4, got %d", result.TotalChecked)
	}

	// Check failure details
	if result.Failures[0].BenchmarkName != "BenchmarkB" {
		t.Errorf("Expected first failure to be BenchmarkB, got %s", result.Failures[0].BenchmarkName)
	}
	if result.Failures[0].DeltaPercent != 10.0 {
		t.Errorf("Expected first failure DeltaPercent 10.0, got %f", result.Failures[0].DeltaPercent)
	}
	if result.Failures[0].Threshold != 5.0 {
		t.Errorf("Expected threshold 5.0, got %f", result.Failures[0].Threshold)
	}

	if result.Failures[1].BenchmarkName != "BenchmarkD" {
		t.Errorf("Expected second failure to be BenchmarkD, got %s", result.Failures[1].BenchmarkName)
	}
}

func TestCheckEmpty(t *testing.T) {
	c := NewChecker(5.0)

	comparisons := []models.Comparison{}

	result := c.Check(comparisons)

	if !result.Passed {
		t.Error("Expected empty check to pass")
	}
	if len(result.Failures) != 0 {
		t.Errorf("Expected 0 failures, got %d", len(result.Failures))
	}
	if result.TotalChecked != 0 {
		t.Errorf("Expected TotalChecked 0, got %d", result.TotalChecked)
	}
}

func TestFormatResultPassed(t *testing.T) {
	result := &Result{
		Passed:       true,
		TotalChecked: 10,
		Failures:     []Failure{},
	}

	formatted := FormatResult(result)

	expectedContains := []string{"✓", "All", "10", "passed"}
	for _, expected := range expectedContains {
		if !strings.Contains(formatted, expected) {
			t.Errorf("Expected result to contain %q, got: %s", expected, formatted)
		}
	}
}

func TestFormatResultFailed(t *testing.T) {
	result := &Result{
		Passed:       false,
		TotalChecked: 10,
		Failures: []Failure{
			{
				BenchmarkName: "BenchmarkA",
				DeltaPercent:  15.0,
				Threshold:     10.0,
				Message:       "Performance degraded by 15.00% (threshold: 10.00%)",
			},
			{
				BenchmarkName: "BenchmarkB",
				DeltaPercent:  12.5,
				Threshold:     10.0,
				Message:       "Performance degraded by 12.50% (threshold: 10.00%)",
			},
		},
	}

	formatted := FormatResult(result)

	expectedContains := []string{
		"✗",
		"2/10",
		"failed",
		"BenchmarkA",
		"BenchmarkB",
		"15.00%",
		"12.50%",
		"10.00%",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(formatted, expected) {
			t.Errorf("Expected result to contain %q, got: %s", expected, formatted)
		}
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name         string
		passed       bool
		expectedCode int
	}{
		{"passed", true, 0},
		{"failed", false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &Result{Passed: tt.passed}
			code := result.ExitCode()
			if code != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectedCode, code)
			}
		})
	}
}

func TestFailureMessage(t *testing.T) {
	failure := Failure{
		BenchmarkName: "BenchmarkTest",
		DeltaPercent:  15.5,
		Threshold:     10.0,
		Message:       "Performance degraded by 15.50% (threshold: 10.00%)",
	}

	if failure.BenchmarkName != "BenchmarkTest" {
		t.Errorf("Expected BenchmarkName BenchmarkTest, got %s", failure.BenchmarkName)
	}
	if failure.DeltaPercent != 15.5 {
		t.Errorf("Expected DeltaPercent 15.5, got %f", failure.DeltaPercent)
	}
	if failure.Threshold != 10.0 {
		t.Errorf("Expected Threshold 10.0, got %f", failure.Threshold)
	}
	if !strings.Contains(failure.Message, "15.50%") {
		t.Errorf("Expected message to contain '15.50%%', got: %s", failure.Message)
	}
}

func TestCheckBoundaryValues(t *testing.T) {
	c := NewChecker(10.0)

	tests := []struct {
		name         string
		deltaPercent float64
		shouldPass   bool
	}{
		{"exactly at threshold", 10.0, true}, // At threshold passes (uses > not >=)
		{"just below threshold", 9.99, true},
		{"just above threshold", 10.01, false},
		{"zero degradation", 0.0, true},
		{"negative (improvement)", -5.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparisons := []models.Comparison{
				{Name: "Test", DeltaPercent: tt.deltaPercent},
			}

			result := c.Check(comparisons)

			if result.Passed != tt.shouldPass {
				t.Errorf("Expected Passed=%v for deltaPercent=%f, got %v",
					tt.shouldPass, tt.deltaPercent, result.Passed)
			}
		})
	}
}
