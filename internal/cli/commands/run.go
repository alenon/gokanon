package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/runner"
	"github.com/alenon/gokanon/internal/storage"
	"github.com/alenon/gokanon/internal/ui"
)

// Run handles the 'run' subcommand
func Run() error {
	runFlags := flag.NewFlagSet("run", flag.ExitOnError)
	benchFilter := runFlags.String("bench", ".", "Benchmark filter (passed to -bench)")
	packagePath := runFlags.String("pkg", "", "Package path (default: current directory)")
	storageDir := runFlags.String("storage", ".gokanon", "Storage directory for results")
	profileFlag := runFlags.String("profile", "", "Enable profiling: cpu, mem, or cpu,mem")
	verbose := runFlags.Bool("verbose", false, "Show detailed benchmark output")
	cpuFlag := runFlags.String("cpu", "", "CPU list (passed to -cpu)")
	benchtimeFlag := runFlags.String("benchtime", "", "Benchmark time (passed to -benchtime)")
	runFlags.Parse(os.Args[2:])

	ui.PrintHeader("Running Benchmarks")
	fmt.Println()

	// Parse profile options
	var profileOpts *runner.ProfileOptions
	if *profileFlag != "" {
		store := storage.NewStorage(*storageDir)
		profileOpts = &runner.ProfileOptions{
			Storage: store,
		}

		profiles := strings.Split(*profileFlag, ",")
		for _, p := range profiles {
			p = strings.TrimSpace(strings.ToLower(p))
			switch p {
			case "cpu":
				profileOpts.EnableCPU = true
			case "mem", "memory":
				profileOpts.EnableMemory = true
			default:
				return ui.NewError(
					fmt.Sprintf("Unknown profile type: %s", p),
					nil,
					"Valid profile types: cpu, mem",
					"Example: -profile=cpu,mem",
				)
			}
		}

		if profileOpts.EnableCPU || profileOpts.EnableMemory {
			var enabled []string
			if profileOpts.EnableCPU {
				enabled = append(enabled, "CPU")
			}
			if profileOpts.EnableMemory {
				enabled = append(enabled, "Memory")
			}
			ui.PrintInfo("Profiling enabled: %s", strings.Join(enabled, ", "))
		}
	}

	// Run benchmarks
	var spinner *ui.Spinner
	if !*verbose {
		spinner = ui.NewSpinner("Executing benchmarks")
		spinner.Start()
	}

	r := runner.NewRunner(*packagePath, *benchFilter)

	// Set CPU and benchtime flags if provided
	if *cpuFlag != "" {
		r = r.WithCPU(*cpuFlag)
	}
	if *benchtimeFlag != "" {
		r = r.WithBenchtime(*benchtimeFlag)
	}

	// Set up progress callback for non-verbose mode
	if !*verbose {
		progressCallback := func(result models.BenchmarkResult) {
			// Format the message with full benchmark details
			msg := fmt.Sprintf("Completed: Benchmark%s | %s iters | %s | %s | %s allocs",
				result.Name,
				formatIterations(result.Iterations),
				formatNsPerOp(result.NsPerOp),
				formatBytes(result.BytesPerOp),
				formatCount(result.AllocsPerOp),
			)
			spinner.UpdateMessage(msg)
		}
		r = r.WithProgress(progressCallback)
	} else {
		// In verbose mode, show raw output
		r = r.WithVerbose(os.Stdout)
	}

	if profileOpts != nil {
		r = r.WithProfiling(profileOpts)
	}

	run, err := r.Run()

	if spinner != nil {
		spinner.Stop()
	}

	if err != nil {
		return ui.ErrBenchmarkFailed(err)
	}

	// Save results
	ui.PrintInfo("Saving results...")
	store := storage.NewStorage(*storageDir)
	if err := store.Save(run); err != nil {
		return ui.NewError(
			"Failed to save results",
			err,
			"Check file permissions on storage directory",
			"Ensure you have write access to: "+*storageDir,
		)
	}

	// Display results
	fmt.Println()
	ui.PrintSuccess("Benchmarks completed successfully!")
	fmt.Printf("Results saved with ID: %s\n\n", ui.Bold(run.ID))

	ui.PrintSection(ui.ChartEmoji, "Run Information")
	fmt.Printf("  Timestamp:  %s\n", ui.Dim(run.Timestamp.Format(time.RFC3339)))
	fmt.Printf("  Duration:   %s\n", ui.Info(run.Duration.String()))
	fmt.Printf("  Go Version: %s\n", ui.Info(run.GoVersion))

	// Display profile info if available
	if run.CPUProfile != "" || run.MemoryProfile != "" {
		fmt.Printf("\nProfiles:\n")
		if run.CPUProfile != "" {
			fmt.Printf("  CPU: %s\n", run.CPUProfile)
		}
		if run.MemoryProfile != "" {
			fmt.Printf("  Memory: %s\n", run.MemoryProfile)
		}
	}
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Benchmark\tIterations\tns/op\tB/op\tallocs/op")
	fmt.Fprintln(w, "---------\t----------\t-----\t----\t---------")
	for _, result := range run.Results {
		fmt.Fprintf(w, "%s\t%d\t%.2f\t%d\t%d\n",
			result.Name,
			result.Iterations,
			result.NsPerOp,
			result.BytesPerOp,
			result.AllocsPerOp,
		)
	}
	w.Flush()

	// Display profile summary if available
	if run.ProfileSummary != nil {
		displayProfileSummary(run.ProfileSummary)
	}

	fmt.Printf("\nResults saved to: %s\n", *storageDir)

	// Hint about viewing flame graphs
	if run.CPUProfile != "" || run.MemoryProfile != "" {
		fmt.Printf("\nView flame graphs: gokanon flamegraph %s\n", run.ID)
	}

	return nil
}

// displayProfileSummary displays profile analysis summary
func displayProfileSummary(summary *models.ProfileSummary) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("PROFILE ANALYSIS")
	fmt.Println(strings.Repeat("=", 80))

	// CPU Profile Summary
	if len(summary.CPUTopFunctions) > 0 {
		fmt.Printf("\n🔥 CPU Hot Functions (Total samples: %d)\n", summary.TotalCPUSamples)
		fmt.Println(strings.Repeat("-", 80))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Function\tFlat%\tCum%")
		for _, fn := range summary.CPUTopFunctions {
			if len(fn.Name) > 50 {
				fn.Name = fn.Name[:47] + "..."
			}
			fmt.Fprintf(w, "%s\t%.1f%%\t%.1f%%\n",
				fn.Name,
				fn.FlatPercent,
				fn.CumPercent,
			)
		}
		w.Flush()
	}

	// Memory Profile Summary
	if len(summary.MemoryTopFunctions) > 0 {
		fmt.Printf("\n💾 Memory Hot Functions (Total: %s)\n", formatBytes(summary.TotalMemoryBytes))
		fmt.Println(strings.Repeat("-", 80))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Function\tFlat%\tBytes")
		for _, fn := range summary.MemoryTopFunctions {
			if len(fn.Name) > 50 {
				fn.Name = fn.Name[:47] + "..."
			}
			fmt.Fprintf(w, "%s\t%.1f%%\t%s\n",
				fn.Name,
				fn.FlatPercent,
				formatBytes(fn.FlatValue),
			)
		}
		w.Flush()
	}

	// Hot Paths
	if len(summary.HotPaths) > 0 {
		fmt.Println("\n🎯 Hot Execution Paths")
		fmt.Println(strings.Repeat("-", 80))

		for i, path := range summary.HotPaths {
			fmt.Printf("\n%d. %.1f%% of execution time (%d samples)\n",
				i+1, path.Percentage, path.Occurrences)
			fmt.Printf("   %s\n", path.Description)
			fmt.Printf("   Path: %s\n", strings.Join(path.Path, " → "))
		}
	}

	// Memory Leaks
	if len(summary.MemoryLeaks) > 0 {
		fmt.Println("\n⚠️  Potential Memory Issues")
		fmt.Println(strings.Repeat("-", 80))

		for _, leak := range summary.MemoryLeaks {
			severityIcon := "⚠️ "
			switch leak.Severity {
			case "high":
				severityIcon = "🔴"
			case "medium":
				severityIcon = "🟡"
			case "low":
				severityIcon = "🟢"
			}

			fmt.Printf("\n%s %s (%s)\n",
				severityIcon,
				leak.Function,
				leak.Severity,
			)
			fmt.Printf("   Allocations: %d (%s)\n", leak.Allocations, formatBytes(leak.Bytes))
			fmt.Printf("   %s\n", leak.Description)
		}
	}

	// Optimization Suggestions
	if len(summary.Suggestions) > 0 {
		fmt.Println("\n💡 Optimization Suggestions")
		fmt.Println(strings.Repeat("=", 80))

		for i, sug := range summary.Suggestions {
			severityIcon := "💡"
			switch sug.Severity {
			case "high":
				severityIcon = "🔴"
			case "medium":
				severityIcon = "🟡"
			case "low":
				severityIcon = "🟢"
			}

			fmt.Printf("\n%d. %s [%s] %s\n", i+1, severityIcon, strings.ToUpper(sug.Type), sug.Function)
			fmt.Printf("   Issue: %s\n", sug.Issue)
			fmt.Printf("   Suggestion: %s\n", sug.Suggestion)
			if sug.Impact != "" {
				fmt.Printf("   Potential Impact: %s\n", sug.Impact)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

// formatBytes formats bytes in human-readable format
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B/op"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B/op", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB/op", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatIterations formats iteration count in human-readable format
func formatIterations(iters int64) string {
	if iters == 0 {
		return "0"
	}
	if iters < 1000 {
		return fmt.Sprintf("%d", iters)
	}
	if iters < 1000000 {
		return fmt.Sprintf("%.1fK", float64(iters)/1000)
	}
	if iters < 1000000000 {
		return fmt.Sprintf("%.1fM", float64(iters)/1000000)
	}
	return fmt.Sprintf("%.1fB", float64(iters)/1000000000)
}

// formatNsPerOp formats nanoseconds per operation in human-readable format
func formatNsPerOp(ns float64) string {
	if ns == 0 {
		return "0 ns/op"
	}
	if ns < 1000 {
		return fmt.Sprintf("%.2f ns/op", ns)
	}
	if ns < 1000000 {
		return fmt.Sprintf("%.2f µs/op", ns/1000)
	}
	if ns < 1000000000 {
		return fmt.Sprintf("%.2f ms/op", ns/1000000)
	}
	return fmt.Sprintf("%.2f s/op", ns/1000000000)
}

// formatCount formats allocation count
func formatCount(count int64) string {
	if count == 0 {
		return "0"
	}
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}
	if count < 1000000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(count)/1000000)
}
