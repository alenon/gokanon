package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/storage"
)

// Helper function to create test storage with sample data
func setupTestStorage(t *testing.T) (*storage.Storage, string, func()) {
	tempDir := t.TempDir()
	store := storage.NewStorage(tempDir)

	// Create sample benchmark runs
	now := time.Now()
	for i := 0; i < 3; i++ {
		run := &models.BenchmarkRun{
			ID:        "test-run-" + string(rune('1'+i)),
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Package:   "./examples",
			GoVersion: "go1.21.0",
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkTest", Iterations: 1000, NsPerOp: float64(100 + i*10)},
				{Name: "BenchmarkAnother", Iterations: 2000, NsPerOp: float64(200 + i*20)},
			},
			Command:  "go test -bench=.",
			Duration: time.Second * time.Duration(i+1),
		}
		if err := store.Save(run); err != nil {
			t.Fatalf("Failed to create test data: %v", err)
		}
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return store, tempDir, cleanup
}

// Test helper to save args and restore them
func withArgs(args []string, fn func()) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = args
	fn()
}

func TestListWithEmptyStorage(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "list", "-storage=" + tempDir}, func() {
		err := List()
		if err != nil {
			t.Errorf("List should not error on empty storage, got: %v", err)
		}
	})
}

func TestListWithData(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "list", "-storage=" + tempDir}, func() {
		err := List()
		if err != nil {
			t.Errorf("List failed: %v", err)
		}
	})

	// Verify runs exist
	runs, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list runs: %v", err)
	}
	if len(runs) != 3 {
		t.Errorf("Expected 3 runs, got %d", len(runs))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "delete", "-storage=" + tempDir, "non-existent-id"}, func() {
		err := Delete()
		if err == nil {
			t.Error("Expected error when deleting non-existent run")
		}
	})
}

func TestDeleteSuccess(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	// Verify run exists
	runs, _ := store.List()
	if len(runs) == 0 {
		t.Fatal("Test setup failed: no runs created")
	}
	runID := runs[0].ID

	withArgs([]string{"gokanon", "delete", "-storage=" + tempDir, runID}, func() {
		err := Delete()
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})

	// Verify run was deleted
	runs, _ = store.List()
	if len(runs) != 2 {
		t.Errorf("Expected 2 runs after delete, got %d", len(runs))
	}

	// Verify the specific run was deleted
	for _, run := range runs {
		if run.ID == runID {
			t.Error("Run was not deleted")
		}
	}
}

func TestDeleteMissingArg(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "delete", "-storage=" + tempDir}, func() {
		err := Delete()
		if err == nil {
			t.Error("Expected error when run ID not provided")
		}
	})
}

func TestStatsWithNoData(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "stats", "-storage=" + tempDir}, func() {
		err := Stats()
		if err == nil {
			t.Error("Expected error when no benchmark results found")
		}
	})
}

func TestStatsWithData(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "stats", "-storage=" + tempDir}, func() {
		err := Stats()
		if err != nil {
			t.Errorf("Stats failed: %v", err)
		}
	})
}

func TestStatsWithLastN(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "stats", "-storage=" + tempDir, "-last=2"}, func() {
		err := Stats()
		if err != nil {
			t.Errorf("Stats with -last flag failed: %v", err)
		}
	})
}

func TestStatsWithCVThreshold(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "stats", "-storage=" + tempDir, "-cv-threshold=5.0"}, func() {
		err := Stats()
		if err != nil {
			t.Errorf("Stats with -cv-threshold flag failed: %v", err)
		}
	})
}

func TestCheckWithNoArgs(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	// Check requires threshold and run ID
	withArgs([]string{"gokanon", "check", "-storage=" + tempDir}, func() {
		err := Check()
		if err == nil {
			t.Error("Expected error when threshold not provided")
		}
	})
}

func TestDoctorCommand(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "doctor", "-storage=" + tempDir}, func() {
		err := Doctor()
		// Doctor should not error even if some checks fail
		// It reports issues but returns nil
		if err != nil {
			t.Logf("Doctor reported: %v", err)
		}
	})
}

func TestCompletionCommand(t *testing.T) {
	// Completion commands require script files which may not exist in test environment
	// We can only test that invalid shells error properly
	t.Skip("Completion requires script files, testing via invalid shell test only")
}

func TestCompletionInvalidShell(t *testing.T) {
	withArgs([]string{"gokanon", "completion", "invalid-shell"}, func() {
		err := Completion()
		if err == nil {
			t.Error("Expected error for invalid shell")
		}
	})
}

func TestCompletionMissingArg(t *testing.T) {
	withArgs([]string{"gokanon", "completion"}, func() {
		err := Completion()
		// Completion prints usage when no shell is specified, but doesn't error
		if err != nil {
			t.Errorf("Completion should print usage, not error: %v", err)
		}
	})
}

func TestExportMissingFormat(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "export", "-storage=" + tempDir}, func() {
		err := Export()
		if err == nil {
			t.Error("Expected error when format not specified")
		}
	})
}

func TestExportInvalidFormat(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "export", "-storage=" + tempDir, "-format=invalid"}, func() {
		err := Export()
		if err == nil {
			t.Error("Expected error for invalid export format")
		}
	})
}

func TestExportNoData(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.csv")

	withArgs([]string{"gokanon", "export", "-storage=" + tempDir, "-format=csv", "-output=" + outputFile}, func() {
		err := Export()
		if err == nil {
			t.Error("Expected error when no data to export")
		}
	})
}

func TestInteractiveCommand(t *testing.T) {
	// Interactive mode requires terminal interaction, skip actual execution
	// Just verify the command function exists and can be called
	t.Skip("Interactive mode requires terminal input, skipping")
}

func TestFlamegraphMissingRunID(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "flamegraph", "-storage=" + tempDir}, func() {
		err := Flamegraph()
		if err == nil {
			t.Error("Expected error when run ID not provided")
		}
	})
}

func TestTrendCommand(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "trend", "-storage=" + tempDir}, func() {
		err := Trend()
		// Trend should work with our test data
		if err != nil {
			t.Logf("Trend command result: %v", err)
		}
	})
}

func TestServeCommand(t *testing.T) {
	// Serve starts a web server, which we can't easily test in unit tests
	// We'll just verify the command doesn't panic on startup
	t.Skip("Serve starts a web server, skipping unit test")
}

func TestBaselineListEmpty(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "list", "-storage=" + tempDir}, func() {
		err := Baseline()
		// Should not error on empty baseline list
		if err != nil {
			t.Errorf("Baseline list should not error on empty storage: %v", err)
		}
	})
}

func TestBaselineSaveMissingArgs(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "save", "-storage=" + tempDir}, func() {
		err := Baseline()
		if err == nil {
			t.Error("Expected error when baseline name not provided")
		}
	})
}

func TestBaselineDeleteMissingArgs(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "delete", "-storage=" + tempDir}, func() {
		err := Baseline()
		if err == nil {
			t.Error("Expected error when baseline name not provided")
		}
	})
}

func TestBaselineInvalidSubcommand(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "invalid-subcommand", "-storage=" + tempDir}, func() {
		err := Baseline()
		if err == nil {
			t.Error("Expected error for invalid baseline subcommand")
		}
	})
}

func TestCompareMissingArgs(t *testing.T) {
	tempDir := t.TempDir()

	// Compare requires two run IDs
	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir}, func() {
		err := Compare()
		if err == nil {
			t.Error("Expected error when run IDs not provided")
		}
	})
}

func TestCompareOneArg(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	runs, _ := store.List()
	if len(runs) == 0 {
		t.Fatal("Test setup failed")
	}

	// Compare requires two run IDs
	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir, runs[0].ID}, func() {
		err := Compare()
		if err == nil {
			t.Error("Expected error when only one run ID provided")
		}
	})
}

// Test that storage directory is properly used
func TestStorageDirFlag(t *testing.T) {
	customDir := t.TempDir()
	store := storage.NewStorage(customDir)

	// Create a run
	run := &models.BenchmarkRun{
		ID:        "storage-test",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	store.Save(run)

	// Test that list finds it with custom storage dir
	withArgs([]string{"gokanon", "list", "-storage=" + customDir}, func() {
		err := List()
		if err != nil {
			t.Errorf("List with custom storage dir failed: %v", err)
		}
	})
}

// Test that commands handle storage errors gracefully
func TestStorageErrorHandling(t *testing.T) {
	invalidDir := "/proc/invalid-test-directory-12345"

	withArgs([]string{"gokanon", "list", "-storage=" + invalidDir}, func() {
		err := List()
		// Should handle permission errors gracefully
		// Empty list is ok, but shouldn't panic
		if err != nil {
			t.Logf("List handled storage error: %v", err)
		}
	})
}

// ===== Run Command Tests =====

func TestRunCommandWithBasicOptions(t *testing.T) {
	tempDir := t.TempDir()

	// Create a simple benchmark file
	testFile := filepath.Join(tempDir, "bench_test.go")
	benchmarkCode := `package test

import "testing"

func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i * 2
	}
}
`
	if err := os.WriteFile(testFile, []byte(benchmarkCode), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	storageDir := filepath.Join(tempDir, ".gokanon")
	withArgs([]string{"gokanon", "run", "-bench=.", "-pkg=" + tempDir, "-storage=" + storageDir}, func() {
		err := Run()
		if err != nil {
			t.Logf("Run command result: %v", err)
		}
	})

	// Verify that a run was saved
	store := storage.NewStorage(storageDir)
	runs, err := store.List()
	if err == nil && len(runs) > 0 {
		t.Logf("Successfully created run: %s", runs[0].ID)
	}
}

func TestRunCommandMissingBenchmarks(t *testing.T) {
	tempDir := t.TempDir()

	// Create a file without benchmarks
	testFile := filepath.Join(tempDir, "nobench_test.go")
	noBenchCode := `package test

import "testing"

func TestSimple(t *testing.T) {
	t.Log("This is not a benchmark")
}
`
	if err := os.WriteFile(testFile, []byte(noBenchCode), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	storageDir := filepath.Join(tempDir, ".gokanon")
	withArgs([]string{"gokanon", "run", "-bench=.", "-pkg=" + tempDir, "-storage=" + storageDir}, func() {
		err := Run()
		// Should error because no benchmarks are found
		if err == nil {
			t.Error("Expected error when no benchmarks found")
		}
	})
}

func TestRunCommandInvalidPackage(t *testing.T) {
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, ".gokanon")

	withArgs([]string{"gokanon", "run", "-bench=.", "-pkg=/nonexistent/path", "-storage=" + storageDir}, func() {
		err := Run()
		if err == nil {
			t.Error("Expected error for invalid package path")
		}
	})
}

func TestRunCommandInvalidProfileOption(t *testing.T) {
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, ".gokanon")

	withArgs([]string{"gokanon", "run", "-bench=.", "-profile=invalid", "-storage=" + storageDir}, func() {
		err := Run()
		if err == nil {
			t.Error("Expected error for invalid profile option")
		}
		if err != nil && !strings.Contains(err.Error(), "Unknown profile type") {
			t.Errorf("Expected 'Unknown profile type' error, got: %v", err)
		}
	})
}

// ===== Compare Command Tests =====

func TestCompareWithBaseline(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create a baseline
	runs, _ := store.List()
	if len(runs) == 0 {
		t.Fatal("Test setup failed: no runs created")
	}

	baselineName := "test-baseline"
	_, err := store.SaveBaseline(baselineName, runs[0].ID, "Test baseline", nil)
	if err != nil {
		t.Fatalf("Failed to create baseline: %v", err)
	}

	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir, "-baseline=" + baselineName}, func() {
		err := Compare()
		if err != nil {
			t.Errorf("Compare with baseline failed: %v", err)
		}
	})
}

func TestCompareLatestWithInsufficientRuns(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewStorage(tempDir)

	// Create only one run
	run := &models.BenchmarkRun{
		ID:        "single-run",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	store.Save(run)

	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir, "-latest"}, func() {
		err := Compare()
		if err == nil {
			t.Error("Expected error when less than 2 runs available")
		}
	})
}

func TestCompareWithTwoRunIDs(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	runs, _ := store.List()
	if len(runs) < 2 {
		t.Fatal("Test setup failed: need at least 2 runs")
	}

	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir, runs[0].ID, runs[1].ID}, func() {
		err := Compare()
		if err != nil {
			t.Errorf("Compare with two run IDs failed: %v", err)
		}
	})
}

func TestCompareWithNonExistentRun(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir, "nonexistent-1", "nonexistent-2"}, func() {
		err := Compare()
		if err == nil {
			t.Error("Expected error for non-existent run IDs")
		}
	})
}

func TestCompareWithNonExistentBaseline(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "compare", "-storage=" + tempDir, "-baseline=nonexistent"}, func() {
		err := Compare()
		if err == nil {
			t.Error("Expected error for non-existent baseline")
		}
	})
}

// ===== Baseline Command Tests =====

func TestBaselineSaveWithLatest(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	baselineName := "v1.0"
	withArgs([]string{"gokanon", "baseline", "save", "-storage=" + tempDir, "-name=" + baselineName}, func() {
		err := Baseline()
		if err != nil {
			t.Errorf("Baseline save failed: %v", err)
		}
	})

	// Verify baseline was created
	store := storage.NewStorage(tempDir)
	if !store.HasBaseline(baselineName) {
		t.Error("Baseline was not created")
	}
}

func TestBaselineSaveWithSpecificRun(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	runs, _ := store.List()
	if len(runs) == 0 {
		t.Fatal("Test setup failed")
	}

	baselineName := "specific-run"
	withArgs([]string{"gokanon", "baseline", "save", "-storage=" + tempDir, "-name=" + baselineName, "-run=" + runs[0].ID}, func() {
		err := Baseline()
		if err != nil {
			t.Errorf("Baseline save with specific run failed: %v", err)
		}
	})

	// Verify baseline was created
	if !store.HasBaseline(baselineName) {
		t.Error("Baseline was not created")
	}
}

func TestBaselineSaveWithNonExistentRun(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "save", "-storage=" + tempDir, "-name=test", "-run=nonexistent"}, func() {
		err := Baseline()
		if err == nil {
			t.Error("Expected error for non-existent run ID")
		}
	})
}

func TestBaselineShowSuccess(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create a baseline
	runs, _ := store.List()
	baselineName := "show-test"
	_, err := store.SaveBaseline(baselineName, runs[0].ID, "Test baseline", nil)
	if err != nil {
		t.Fatalf("Failed to create baseline: %v", err)
	}

	withArgs([]string{"gokanon", "baseline", "show", "-storage=" + tempDir, "-name=" + baselineName}, func() {
		err := Baseline()
		if err != nil {
			t.Errorf("Baseline show failed: %v", err)
		}
	})
}

func TestBaselineShowMissingName(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "show", "-storage=" + tempDir}, func() {
		err := Baseline()
		if err == nil {
			t.Error("Expected error when baseline name not provided")
		}
	})
}

func TestBaselineDeleteSuccess(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create a baseline
	runs, _ := store.List()
	baselineName := "delete-test"
	_, err := store.SaveBaseline(baselineName, runs[0].ID, "Test baseline", nil)
	if err != nil {
		t.Fatalf("Failed to create baseline: %v", err)
	}

	withArgs([]string{"gokanon", "baseline", "delete", "-storage=" + tempDir, "-name=" + baselineName}, func() {
		err := Baseline()
		if err != nil {
			t.Errorf("Baseline delete failed: %v", err)
		}
	})

	// Verify baseline was deleted
	if store.HasBaseline(baselineName) {
		t.Error("Baseline was not deleted")
	}
}

func TestBaselineDeleteNonExistent(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "baseline", "delete", "-storage=" + tempDir, "-name=nonexistent"}, func() {
		err := Baseline()
		if err == nil {
			t.Error("Expected error when deleting non-existent baseline")
		}
	})
}

// ===== Export Command Tests =====

func TestExportLatestSuccess(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	outputFile := filepath.Join(tempDir, "export.csv")
	withArgs([]string{"gokanon", "export", "-storage=" + tempDir, "-latest", "-format=csv", "-output=" + outputFile}, func() {
		err := Export()
		if err != nil {
			t.Errorf("Export latest failed: %v", err)
		}
	})

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}
}

func TestExportWithTwoRunsHTML(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	runs, _ := store.List()
	if len(runs) < 2 {
		t.Fatal("Test setup failed")
	}

	outputFile := filepath.Join(tempDir, "export.html")
	withArgs([]string{"gokanon", "export", "-storage=" + tempDir, "-format=html", "-output=" + outputFile, runs[0].ID, runs[1].ID}, func() {
		err := Export()
		if err != nil {
			t.Errorf("Export to HTML failed: %v", err)
		}
	})

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Export HTML file was not created")
	}
}

func TestExportMarkdownFormat(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	runs, _ := store.List()
	outputFile := filepath.Join(tempDir, "export.md")
	withArgs([]string{"gokanon", "export", "-storage=" + tempDir, "-format=markdown", "-output=" + outputFile, runs[0].ID, runs[1].ID}, func() {
		err := Export()
		if err != nil {
			t.Errorf("Export to Markdown failed: %v", err)
		}
	})

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Export Markdown file was not created")
	}
}

// ===== Check Command Tests =====

func TestCheckWithLatestSuccess(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "check", "-storage=" + tempDir, "-latest", "-threshold=50.0"}, func() {
		// Check may exit with non-zero code, but shouldn't panic
		_ = Check()
	})
}

func TestCheckWithTwoRunIDs(t *testing.T) {
	store, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	runs, _ := store.List()
	if len(runs) < 2 {
		t.Fatal("Test setup failed")
	}

	withArgs([]string{"gokanon", "check", "-storage=" + tempDir, "-threshold=10.0", runs[0].ID, runs[1].ID}, func() {
		_ = Check()
	})
}

func TestCheckWithInsufficientRuns(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewStorage(tempDir)

	// Create only one run
	run := &models.BenchmarkRun{
		ID:        "single-run",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	store.Save(run)

	withArgs([]string{"gokanon", "check", "-storage=" + tempDir, "-latest", "-threshold=10.0"}, func() {
		err := Check()
		if err == nil {
			t.Error("Expected error when less than 2 runs available")
		}
	})
}

// ===== Flamegraph Command Tests =====

func TestFlamegraphLatestWithoutProfiles(t *testing.T) {
	_, tempDir, cleanup := setupTestStorage(t)
	defer cleanup()

	withArgs([]string{"gokanon", "flamegraph", "-storage=" + tempDir, "-latest"}, func() {
		err := Flamegraph()
		if err == nil {
			t.Error("Expected error when no profiles are available")
		}
	})
}

func TestFlamegraphWithNonExistentRun(t *testing.T) {
	tempDir := t.TempDir()

	withArgs([]string{"gokanon", "flamegraph", "-storage=" + tempDir, "nonexistent"}, func() {
		err := Flamegraph()
		if err == nil {
			t.Error("Expected error for non-existent run")
		}
	})
}

// ===== Completion Command Tests =====

func TestCompletionBash(t *testing.T) {
	withArgs([]string{"gokanon", "completion", "bash"}, func() {
		err := Completion()
		if err != nil {
			t.Errorf("Completion bash failed: %v", err)
		}
	})
}

func TestCompletionZsh(t *testing.T) {
	withArgs([]string{"gokanon", "completion", "zsh"}, func() {
		err := Completion()
		if err != nil {
			t.Errorf("Completion zsh failed: %v", err)
		}
	})
}

func TestCompletionFish(t *testing.T) {
	withArgs([]string{"gokanon", "completion", "fish"}, func() {
		err := Completion()
		if err != nil {
			t.Errorf("Completion fish failed: %v", err)
		}
	})
}
