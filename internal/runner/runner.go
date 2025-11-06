package runner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alenon/gokanon/internal/aianalyzer"
	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/profiler"
	"github.com/alenon/gokanon/internal/storage"
)

// ProfileOptions configures profiling behavior
type ProfileOptions struct {
	EnableCPU    bool
	EnableMemory bool
	Storage      *storage.Storage
}

// Runner handles benchmark execution
type Runner struct {
	packagePath    string
	benchFilter    string
	profileOptions *ProfileOptions
}

// NewRunner creates a new benchmark runner
func NewRunner(packagePath, benchFilter string) *Runner {
	return &Runner{
		packagePath: packagePath,
		benchFilter: benchFilter,
	}
}

// WithProfiling configures the runner to enable profiling
func (r *Runner) WithProfiling(opts *ProfileOptions) *Runner {
	r.profileOptions = opts
	return r
}

// Run executes the benchmarks and returns parsed results
func (r *Runner) Run() (*models.BenchmarkRun, error) {
	startTime := time.Now()

	// Get Go version
	goVersion, err := r.getGoVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get Go version: %w", err)
	}

	// Generate unique ID for this run
	runID := generateID()

	// Create temporary directory for profile files
	tempDir, err := os.MkdirTemp("", "gokanon-profile-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Build the benchmark command
	args := []string{"test", "-bench", r.benchFilter, "-benchmem"}

	// Add profiling flags if enabled
	var cpuProfilePath, memProfilePath string
	if r.profileOptions != nil {
		if r.profileOptions.EnableCPU {
			cpuProfilePath = filepath.Join(tempDir, "cpu.prof")
			args = append(args, "-cpuprofile", cpuProfilePath)
		}
		if r.profileOptions.EnableMemory {
			memProfilePath = filepath.Join(tempDir, "mem.prof")
			args = append(args, "-memprofile", memProfilePath)
		}
	}

	if r.packagePath != "" {
		args = append(args, r.packagePath)
	} else {
		args = append(args, "./...")
	}

	// Execute benchmark
	cmd := exec.Command("go", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("benchmark execution failed: %w\nStderr: %s", err, stderr.String())
	}

	// Parse results
	results, err := r.parseOutput(stdout.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse benchmark output: %w", err)
	}

	duration := time.Since(startTime)

	run := &models.BenchmarkRun{
		ID:        runID,
		Timestamp: startTime,
		Package:   r.packagePath,
		GoVersion: goVersion,
		Results:   results,
		Command:   fmt.Sprintf("go %s", strings.Join(args, " ")),
		Duration:  duration,
	}

	// Handle profile files if profiling was enabled
	if r.profileOptions != nil && r.profileOptions.Storage != nil {
		if err := r.handleProfiles(run, cpuProfilePath, memProfilePath); err != nil {
			// Log warning but don't fail the run
			fmt.Fprintf(os.Stderr, "Warning: failed to process profiles: %v\n", err)
		}
	}

	return run, nil
}

// parseOutput parses the benchmark output from go test -bench
func (r *Runner) parseOutput(output string) ([]models.BenchmarkResult, error) {
	var results []models.BenchmarkResult

	// Regex to match benchmark lines
	// Example: BenchmarkFoo-8   1000000   1234 ns/op   512 B/op   10 allocs/op
	benchRegex := regexp.MustCompile(`^Benchmark(\S+)\s+(\d+)\s+([\d.]+)\s+ns/op(?:\s+([\d.]+)\s+MB/s)?(?:\s+(\d+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	// Increase buffer size to handle long output lines (default is 64KB, set to 1MB)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB max token size
	for scanner.Scan() {
		line := scanner.Text()
		matches := benchRegex.FindStringSubmatch(line)

		if matches != nil {
			name := matches[1]
			iterations, _ := strconv.ParseInt(matches[2], 10, 64)
			nsPerOp, _ := strconv.ParseFloat(matches[3], 64)

			result := models.BenchmarkResult{
				Name:       name,
				Iterations: iterations,
				NsPerOp:    nsPerOp,
			}

			// Parse optional MB/s
			if matches[4] != "" {
				result.MBPerSec, _ = strconv.ParseFloat(matches[4], 64)
			}

			// Parse optional B/op
			if matches[5] != "" {
				result.BytesPerOp, _ = strconv.ParseInt(matches[5], 10, 64)
			}

			// Parse optional allocs/op
			if matches[6] != "" {
				result.AllocsPerOp, _ = strconv.ParseInt(matches[6], 10, 64)
			}

			results = append(results, result)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no benchmark results found in output")
	}

	return results, nil
}

// getGoVersion returns the current Go version
func (r *Runner) getGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// generateID generates a unique ID for a benchmark run
func generateID() string {
	return fmt.Sprintf("run-%d", time.Now().Unix())
}

// handleProfiles processes and stores profile files, and analyzes them
func (r *Runner) handleProfiles(run *models.BenchmarkRun, cpuProfilePath, memProfilePath string) error {
	store := r.profileOptions.Storage
	analyzer := profiler.NewAnalyzer()

	// Process CPU profile
	if cpuProfilePath != "" {
		if _, err := os.Stat(cpuProfilePath); err == nil {
			// Read profile data
			cpuData, err := os.ReadFile(cpuProfilePath)
			if err != nil {
				return fmt.Errorf("failed to read CPU profile: %w", err)
			}

			// Save to storage
			if err := store.SaveProfile(run.ID, "cpu", bytes.NewReader(cpuData)); err != nil {
				return fmt.Errorf("failed to save CPU profile: %w", err)
			}

			// Set profile path in run
			run.CPUProfile = store.GetCPUProfilePath(run.ID)

			// Load into analyzer
			if err := analyzer.LoadCPUProfile(cpuData); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to analyze CPU profile: %v\n", err)
			}
		}
	}

	// Process memory profile
	if memProfilePath != "" {
		if _, err := os.Stat(memProfilePath); err == nil {
			// Read profile data
			memData, err := os.ReadFile(memProfilePath)
			if err != nil {
				return fmt.Errorf("failed to read memory profile: %w", err)
			}

			// Save to storage
			if err := store.SaveProfile(run.ID, "memory", bytes.NewReader(memData)); err != nil {
				return fmt.Errorf("failed to save memory profile: %w", err)
			}

			// Set profile path in run
			run.MemoryProfile = store.GetMemoryProfilePath(run.ID)

			// Load into analyzer
			if err := analyzer.LoadMemoryProfile(memData); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to analyze memory profile: %v\n", err)
			}
		}
	}

	// Analyze profiles and generate summary
	if run.CPUProfile != "" || run.MemoryProfile != "" {
		summary, err := analyzer.Analyze()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to analyze profiles: %v\n", err)
		} else {
			// Enhance with AI analysis if enabled
			aiAnalyzer, err := aianalyzer.NewFromEnv()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize AI analyzer: %v\n", err)
				run.ProfileSummary = summary
			} else {
				enhanced, err := aiAnalyzer.EnhanceProfileSummary(summary)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: AI analysis failed: %v\n", err)
					run.ProfileSummary = summary
				} else {
					run.ProfileSummary = enhanced
				}
			}
		}
	}

	return nil
}
