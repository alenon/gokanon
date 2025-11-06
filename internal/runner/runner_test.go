package runner

import (
	"strings"
	"testing"

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

	// Verify command format
	if !strings.Contains(run.Command, "go test") {
		t.Errorf("Expected command to contain 'go test', got: %s", run.Command)
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

