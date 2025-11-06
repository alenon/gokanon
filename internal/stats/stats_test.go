package stats

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/alenon/gokanon/internal/models"
)

func TestNewAnalyzer(t *testing.T) {
	a := NewAnalyzer()
	if a == nil {
		t.Fatal("Expected non-nil analyzer")
	}
}

func TestCalculateStats(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		name           string
		values         []float64
		expectedMean   float64
		expectedMedian float64
		expectedMin    float64
		expectedMax    float64
	}{
		{
			name:           "simple values",
			values:         []float64{100, 200, 300},
			expectedMean:   200.0,
			expectedMedian: 200.0,
			expectedMin:    100.0,
			expectedMax:    300.0,
		},
		{
			name:           "even count",
			values:         []float64{10, 20, 30, 40},
			expectedMean:   25.0,
			expectedMedian: 25.0,
			expectedMin:    10.0,
			expectedMax:    40.0,
		},
		{
			name:           "single value",
			values:         []float64{42.0},
			expectedMean:   42.0,
			expectedMedian: 42.0,
			expectedMin:    42.0,
			expectedMax:    42.0,
		},
		{
			name:           "identical values",
			values:         []float64{100, 100, 100, 100},
			expectedMean:   100.0,
			expectedMedian: 100.0,
			expectedMin:    100.0,
			expectedMax:    100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := a.calculateStats("Test", tt.values)

			if stats.Mean != tt.expectedMean {
				t.Errorf("Expected mean %f, got %f", tt.expectedMean, stats.Mean)
			}
			if stats.Median != tt.expectedMedian {
				t.Errorf("Expected median %f, got %f", tt.expectedMedian, stats.Median)
			}
			if stats.Min != tt.expectedMin {
				t.Errorf("Expected min %f, got %f", tt.expectedMin, stats.Min)
			}
			if stats.Max != tt.expectedMax {
				t.Errorf("Expected max %f, got %f", tt.expectedMax, stats.Max)
			}
			if stats.Count != len(tt.values) {
				t.Errorf("Expected count %d, got %d", len(tt.values), stats.Count)
			}
		})
	}
}

func TestCalculateStatsStdDev(t *testing.T) {
	a := NewAnalyzer()

	// Test with known standard deviation
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	stats := a.calculateStats("Test", values)

	// Mean should be 5.0
	if stats.Mean != 5.0 {
		t.Errorf("Expected mean 5.0, got %f", stats.Mean)
	}

	// Standard deviation should be 2.0
	expectedStdDev := 2.0
	if math.Abs(stats.StdDev-expectedStdDev) > 0.01 {
		t.Errorf("Expected stddev ~%f, got %f", expectedStdDev, stats.StdDev)
	}

	// Variance should be 4.0
	expectedVariance := 4.0
	if math.Abs(stats.Variance-expectedVariance) > 0.01 {
		t.Errorf("Expected variance ~%f, got %f", expectedVariance, stats.Variance)
	}
}

func TestCalculateStatsEmpty(t *testing.T) {
	a := NewAnalyzer()
	stats := a.calculateStats("Test", []float64{})

	if stats != nil {
		t.Error("Expected nil for empty values")
	}
}

func TestAnalyzeMultiple(t *testing.T) {
	a := NewAnalyzer()

	runs := []models.BenchmarkRun{
		{
			Timestamp: time.Now().Add(-2 * time.Hour),
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkA", NsPerOp: 100.0},
				{Name: "BenchmarkB", NsPerOp: 200.0},
			},
		},
		{
			Timestamp: time.Now().Add(-1 * time.Hour),
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkA", NsPerOp: 110.0},
				{Name: "BenchmarkB", NsPerOp: 210.0},
			},
		},
		{
			Timestamp: time.Now(),
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkA", NsPerOp: 105.0},
				{Name: "BenchmarkB", NsPerOp: 205.0},
			},
		},
	}

	stats := a.AnalyzeMultiple(runs)

	if len(stats) != 2 {
		t.Fatalf("Expected 2 benchmark stats, got %d", len(stats))
	}

	// Check BenchmarkA stats
	if statsA, ok := stats["BenchmarkA"]; ok {
		if statsA.Count != 3 {
			t.Errorf("Expected count 3, got %d", statsA.Count)
		}
		expectedMean := (100.0 + 110.0 + 105.0) / 3.0
		if statsA.Mean != expectedMean {
			t.Errorf("Expected mean %f, got %f", expectedMean, statsA.Mean)
		}
	} else {
		t.Error("Expected BenchmarkA in stats")
	}

	// Check BenchmarkB stats
	if statsB, ok := stats["BenchmarkB"]; ok {
		if statsB.Count != 3 {
			t.Errorf("Expected count 3, got %d", statsB.Count)
		}
		expectedMean := (200.0 + 210.0 + 205.0) / 3.0
		if statsB.Mean != expectedMean {
			t.Errorf("Expected mean %f, got %f", expectedMean, statsB.Mean)
		}
	} else {
		t.Error("Expected BenchmarkB in stats")
	}
}

func TestCoefficientOfVariation(t *testing.T) {
	a := NewAnalyzer()

	// Test with 10% CV
	values := []float64{90, 95, 100, 105, 110}
	stats := a.calculateStats("Test", values)

	// CV should be approximately (stddev/mean) * 100
	expectedCV := (stats.StdDev / stats.Mean) * 100
	if math.Abs(stats.CV-expectedCV) > 0.01 {
		t.Errorf("Expected CV %f, got %f", expectedCV, stats.CV)
	}
}

func TestIsStable(t *testing.T) {
	tests := []struct {
		name      string
		cv        float64
		threshold float64
		expected  bool
	}{
		{"stable - below threshold", 5.0, 10.0, true},
		{"stable - at threshold", 10.0, 10.0, true},
		{"unstable - above threshold", 15.0, 10.0, false},
		{"very stable", 1.0, 10.0, true},
		{"very unstable", 50.0, 10.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &Stats{CV: tt.cv}
			result := stats.IsStable(tt.threshold)
			if result != tt.expected {
				t.Errorf("Expected IsStable=%v for CV=%f and threshold=%f, got %v",
					tt.expected, tt.cv, tt.threshold, result)
			}
		})
	}
}

func TestFormatStats(t *testing.T) {
	stats := &Stats{
		Name:     "BenchmarkTest",
		Count:    10,
		Mean:     100.5,
		Median:   101.0,
		Min:      95.0,
		Max:      105.0,
		StdDev:   3.5,
		Variance: 12.25,
		CV:       3.48,
	}

	formatted := FormatStats(stats)

	expectedContains := []string{
		"BenchmarkTest",
		"10",
		"100.50",
		"101.00",
		"3.50",
		"3.5%",
		"95.00",
		"105.00",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(formatted, expected) {
			t.Errorf("Expected formatted stats to contain %q, got: %s", expected, formatted)
		}
	}
}

func TestLinearRegression(t *testing.T) {
	// Test with perfect linear relationship y = 2x + 1
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{3, 5, 7, 9, 11}

	slope, intercept, rSquared := linearRegression(x, y)

	// Slope should be 2
	if math.Abs(slope-2.0) > 0.01 {
		t.Errorf("Expected slope 2.0, got %f", slope)
	}

	// Intercept should be 1
	if math.Abs(intercept-1.0) > 0.01 {
		t.Errorf("Expected intercept 1.0, got %f", intercept)
	}

	// R-squared should be 1.0 (perfect fit)
	if math.Abs(rSquared-1.0) > 0.01 {
		t.Errorf("Expected R² 1.0, got %f", rSquared)
	}
}

func TestLinearRegressionImperfect(t *testing.T) {
	// Test with imperfect relationship
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2.1, 3.9, 6.2, 7.8, 10.1}

	slope, _, rSquared := linearRegression(x, y)

	// Slope should be approximately 2
	if math.Abs(slope-2.0) > 0.5 {
		t.Errorf("Expected slope ~2.0, got %f", slope)
	}

	// R-squared should be high but not perfect
	if rSquared < 0.9 || rSquared > 1.0 {
		t.Errorf("Expected R² between 0.9 and 1.0, got %f", rSquared)
	}
}

func TestAnalyzeTrend(t *testing.T) {
	a := NewAnalyzer()

	// Create runs with improving trend
	runs := []models.BenchmarkRun{
		{
			Timestamp: time.Now().Add(-3 * time.Hour),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 110.0}},
		},
		{
			Timestamp: time.Now().Add(-2 * time.Hour),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 105.0}},
		},
		{
			Timestamp: time.Now().Add(-1 * time.Hour),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 100.0}},
		},
		{
			Timestamp: time.Now(),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 95.0}},
		},
	}

	trend := a.AnalyzeTrend(runs, "Test")

	if trend == nil {
		t.Fatal("Expected non-nil trend")
	}

	if trend.BenchmarkName != "Test" {
		t.Errorf("Expected BenchmarkName Test, got %s", trend.BenchmarkName)
	}

	// Trend should be improving (negative slope)
	if trend.Direction != "improving" {
		t.Errorf("Expected direction improving, got %s", trend.Direction)
	}

	if trend.TrendLine >= 0 {
		t.Errorf("Expected negative slope for improving trend, got %f", trend.TrendLine)
	}
}

func TestAnalyzeTrendDegrading(t *testing.T) {
	a := NewAnalyzer()

	// Create runs with degrading trend
	runs := []models.BenchmarkRun{
		{
			Timestamp: time.Now().Add(-3 * time.Hour),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 95.0}},
		},
		{
			Timestamp: time.Now().Add(-2 * time.Hour),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 100.0}},
		},
		{
			Timestamp: time.Now().Add(-1 * time.Hour),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 105.0}},
		},
		{
			Timestamp: time.Now(),
			Results:   []models.BenchmarkResult{{Name: "Test", NsPerOp: 110.0}},
		},
	}

	trend := a.AnalyzeTrend(runs, "Test")

	if trend == nil {
		t.Fatal("Expected non-nil trend")
	}

	// Trend should be degrading (positive slope)
	if trend.Direction != "degrading" {
		t.Errorf("Expected direction degrading, got %s", trend.Direction)
	}

	if trend.TrendLine <= 0 {
		t.Errorf("Expected positive slope for degrading trend, got %f", trend.TrendLine)
	}
}

func TestAnalyzeTrendStable(t *testing.T) {
	a := NewAnalyzer()

	// Create runs with stable trend (small variations)
	runs := []models.BenchmarkRun{
		{Results: []models.BenchmarkResult{{Name: "Test", NsPerOp: 100.0}}},
		{Results: []models.BenchmarkResult{{Name: "Test", NsPerOp: 100.5}}},
		{Results: []models.BenchmarkResult{{Name: "Test", NsPerOp: 99.5}}},
		{Results: []models.BenchmarkResult{{Name: "Test", NsPerOp: 100.2}}},
	}

	trend := a.AnalyzeTrend(runs, "Test")

	if trend == nil {
		t.Fatal("Expected non-nil trend")
	}

	// Trend should be stable (slope close to 0)
	if trend.Direction != "stable" {
		t.Errorf("Expected direction stable, got %s", trend.Direction)
	}
}

func TestAnalyzeTrendNotEnoughData(t *testing.T) {
	a := NewAnalyzer()

	// Only one data point
	runs := []models.BenchmarkRun{
		{Results: []models.BenchmarkResult{{Name: "Test", NsPerOp: 100.0}}},
	}

	trend := a.AnalyzeTrend(runs, "Test")

	if trend != nil {
		t.Error("Expected nil trend for insufficient data")
	}
}

func TestAnalyzeTrendNonExistentBenchmark(t *testing.T) {
	a := NewAnalyzer()

	runs := []models.BenchmarkRun{
		{Results: []models.BenchmarkResult{{Name: "BenchmarkA", NsPerOp: 100.0}}},
		{Results: []models.BenchmarkResult{{Name: "BenchmarkA", NsPerOp: 105.0}}},
	}

	trend := a.AnalyzeTrend(runs, "BenchmarkB")

	if trend != nil {
		t.Error("Expected nil trend for non-existent benchmark")
	}
}
