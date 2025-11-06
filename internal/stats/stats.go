package stats

import (
	"fmt"
	"math"
	"sort"

	"github.com/alenon/gokanon/internal/models"
)

// Stats represents statistical data for a benchmark across multiple runs
type Stats struct {
	Name     string
	Count    int
	Mean     float64
	Median   float64
	Min      float64
	Max      float64
	StdDev   float64
	Variance float64
	CV       float64 // Coefficient of Variation (StdDev/Mean)
}

// Analyzer handles statistical analysis of benchmarks
type Analyzer struct{}

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeMultiple analyzes multiple benchmark runs and provides statistics
func (a *Analyzer) AnalyzeMultiple(runs []models.BenchmarkRun) map[string]*Stats {
	// Group results by benchmark name
	grouped := make(map[string][]float64)

	for _, run := range runs {
		for _, result := range run.Results {
			grouped[result.Name] = append(grouped[result.Name], result.NsPerOp)
		}
	}

	// Calculate statistics for each benchmark
	stats := make(map[string]*Stats)
	for name, values := range grouped {
		stats[name] = a.calculateStats(name, values)
	}

	return stats
}

// calculateStats calculates statistical measures for a set of values
func (a *Analyzer) calculateStats(name string, values []float64) *Stats {
	if len(values) == 0 {
		return nil
	}

	// Sort values for median calculation
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	stats := &Stats{
		Name:  name,
		Count: len(values),
		Min:   sorted[0],
		Max:   sorted[len(sorted)-1],
	}

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	stats.Mean = sum / float64(len(values))

	// Calculate median
	if len(sorted)%2 == 0 {
		stats.Median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	} else {
		stats.Median = sorted[len(sorted)/2]
	}

	// Calculate variance and standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - stats.Mean
		sumSquaredDiff += diff * diff
	}
	stats.Variance = sumSquaredDiff / float64(len(values))
	stats.StdDev = math.Sqrt(stats.Variance)

	// Calculate coefficient of variation
	if stats.Mean != 0 {
		stats.CV = (stats.StdDev / stats.Mean) * 100
	}

	return stats
}

// FormatStats formats statistics for display
func FormatStats(stats *Stats) string {
	return fmt.Sprintf(
		"%-40s Count: %3d | Mean: %10.2f ns/op | Median: %10.2f ns/op | StdDev: %8.2f (±%.1f%%) | Range: [%.2f - %.2f]",
		stats.Name,
		stats.Count,
		stats.Mean,
		stats.Median,
		stats.StdDev,
		stats.CV,
		stats.Min,
		stats.Max,
	)
}

// IsStable returns true if the benchmark is considered stable (low variation)
func (s *Stats) IsStable(threshold float64) bool {
	return s.CV <= threshold
}

// TrendAnalysis represents performance trend over time
type TrendAnalysis struct {
	BenchmarkName string
	Direction     string  // "improving", "degrading", "stable"
	TrendLine     float64 // Slope of the trend line
	Confidence    float64 // R-squared value (0-1)
}

// AnalyzeTrend analyzes the performance trend over time
func (a *Analyzer) AnalyzeTrend(runs []models.BenchmarkRun, benchmarkName string) *TrendAnalysis {
	var values []float64
	var times []float64

	for i, run := range runs {
		for _, result := range run.Results {
			if result.Name == benchmarkName {
				values = append(values, result.NsPerOp)
				times = append(times, float64(i))
				break
			}
		}
	}

	if len(values) < 2 {
		return nil
	}

	// Calculate linear regression
	slope, _, rSquared := linearRegression(times, values)

	direction := "stable"
	if math.Abs(slope) > 1.0 { // Threshold for meaningful change
		if slope < 0 {
			direction = "improving"
		} else {
			direction = "degrading"
		}
	}

	return &TrendAnalysis{
		BenchmarkName: benchmarkName,
		Direction:     direction,
		TrendLine:     slope,
		Confidence:    rSquared,
	}
}

// linearRegression calculates the linear regression for the given data
// Returns: slope, intercept, r-squared
func linearRegression(x, y []float64) (float64, float64, float64) {
	n := float64(len(x))

	var sumX, sumY, sumXY, sumX2, sumY2 float64
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	// Calculate slope and intercept
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// Calculate R-squared
	meanY := sumY / n
	ssRes := 0.0 // Sum of squares of residuals
	ssTot := 0.0 // Total sum of squares

	for i := 0; i < len(x); i++ {
		predicted := slope*x[i] + intercept
		ssRes += (y[i] - predicted) * (y[i] - predicted)
		ssTot += (y[i] - meanY) * (y[i] - meanY)
	}

	rSquared := 1.0
	if ssTot != 0 {
		rSquared = 1.0 - (ssRes / ssTot)
	}

	return slope, intercept, rSquared
}
