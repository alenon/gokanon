package profiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/pprof/profile"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer()
	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}
	if analyzer.cpuProfile != nil {
		t.Error("Expected cpuProfile to be nil initially")
	}
	if analyzer.memoryProfile != nil {
		t.Error("Expected memoryProfile to be nil initially")
	}
}

func TestLoadCPUProfile_Invalid(t *testing.T) {
	analyzer := NewAnalyzer()
	invalidData := []byte("not a valid profile")

	err := analyzer.LoadCPUProfile(invalidData)
	if err == nil {
		t.Error("Expected error for invalid CPU profile data")
	}
}

func TestLoadMemoryProfile_Invalid(t *testing.T) {
	analyzer := NewAnalyzer()
	invalidData := []byte("not a valid profile")

	err := analyzer.LoadMemoryProfile(invalidData)
	if err == nil {
		t.Error("Expected error for invalid memory profile data")
	}
}

func TestCleanFunctionName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple function",
			input:    "main.foo",
			expected: "main.foo",
		},
		{
			name:     "function with package path",
			input:    "github.com/user/project/pkg.Function",
			expected: "pkg.Function",
		},
		{
			name:     "function with generics",
			input:    "main.Process[int]",
			expected: "main.Process",
		},
		{
			name:     "complex path with generics",
			input:    "github.com/user/project/internal/pkg.Handler[*Request]",
			expected: "pkg.Handler",
		},
		{
			name:     "runtime function",
			input:    "runtime.mallocgc",
			expected: "runtime.mallocgc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanFunctionName(tt.input)
			if result != tt.expected {
				t.Errorf("cleanFunctionName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "bytes",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "kilobytes",
			bytes:    1024,
			expected: "1.0 KB",
		},
		{
			name:     "megabytes",
			bytes:    1024 * 1024,
			expected: "1.0 MB",
		},
		{
			name:     "gigabytes",
			bytes:    1024 * 1024 * 1024,
			expected: "1.0 GB",
		},
		{
			name:     "mixed megabytes",
			bytes:    1536 * 1024,
			expected: "1.5 MB",
		},
		{
			name:     "zero",
			bytes:    0,
			expected: "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

// createTestCPUProfile creates a simple CPU profile for testing
func createTestCPUProfile() []byte {
	// Create functions
	fooFunc := &profile.Function{
		ID:   1,
		Name: "main.foo",
	}
	barFunc := &profile.Function{
		ID:   2,
		Name: "main.bar",
	}

	// Create locations
	fooLoc := &profile.Location{
		ID:      1,
		Address: 0x1000,
		Line: []profile.Line{
			{Function: fooFunc},
		},
	}
	barLoc := &profile.Location{
		ID:      2,
		Address: 0x2000,
		Line: []profile.Line{
			{Function: barFunc},
		},
	}

	prof := &profile.Profile{
		SampleType: []*profile.ValueType{
			{Type: "samples", Unit: "count"},
			{Type: "cpu", Unit: "nanoseconds"},
		},
		Sample: []*profile.Sample{
			{
				Location: []*profile.Location{fooLoc},
				Value:    []int64{100, 1000000},
			},
			{
				Location: []*profile.Location{barLoc},
				Value:    []int64{50, 500000},
			},
		},
		Location:      []*profile.Location{fooLoc, barLoc},
		Function:      []*profile.Function{fooFunc, barFunc},
		TimeNanos:     1234567890,
		DurationNanos: 1000000000,
		PeriodType:    &profile.ValueType{Type: "cpu", Unit: "nanoseconds"},
		Period:        10000000,
	}

	var buf bytes.Buffer
	if err := prof.Write(&buf); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// createTestMemoryProfile creates a simple memory profile for testing
func createTestMemoryProfile() []byte {
	// Create functions
	allocFunc := &profile.Function{
		ID:   1,
		Name: "main.allocate",
	}
	bufferFunc := &profile.Function{
		ID:   2,
		Name: "main.buffer",
	}

	// Create locations
	allocLoc := &profile.Location{
		ID:      1,
		Address: 0x3000,
		Line: []profile.Line{
			{Function: allocFunc},
		},
	}
	bufferLoc := &profile.Location{
		ID:      2,
		Address: 0x4000,
		Line: []profile.Line{
			{Function: bufferFunc},
		},
	}

	prof := &profile.Profile{
		SampleType: []*profile.ValueType{
			{Type: "alloc_objects", Unit: "count"},
			{Type: "alloc_space", Unit: "bytes"},
			{Type: "inuse_objects", Unit: "count"},
			{Type: "inuse_space", Unit: "bytes"},
		},
		Sample: []*profile.Sample{
			{
				Location: []*profile.Location{allocLoc},
				Value:    []int64{100, 1024000, 10, 102400},
			},
			{
				Location: []*profile.Location{bufferLoc},
				Value:    []int64{50, 5120000, 5, 512000},
			},
		},
		Location:      []*profile.Location{allocLoc, bufferLoc},
		Function:      []*profile.Function{allocFunc, bufferFunc},
		TimeNanos:     1234567890,
		DurationNanos: 1000000000,
		PeriodType:    &profile.ValueType{Type: "space", Unit: "bytes"},
		Period:        524288,
	}

	var buf bytes.Buffer
	if err := prof.Write(&buf); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestLoadAndAnalyzeCPUProfile(t *testing.T) {
	analyzer := NewAnalyzer()
	cpuData := createTestCPUProfile()

	err := analyzer.LoadCPUProfile(cpuData)
	if err != nil {
		t.Fatalf("LoadCPUProfile() failed: %v", err)
	}

	if analyzer.cpuProfile == nil {
		t.Fatal("CPU profile not loaded")
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	if summary == nil {
		t.Fatal("Analyze() returned nil summary")
	}

	if len(summary.CPUTopFunctions) == 0 {
		t.Error("Expected CPU top functions, got none")
	}

	if summary.TotalCPUSamples == 0 {
		t.Error("Expected TotalCPUSamples > 0")
	}

	// Check that function names are cleaned
	for _, fn := range summary.CPUTopFunctions {
		if strings.Contains(fn.Name, "/") {
			t.Errorf("Function name not cleaned: %s", fn.Name)
		}
	}
}

func TestLoadAndAnalyzeMemoryProfile(t *testing.T) {
	analyzer := NewAnalyzer()
	memData := createTestMemoryProfile()

	err := analyzer.LoadMemoryProfile(memData)
	if err != nil {
		t.Fatalf("LoadMemoryProfile() failed: %v", err)
	}

	if analyzer.memoryProfile == nil {
		t.Fatal("Memory profile not loaded")
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	if summary == nil {
		t.Fatal("Analyze() returned nil summary")
	}

	if len(summary.MemoryTopFunctions) == 0 {
		t.Error("Expected memory top functions, got none")
	}

	if summary.TotalMemoryBytes == 0 {
		t.Error("Expected TotalMemoryBytes > 0")
	}
}

func TestAnalyzeBothProfiles(t *testing.T) {
	analyzer := NewAnalyzer()

	cpuData := createTestCPUProfile()
	err := analyzer.LoadCPUProfile(cpuData)
	if err != nil {
		t.Fatalf("LoadCPUProfile() failed: %v", err)
	}

	memData := createTestMemoryProfile()
	err = analyzer.LoadMemoryProfile(memData)
	if err != nil {
		t.Fatalf("LoadMemoryProfile() failed: %v", err)
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	if len(summary.CPUTopFunctions) == 0 {
		t.Error("Expected CPU top functions")
	}

	if len(summary.MemoryTopFunctions) == 0 {
		t.Error("Expected memory top functions")
	}

	if summary.TotalCPUSamples == 0 {
		t.Error("Expected TotalCPUSamples > 0")
	}

	if summary.TotalMemoryBytes == 0 {
		t.Error("Expected TotalMemoryBytes > 0")
	}
}

func TestAnalyzeWithNoProfiles(t *testing.T) {
	analyzer := NewAnalyzer()

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	if summary == nil {
		t.Fatal("Expected non-nil summary")
	}

	if len(summary.CPUTopFunctions) != 0 {
		t.Error("Expected no CPU functions")
	}

	if len(summary.MemoryTopFunctions) != 0 {
		t.Error("Expected no memory functions")
	}

	if summary.TotalCPUSamples != 0 {
		t.Error("Expected TotalCPUSamples = 0")
	}

	if summary.TotalMemoryBytes != 0 {
		t.Error("Expected TotalMemoryBytes = 0")
	}
}

func TestFunctionProfilePercentages(t *testing.T) {
	analyzer := NewAnalyzer()
	cpuData := createTestCPUProfile()

	err := analyzer.LoadCPUProfile(cpuData)
	if err != nil {
		t.Fatalf("LoadCPUProfile() failed: %v", err)
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	// Check that percentages are reasonable
	for _, fn := range summary.CPUTopFunctions {
		if fn.FlatPercent < 0 || fn.FlatPercent > 100 {
			t.Errorf("FlatPercent out of range: %f for %s", fn.FlatPercent, fn.Name)
		}
		if fn.CumPercent < 0 || fn.CumPercent > 100 {
			t.Errorf("CumPercent out of range: %f for %s", fn.CumPercent, fn.Name)
		}
		if fn.CumPercent < fn.FlatPercent {
			t.Errorf("CumPercent (%f) should be >= FlatPercent (%f) for %s",
				fn.CumPercent, fn.FlatPercent, fn.Name)
		}
	}
}

func TestSuggestionGeneration(t *testing.T) {
	analyzer := NewAnalyzer()

	cpuData := createTestCPUProfile()
	err := analyzer.LoadCPUProfile(cpuData)
	if err != nil {
		t.Fatalf("LoadCPUProfile() failed: %v", err)
	}

	memData := createTestMemoryProfile()
	err = analyzer.LoadMemoryProfile(memData)
	if err != nil {
		t.Fatalf("LoadMemoryProfile() failed: %v", err)
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	// Suggestions are optional, but check structure if they exist
	for _, sug := range summary.Suggestions {
		if sug.Type == "" {
			t.Error("Suggestion has empty Type")
		}
		if sug.Severity == "" {
			t.Error("Suggestion has empty Severity")
		}
		if sug.Issue == "" {
			t.Error("Suggestion has empty Issue")
		}
		if sug.Suggestion == "" {
			t.Error("Suggestion has empty Suggestion")
		}

		// Check valid types
		if sug.Type != "cpu" && sug.Type != "memory" && sug.Type != "algorithm" {
			t.Errorf("Invalid suggestion type: %s", sug.Type)
		}

		// Check valid severities
		if sug.Severity != "high" && sug.Severity != "medium" && sug.Severity != "low" {
			t.Errorf("Invalid suggestion severity: %s", sug.Severity)
		}
	}
}

func TestGetProfileTypes(t *testing.T) {
	types := GetProfileTypes()
	if len(types) == 0 {
		t.Error("Expected at least some profile types")
	}

	// Check that common types are present
	foundTypes := make(map[string]bool)
	for _, pt := range types {
		foundTypes[pt] = true
	}

	// These are standard Go pprof types
	expectedTypes := []string{"goroutine", "threadcreate", "heap", "allocs", "block", "mutex"}
	for _, expected := range expectedTypes {
		if !foundTypes[expected] {
			t.Logf("Note: Expected profile type %q not found (this may be OK depending on Go version)", expected)
		}
	}
}

func TestMemoryLeakDetection(t *testing.T) {
	// Create a profile with a potential leak (high allocation, low in-use)
	leakyFunc := &profile.Function{
		ID:   1,
		Name: "main.leakyFunction",
	}

	leakyLoc := &profile.Location{
		ID:      1,
		Address: 0x5000,
		Line: []profile.Line{
			{Function: leakyFunc},
		},
	}

	prof := &profile.Profile{
		SampleType: []*profile.ValueType{
			{Type: "alloc_objects", Unit: "count"},
			{Type: "alloc_space", Unit: "bytes"},
			{Type: "inuse_objects", Unit: "count"},
			{Type: "inuse_space", Unit: "bytes"},
		},
		Sample: []*profile.Sample{
			{
				Location: []*profile.Location{leakyLoc},
				// High allocation, low in-use (potential leak)
				Value: []int64{1000, 10 * 1024 * 1024, 10, 100 * 1024},
			},
		},
		Location:      []*profile.Location{leakyLoc},
		Function:      []*profile.Function{leakyFunc},
		TimeNanos:     1234567890,
		DurationNanos: 1000000000,
		PeriodType:    &profile.ValueType{Type: "space", Unit: "bytes"},
		Period:        524288,
	}

	var buf bytes.Buffer
	if err := prof.Write(&buf); err != nil {
		t.Fatalf("Failed to write profile: %v", err)
	}

	analyzer := NewAnalyzer()
	err := analyzer.LoadMemoryProfile(buf.Bytes())
	if err != nil {
		t.Fatalf("LoadMemoryProfile() failed: %v", err)
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}

	// Should detect at least one potential leak
	if len(summary.MemoryLeaks) == 0 {
		t.Error("Expected to detect memory leak")
	}

	// Check leak structure
	for _, leak := range summary.MemoryLeaks {
		if leak.Function == "" {
			t.Error("Leak has empty Function")
		}
		if leak.Severity == "" {
			t.Error("Leak has empty Severity")
		}
		if leak.Bytes == 0 {
			t.Error("Leak has zero Bytes")
		}
		if leak.Description == "" {
			t.Error("Leak has empty Description")
		}

		// Check valid severity
		if leak.Severity != "high" && leak.Severity != "medium" && leak.Severity != "low" {
			t.Errorf("Invalid leak severity: %s", leak.Severity)
		}
	}
}
