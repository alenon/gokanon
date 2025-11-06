package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/stats"
	"github.com/alenon/gokanon/internal/storage"
)

// Trend handles the 'trend' subcommand
func Trend() error {
	trendFlags := flag.NewFlagSet("trend", flag.ExitOnError)
	storageDir := trendFlags.String("storage", ".gokanon", "Storage directory for results")
	lastN := trendFlags.Int("last", 10, "Analyze last N runs")
	benchmark := trendFlags.String("benchmark", "", "Specific benchmark to analyze (empty = all)")
	trendFlags.Parse(os.Args[2:])

	store := storage.NewStorage(*storageDir)
	runs, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to list results: %w", err)
	}

	if len(runs) < 2 {
		return fmt.Errorf("need at least 2 benchmark runs for trend analysis")
	}

	// Limit to last N runs
	if *lastN > 0 && *lastN < len(runs) {
		runs = runs[:*lastN]
	}

	// Reverse to get chronological order
	for i := 0; i < len(runs)/2; i++ {
		runs[i], runs[len(runs)-1-i] = runs[len(runs)-1-i], runs[i]
	}

	fmt.Printf("Performance Trend Analysis (%d runs)\n", len(runs))
	fmt.Printf("Period: %s to %s\n\n",
		runs[0].Timestamp.Format("2006-01-02 15:04:05"),
		runs[len(runs)-1].Timestamp.Format("2006-01-02 15:04:05"),
	)

	analyzer := stats.NewAnalyzer()

	// Get all unique benchmark names
	benchmarkNames := make(map[string]bool)
	for _, run := range runs {
		for _, result := range run.Results {
			if *benchmark == "" || result.Name == *benchmark {
				benchmarkNames[result.Name] = true
			}
		}
	}

	// Analyze trend for each benchmark
	for name := range benchmarkNames {
		trend := analyzer.AnalyzeTrend(runs, name)
		if trend == nil {
			continue
		}

		fmt.Printf("Benchmark: %s\n", name)

		// Show direction
		directionSymbol := "→"
		directionColor := ""
		switch trend.Direction {
		case "improving":
			directionSymbol = "↓"
			directionColor = "🟢"
		case "degrading":
			directionSymbol = "↑"
			directionColor = "🔴"
		default:
			directionColor = "⚪"
		}

		fmt.Printf("  %s Trend: %s %s (slope: %.2f ns/op per run)\n",
			directionColor,
			trend.Direction,
			directionSymbol,
			trend.TrendLine,
		)

		fmt.Printf("  Confidence: %.1f%% (R²)\n", trend.Confidence*100)

		// Show data points
		fmt.Printf("  Data points: ")
		var values []float64
		for _, run := range runs {
			for _, result := range run.Results {
				if result.Name == name {
					values = append(values, result.NsPerOp)
					break
				}
			}
		}

		// Show sparkline-like representation
		if len(values) > 0 {
			min, max := values[0], values[0]
			for _, v := range values {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}

			fmt.Printf("%.2f", values[0])
			for i := 1; i < len(values); i++ {
				change := ((values[i] - values[i-1]) / values[i-1]) * 100
				if change > 0 {
					fmt.Printf(" → %.2f (+%.1f%%)", values[i], change)
				} else {
					fmt.Printf(" → %.2f (%.1f%%)", values[i], change)
				}
			}
			fmt.Println()
		}

		fmt.Println()
	}

	return nil
}
