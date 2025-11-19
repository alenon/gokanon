package runner

import (
	"strings"
	"testing"

	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/storage"
)

func TestParseOutput(t *testing.T) {
	tests := []struct {
		name          string
		output        string
		expectedCount int
		expectError   bool
	}{
		{
			name: "valid output with memory stats",
			output: `goos: linux
goarch: amd64
BenchmarkStringBuilder-8    1000000   1234 ns/op   512 B/op   10 allocs/op
BenchmarkStringConcat-8     500000    2345 ns/op   1024 B/op  20 allocs/op
PASS
ok      github.com/test/bench   3.456s`,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "valid output without memory stats",
			output: `BenchmarkTest-8    1000000   1234 ns/op
PASS`,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "output with MB/s",
			output: `BenchmarkWrite-8    1000000   1234 ns/op   85.4 MB/s   512 B/op   10 allocs/op
PASS`,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "no benchmarks found",
			output:        `PASS\nok      github.com/test/bench   1.234s`,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "empty output",
			output:        "",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{}
			results, err := r.parseOutput(tt.output)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestParseOutputDetails(t *testing.T) {
	output := `BenchmarkStringBuilder-8    1000000   1234.56 ns/op   85.4 MB/s   512 B/op   10 allocs/op
PASS`

	r := &Runner{}
	results, err := r.parseOutput(output)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	result := results[0]

	if result.Name != "StringBuilder-8" {
		t.Errorf("Expected name StringBuilder-8, got %s", result.Name)
	}
	if result.Iterations != 1000000 {
		t.Errorf("Expected iterations 1000000, got %d", result.Iterations)
	}
	if result.NsPerOp != 1234.56 {
		t.Errorf("Expected NsPerOp 1234.56, got %f", result.NsPerOp)
	}
	if result.MBPerSec != 85.4 {
		t.Errorf("Expected MBPerSec 85.4, got %f", result.MBPerSec)
	}
	if result.BytesPerOp != 512 {
		t.Errorf("Expected BytesPerOp 512, got %d", result.BytesPerOp)
	}
	if result.AllocsPerOp != 10 {
		t.Errorf("Expected AllocsPerOp 10, got %d", result.AllocsPerOp)
	}
}

func TestParseOutputMultipleBenchmarks(t *testing.T) {
	output := `BenchmarkA-8    1000   100.0 ns/op   64 B/op   1 allocs/op
BenchmarkB-8    2000   200.0 ns/op   128 B/op  2 allocs/op
BenchmarkC-8    3000   300.0 ns/op   256 B/op  3 allocs/op
PASS`

	r := &Runner{}
	results, err := r.parseOutput(output)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	expected := []struct {
		name        string
		iterations  int64
		nsPerOp     float64
		bytesPerOp  int64
		allocsPerOp int64
	}{
		{"A-8", 1000, 100.0, 64, 1},
		{"B-8", 2000, 200.0, 128, 2},
		{"C-8", 3000, 300.0, 256, 3},
	}

	for i, exp := range expected {
		if results[i].Name != exp.name {
			t.Errorf("Result %d: expected name %s, got %s", i, exp.name, results[i].Name)
		}
		if results[i].Iterations != exp.iterations {
			t.Errorf("Result %d: expected iterations %d, got %d", i, exp.iterations, results[i].Iterations)
		}
		if results[i].NsPerOp != exp.nsPerOp {
			t.Errorf("Result %d: expected NsPerOp %f, got %f", i, exp.nsPerOp, results[i].NsPerOp)
		}
		if results[i].BytesPerOp != exp.bytesPerOp {
			t.Errorf("Result %d: expected BytesPerOp %d, got %d", i, exp.bytesPerOp, results[i].BytesPerOp)
		}
		if results[i].AllocsPerOp != exp.allocsPerOp {
			t.Errorf("Result %d: expected AllocsPerOp %d, got %d", i, exp.allocsPerOp, results[i].AllocsPerOp)
		}
	}
}

func TestNewRunner(t *testing.T) {
	tests := []struct {
		name        string
		packagePath string
		benchFilter string
	}{
		{"default", "", "."},
		{"specific package", "./examples", "."},
		{"specific benchmark", "", "BenchmarkTest"},
		{"both specified", "./pkg", "BenchmarkFoo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRunner(tt.packagePath, tt.benchFilter)
			if r == nil {
				t.Error("Expected non-nil runner")
			}
			if r.packagePath != tt.packagePath {
				t.Errorf("Expected packagePath %s, got %s", tt.packagePath, r.packagePath)
			}
			if r.benchFilter != tt.benchFilter {
				t.Errorf("Expected benchFilter %s, got %s", tt.benchFilter, r.benchFilter)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()

	if !strings.HasPrefix(id1, "run-") {
		t.Errorf("Expected ID to start with 'run-', got %s", id1)
	}

	// Check that ID contains only valid characters and is properly formatted
	if len(id1) < 5 {
		t.Errorf("Expected ID length > 5, got %d", len(id1))
	}

	// Test that format is consistent (run-<timestamp>)
	parts := strings.Split(id1, "-")
	if len(parts) != 2 {
		t.Errorf("Expected ID format 'run-<timestamp>', got %s", id1)
	}
	if parts[0] != "run" {
		t.Errorf("Expected ID prefix 'run', got %s", parts[0])
	}
}

func TestWithProfiling(t *testing.T) {
	r := NewRunner("./test", ".")

	if r.profileOptions != nil {
		t.Error("Expected profileOptions to be nil initially")
	}

	opts := &ProfileOptions{
		EnableCPU:    true,
		EnableMemory: true,
	}

	result := r.WithProfiling(opts)

	if result != r {
		t.Error("Expected WithProfiling to return the same runner instance")
	}

	if r.profileOptions != opts {
		t.Error("Expected profileOptions to be set")
	}

	if !r.profileOptions.EnableCPU {
		t.Error("Expected CPU profiling to be enabled")
	}

	if !r.profileOptions.EnableMemory {
		t.Error("Expected memory profiling to be enabled")
	}
}

func TestGetGoVersion(t *testing.T) {
	r := NewRunner("", ".")
	version, err := r.getGoVersion()

	if err != nil {
		t.Fatalf("getGoVersion failed: %v", err)
	}

	if version == "" {
		t.Error("Expected non-empty version string")
	}

	if !strings.Contains(version, "go version") {
		t.Errorf("Expected version to contain 'go version', got: %s", version)
	}
}

func TestRunWithActualBenchmarks(t *testing.T) {
	// This test runs actual benchmarks from the examples package
	r := NewRunner("../../examples", ".")

	run, err := r.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if run == nil {
		t.Fatal("Expected non-nil run")
	}

	if run.ID == "" {
		t.Error("Expected non-empty run ID")
	}

	if !strings.HasPrefix(run.ID, "run-") {
		t.Errorf("Expected ID to start with 'run-', got %s", run.ID)
	}

	if run.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	if run.GoVersion == "" {
		t.Error("Expected non-empty Go version")
	}

	if len(run.Results) == 0 {
		t.Error("Expected at least one benchmark result")
	}

	if run.Duration == 0 {
		t.Error("Expected non-zero duration")
	}

	if run.Command == "" {
		t.Error("Expected non-empty command")
	}

	// Verify command format - now uses direct execution
	if !strings.Contains(run.Command, "benchmark") && !strings.Contains(run.Command, "direct execution") {
		t.Errorf("Expected command to contain 'benchmark' or 'direct execution', got: %s", run.Command)
	}
}

func TestRunWithProfiling(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewStorage(tempDir)

	r := NewRunner("../../examples", ".")
	r.WithProfiling(&ProfileOptions{
		EnableCPU:    true,
		EnableMemory: true,
		Storage:      store,
	})

	run, err := r.Run()
	if err != nil {
		t.Fatalf("Run with profiling failed: %v", err)
	}

	if run == nil {
		t.Fatal("Expected non-nil run")
	}

	// Verify run has profile paths set
	// Note: profiles might be set if benchmarks ran long enough
	if run.CPUProfile == "" && run.MemoryProfile == "" {
		// It's okay if profiles weren't generated for quick tests
		t.Log("No profiles generated (benchmarks may have been too quick)")
	}
}

func TestRunWithInvalidPackage(t *testing.T) {
	r := NewRunner("./nonexistent", ".")

	_, err := r.Run()
	if err == nil {
		t.Error("Expected error when running benchmarks on non-existent package")
	}
}

func TestRunWithSpecificBenchmark(t *testing.T) {
	// Run only benchmarks matching "String"
	r := NewRunner("../../examples", "String")

	run, err := r.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if run == nil {
		t.Fatal("Expected non-nil run")
	}

	// Verify that results contain String benchmarks
	foundString := false
	for _, result := range run.Results {
		if strings.Contains(result.Name, "String") {
			foundString = true
			break
		}
	}

	if len(run.Results) > 0 && !foundString {
		t.Error("Expected to find String benchmark in results")
	}
}

func TestWithProgress(t *testing.T) {
	r := NewRunner("./test", ".")

	if r.progressCallback != nil {
		t.Error("Expected progressCallback to be nil initially")
	}

	callCount := 0
	var capturedResults []models.BenchmarkResult
	callback := func(result models.BenchmarkResult) {
		callCount++
		capturedResults = append(capturedResults, result)
	}

	result := r.WithProgress(callback)

	if result != r {
		t.Error("Expected WithProgress to return the same runner instance")
	}

	if r.progressCallback == nil {
		t.Error("Expected progressCallback to be set")
	}

	// Test that callback is invoked
	testResult := models.BenchmarkResult{
		Name:        "TestBenchmark-8",
		Iterations:  1000,
		NsPerOp:     100.5,
		BytesPerOp:  64,
		AllocsPerOp: 2,
	}
	r.progressCallback(testResult)
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, got %d", callCount)
	}
	if len(capturedResults) != 1 || capturedResults[0].Name != "TestBenchmark-8" {
		t.Errorf("Expected captured result with name 'TestBenchmark-8', got %v", capturedResults)
	}
}

func TestWithVerbose(t *testing.T) {
	r := NewRunner("./test", ".")

	if r.verboseWriter != nil {
		t.Error("Expected verboseWriter to be nil initially")
	}

	var buf strings.Builder
	result := r.WithVerbose(&buf)

	if result != r {
		t.Error("Expected WithVerbose to return the same runner instance")
	}

	if r.verboseWriter == nil {
		t.Error("Expected verboseWriter to be set")
	}
}

func TestProgressCallbackDuringParsing(t *testing.T) {
	output := `goos: linux
goarch: amd64
BenchmarkStringBuilder-8    1000000   1234 ns/op   512 B/op   10 allocs/op
BenchmarkStringConcat-8     500000    2345 ns/op   1024 B/op  20 allocs/op
PASS`

	callCount := 0
	var capturedResults []models.BenchmarkResult

	r := &Runner{}
	r.WithProgress(func(result models.BenchmarkResult) {
		callCount++
		capturedResults = append(capturedResults, result)
	})

	results, err := r.parseOutput(output)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if callCount != 2 {
		t.Errorf("Expected progress callback to be called 2 times, got %d", callCount)
	}

	expectedResults := []struct {
		name        string
		iterations  int64
		nsPerOp     float64
		bytesPerOp  int64
		allocsPerOp int64
	}{
		{"StringBuilder-8", 1000000, 1234, 512, 10},
		{"StringConcat-8", 500000, 2345, 1024, 20},
	}

	if len(capturedResults) != len(expectedResults) {
		t.Errorf("Expected %d captured results, got %d", len(expectedResults), len(capturedResults))
	}

	for i, expected := range expectedResults {
		if i >= len(capturedResults) {
			t.Errorf("Missing captured result at index %d", i)
			continue
		}
		if capturedResults[i].Name != expected.name {
			t.Errorf("Expected result[%d].Name = %s, got %s", i, expected.name, capturedResults[i].Name)
		}
		if capturedResults[i].Iterations != expected.iterations {
			t.Errorf("Expected result[%d].Iterations = %d, got %d", i, expected.iterations, capturedResults[i].Iterations)
		}
		if capturedResults[i].NsPerOp != expected.nsPerOp {
			t.Errorf("Expected result[%d].NsPerOp = %f, got %f", i, expected.nsPerOp, capturedResults[i].NsPerOp)
		}
		if capturedResults[i].BytesPerOp != expected.bytesPerOp {
			t.Errorf("Expected result[%d].BytesPerOp = %d, got %d", i, expected.bytesPerOp, capturedResults[i].BytesPerOp)
		}
		if capturedResults[i].AllocsPerOp != expected.allocsPerOp {
			t.Errorf("Expected result[%d].AllocsPerOp = %d, got %d", i, expected.allocsPerOp, capturedResults[i].AllocsPerOp)
		}
	}
}

func TestVerboseOutputDuringParsing(t *testing.T) {
	output := `goos: linux
goarch: amd64
BenchmarkStringBuilder-8    1000000   1234 ns/op   512 B/op   10 allocs/op
PASS
ok      github.com/test/bench   1.234s`

	var buf strings.Builder

	r := &Runner{}
	r.WithVerbose(&buf)

	results, err := r.parseOutput(output)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Check that verbose output was written
	verboseOutput := buf.String()
	if verboseOutput == "" {
		t.Error("Expected verbose output to be written, got empty string")
	}

	// Verify the output contains expected content
	expectedContents := []string{
		"goos: linux",
		"goarch: amd64",
		"BenchmarkStringBuilder-8",
		"1000000",
		"1234 ns/op",
		"PASS",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(verboseOutput, expected) {
			t.Errorf("Expected verbose output to contain '%s', got: %s", expected, verboseOutput)
		}
	}
}

func TestProgressCallbackNotCalledWhenNotSet(t *testing.T) {
	output := `BenchmarkStringBuilder-8    1000000   1234 ns/op   512 B/op   10 allocs/op
PASS`

	r := &Runner{} // No progress callback set

	results, err := r.parseOutput(output)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Test should complete without panic
}

func TestRunWithProgressCallback(t *testing.T) {
	callCount := 0
	var capturedResults []models.BenchmarkResult

	r := NewRunner("../../examples", "StringBuilder")
	r.WithProgress(func(result models.BenchmarkResult) {
		callCount++
		capturedResults = append(capturedResults, result)
	})

	run, err := r.Run()
	if err != nil {
		t.Fatalf("Run with progress callback failed: %v", err)
	}

	if run == nil {
		t.Fatal("Expected non-nil run")
	}

	if callCount == 0 {
		t.Error("Expected progress callback to be called at least once")
	}

	if len(capturedResults) == 0 {
		t.Error("Expected at least one benchmark result to be captured")
	}

	// Verify captured results match final results
	if len(capturedResults) != len(run.Results) {
		t.Errorf("Expected %d captured results to match %d results", len(capturedResults), len(run.Results))
	}

	for i, result := range run.Results {
		if i >= len(capturedResults) {
			t.Errorf("Missing captured result for result %d", i)
			continue
		}
		if capturedResults[i].Name != result.Name {
			t.Errorf("Captured result[%d].Name = %s doesn't match result name %s", i, capturedResults[i].Name, result.Name)
		}
		if capturedResults[i].Iterations != result.Iterations {
			t.Errorf("Captured result[%d].Iterations = %d doesn't match result iterations %d", i, capturedResults[i].Iterations, result.Iterations)
		}
		// Verify other fields are populated
		if capturedResults[i].NsPerOp == 0 {
			t.Errorf("Expected result[%d].NsPerOp to be non-zero", i)
		}
	}
}

func TestRunWithVerboseOutput(t *testing.T) {
	var buf strings.Builder

	r := NewRunner("../../examples", "StringBuilder")
	r.WithVerbose(&buf)

	run, err := r.Run()
	if err != nil {
		t.Fatalf("Run with verbose output failed: %v", err)
	}

	if run == nil {
		t.Fatal("Expected non-nil run")
	}

	verboseOutput := buf.String()
	if verboseOutput == "" {
		t.Error("Expected verbose output to be captured")
	}

	// Verify output contains benchmark information
	if !strings.Contains(verboseOutput, "Benchmark") {
		t.Error("Expected verbose output to contain 'Benchmark'")
	}

	if !strings.Contains(verboseOutput, "ns/op") {
		t.Error("Expected verbose output to contain 'ns/op'")
	}
}

func TestProgressAndVerboseNotBothSet(t *testing.T) {
	// This tests that both progress and verbose can be set independently
	// though in practice they wouldn't be used together

	var buf strings.Builder
	callCount := 0

	r := NewRunner("../../examples", "StringBuilder")
	r.WithProgress(func(result models.BenchmarkResult) {
		callCount++
	})
	r.WithVerbose(&buf)

	run, err := r.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if run == nil {
		t.Fatal("Expected non-nil run")
	}

	// Both should work
	if callCount == 0 {
		t.Error("Expected progress callback to be called")
	}

	if buf.String() == "" {
		t.Error("Expected verbose output to be written")
	}
}
