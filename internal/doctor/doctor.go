package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alenon/gokanon/internal/storage"
	"github.com/alenon/gokanon/internal/ui"
)

// CheckResult represents the result of a diagnostic check
type CheckResult struct {
	Name        string
	Passed      bool
	Message     string
	Suggestions []string
}

// RunDiagnostics runs all diagnostic checks
func RunDiagnostics() []CheckResult {
	results := []CheckResult{}

	ui.PrintHeader("Running gokanon diagnostics...")
	fmt.Println()

	// Check 1: Go installation
	results = append(results, checkGoInstallation())

	// Check 2: Go test command
	results = append(results, checkGoTest())

	// Check 3: Storage directory
	results = append(results, checkStorageDirectory())

	// Check 4: Storage integrity
	results = append(results, checkStorageIntegrity())

	// Check 5: Benchmark files
	results = append(results, checkBenchmarkFiles())

	// Check 6: Git repository (optional)
	results = append(results, checkGitRepo())

	// Check 7: Available memory
	results = append(results, checkSystemResources())

	return results
}

// PrintResults prints the diagnostic results
func PrintResults(results []CheckResult) {
	fmt.Println()
	ui.PrintSection(ui.ChartEmoji, "Diagnostic Results")
	fmt.Println()

	passedCount := 0
	failedCount := 0

	for _, result := range results {
		if result.Passed {
			passedCount++
			ui.PrintSuccess("%s: %s", result.Name, result.Message)
		} else {
			failedCount++
			ui.PrintError("%s: %s", result.Name, result.Message)
			if len(result.Suggestions) > 0 {
				for _, suggestion := range result.Suggestions {
					fmt.Printf("  %s %s\n", ui.Info(ui.ArrowIcon), suggestion)
				}
			}
		}
		fmt.Println()
	}

	fmt.Println()
	ui.PrintHeader("Summary")
	fmt.Printf("%s %d checks passed\n", ui.Success(ui.CheckEmoji), passedCount)
	if failedCount > 0 {
		fmt.Printf("%s %d checks failed\n", ui.Error(ui.CrossEmoji), failedCount)
	}
	fmt.Println()

	if failedCount == 0 {
		ui.PrintSuccess("All checks passed! Your gokanon setup is healthy.")
	} else {
		ui.PrintWarning("Some checks failed. Please review the suggestions above.")
	}
}

func checkGoInstallation() CheckResult {
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return CheckResult{
			Name:    "Go Installation",
			Passed:  false,
			Message: "Go is not installed or not in PATH",
			Suggestions: []string{
				"Install Go from https://golang.org/dl/",
				"Ensure Go is in your PATH environment variable",
			},
		}
	}

	version := strings.TrimSpace(string(output))
	return CheckResult{
		Name:    "Go Installation",
		Passed:  true,
		Message: version,
	}
}

func checkGoTest() CheckResult {
	cmd := exec.Command("go", "test", "-bench=.", "-run=^$", "-count=1")
	cmd.Env = append(os.Environ(), "GOKANON_DRY_RUN=1")

	// Don't actually run tests, just check if the command is available
	err := cmd.Start()
	if err != nil {
		return CheckResult{
			Name:    "Go Test Command",
			Passed:  false,
			Message: "Cannot execute 'go test' command",
			Suggestions: []string{
				"Ensure Go toolchain is properly installed",
				"Check if current directory is a valid Go module",
			},
		}
	}
	cmd.Process.Kill()

	return CheckResult{
		Name:    "Go Test Command",
		Passed:  true,
		Message: "'go test' command is available",
	}
}

func checkStorageDirectory() CheckResult {
	storageDir := ".gokanon"

	info, err := os.Stat(storageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return CheckResult{
				Name:    "Storage Directory",
				Passed:  true,
				Message: fmt.Sprintf("Storage directory will be created at: %s", storageDir),
			}
		}

		return CheckResult{
			Name:    "Storage Directory",
			Passed:  false,
			Message: fmt.Sprintf("Cannot access storage directory: %v", err),
			Suggestions: []string{
				"Check file permissions on " + storageDir,
				"Ensure parent directory exists and is writable",
			},
		}
	}

	if !info.IsDir() {
		return CheckResult{
			Name:    "Storage Directory",
			Passed:  false,
			Message: fmt.Sprintf("%s exists but is not a directory", storageDir),
			Suggestions: []string{
				"Remove the file: rm " + storageDir,
				"Gobench will create the directory automatically",
			},
		}
	}

	return CheckResult{
		Name:    "Storage Directory",
		Passed:  true,
		Message: fmt.Sprintf("Storage directory exists at: %s", storageDir),
	}
}

func checkStorageIntegrity() CheckResult {
	store := storage.NewStorage(".gokanon")
	runs, err := store.List()
	if err != nil {
		return CheckResult{
			Name:    "Storage Integrity",
			Passed:  false,
			Message: fmt.Sprintf("Cannot read storage: %v", err),
			Suggestions: []string{
				"Run 'gokanon list' to see detailed error",
				"Backup and remove .gokanon directory if corrupted",
			},
		}
	}

	if len(runs) == 0 {
		return CheckResult{
			Name:    "Storage Integrity",
			Passed:  true,
			Message: "No benchmark runs stored yet",
		}
	}

	// Check if we can load a recent run
	_, err = store.Load(runs[0].ID)
	if err != nil {
		return CheckResult{
			Name:    "Storage Integrity",
			Passed:  false,
			Message: fmt.Sprintf("Storage may be corrupted: %v", err),
			Suggestions: []string{
				"Try deleting corrupted runs manually",
				"Backup and recreate .gokanon directory if needed",
			},
		}
	}

	return CheckResult{
		Name:    "Storage Integrity",
		Passed:  true,
		Message: fmt.Sprintf("Storage is healthy with %d run(s)", len(runs)),
	}
}

func checkBenchmarkFiles() CheckResult {
	// Look for *_test.go files with Benchmark functions
	matches, err := filepath.Glob("*_test.go")
	if err != nil {
		return CheckResult{
			Name:    "Benchmark Files",
			Passed:  false,
			Message: fmt.Sprintf("Error searching for test files: %v", err),
		}
	}

	if len(matches) == 0 {
		return CheckResult{
			Name:    "Benchmark Files",
			Passed:  false,
			Message: "No test files found in current directory",
			Suggestions: []string{
				"Create benchmark files with *_test.go naming",
				"Benchmark functions should start with 'Benchmark' (e.g., BenchmarkMyFunc)",
				"Run 'gokanon run -pkg=./...' to search subdirectories",
			},
		}
	}

	// Check if any file contains "Benchmark" functions
	hasBenchmarks := false
	for _, file := range matches {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		if strings.Contains(string(content), "func Benchmark") {
			hasBenchmarks = true
			break
		}
	}

	if !hasBenchmarks {
		return CheckResult{
			Name:    "Benchmark Files",
			Passed:  false,
			Message: fmt.Sprintf("Found %d test file(s) but no benchmark functions", len(matches)),
			Suggestions: []string{
				"Add benchmark functions: func BenchmarkXxx(b *testing.B) { ... }",
				"See https://golang.org/pkg/testing/#hdr-Benchmarks for examples",
			},
		}
	}

	return CheckResult{
		Name:    "Benchmark Files",
		Passed:  true,
		Message: fmt.Sprintf("Found %d test file(s) with benchmarks", len(matches)),
	}
}

func checkGitRepo() CheckResult {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()

	if err != nil {
		return CheckResult{
			Name:    "Git Repository (optional)",
			Passed:  true,
			Message: "Not a git repository (this is optional)",
		}
	}

	return CheckResult{
		Name:    "Git Repository",
		Passed:  true,
		Message: "Git repository detected",
	}
}

func checkSystemResources() CheckResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	totalMemMB := m.Sys / 1024 / 1024

	if totalMemMB < 100 {
		return CheckResult{
			Name:    "System Resources",
			Passed:  false,
			Message: fmt.Sprintf("Low memory available: %d MB", totalMemMB),
			Suggestions: []string{
				"Close other applications to free up memory",
				"Benchmarking may be unreliable with low memory",
			},
		}
	}

	return CheckResult{
		Name:    "System Resources",
		Passed:  true,
		Message: fmt.Sprintf("System memory: %d MB, Go runtime: %s", totalMemMB, runtime.Version()),
	}
}
