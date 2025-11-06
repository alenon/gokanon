package models

import "time"

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Name        string  `json:"name"`
	Iterations  int64   `json:"iterations"`
	NsPerOp     float64 `json:"ns_per_op"`
	BytesPerOp  int64   `json:"bytes_per_op,omitempty"`
	AllocsPerOp int64   `json:"allocs_per_op,omitempty"`
	MBPerSec    float64 `json:"mb_per_sec,omitempty"`
}

// BenchmarkRun represents a complete benchmark run with metadata
type BenchmarkRun struct {
	ID             string            `json:"id"`
	Timestamp      time.Time         `json:"timestamp"`
	Package        string            `json:"package"`
	GoVersion      string            `json:"go_version"`
	Results        []BenchmarkResult `json:"results"`
	Command        string            `json:"command"`
	Duration       time.Duration     `json:"duration"`
	CPUProfile     string            `json:"cpu_profile,omitempty"`     // Path to CPU profile file
	MemoryProfile  string            `json:"memory_profile,omitempty"`  // Path to memory profile file
	ProfileSummary *ProfileSummary   `json:"profile_summary,omitempty"` // Summary of profile analysis
}

// Comparison represents the difference between two benchmark results
type Comparison struct {
	Name         string  `json:"name"`
	OldNsPerOp   float64 `json:"old_ns_per_op"`
	NewNsPerOp   float64 `json:"new_ns_per_op"`
	Delta        float64 `json:"delta"`
	DeltaPercent float64 `json:"delta_percent"`
	Status       string  `json:"status"` // "improved", "degraded", "same"
}

// ProfileSummary contains analyzed profile data
type ProfileSummary struct {
	CPUTopFunctions    []FunctionProfile `json:"cpu_top_functions,omitempty"`
	MemoryTopFunctions []FunctionProfile `json:"memory_top_functions,omitempty"`
	MemoryLeaks        []MemoryLeak      `json:"memory_leaks,omitempty"`
	HotPaths           []HotPath         `json:"hot_paths,omitempty"`
	Suggestions        []Suggestion      `json:"suggestions,omitempty"`
	TotalCPUSamples    int64             `json:"total_cpu_samples,omitempty"`
	TotalMemoryBytes   int64             `json:"total_memory_bytes,omitempty"`
}

// FunctionProfile represents a function's profile metrics
type FunctionProfile struct {
	Name        string  `json:"name"`
	FlatPercent float64 `json:"flat_percent"` // Time spent in function itself
	CumPercent  float64 `json:"cum_percent"`  // Time spent in function + callees
	FlatValue   int64   `json:"flat_value"`   // Actual value (samples or bytes)
	CumValue    int64   `json:"cum_value"`    // Cumulative value
}

// MemoryLeak represents a potential memory leak
type MemoryLeak struct {
	Function    string `json:"function"`
	Allocations int64  `json:"allocations"`
	Bytes       int64  `json:"bytes"`
	Severity    string `json:"severity"` // "low", "medium", "high"
	Description string `json:"description"`
}

// HotPath represents a critical execution path
type HotPath struct {
	Path        []string `json:"path"`        // Call stack
	Percentage  float64  `json:"percentage"`  // Percentage of total time
	Occurrences int64    `json:"occurrences"` // Number of samples
	Description string   `json:"description"`
}

// Suggestion represents an optimization suggestion
type Suggestion struct {
	Type       string `json:"type"`     // "cpu", "memory", "algorithm"
	Severity   string `json:"severity"` // "low", "medium", "high"
	Function   string `json:"function"`
	Issue      string `json:"issue"`
	Suggestion string `json:"suggestion"`
	Impact     string `json:"impact"` // Expected performance improvement
}

// Baseline represents a saved baseline benchmark run
type Baseline struct {
	Name        string            `json:"name"`         // Baseline identifier (e.g., "v1.0", "main", "stable")
	RunID       string            `json:"run_id"`       // ID of the benchmark run used as baseline
	CreatedAt   time.Time         `json:"created_at"`   // When the baseline was created
	Description string            `json:"description"`  // Optional description
	Run         *BenchmarkRun     `json:"run,omitempty"` // Full benchmark run data
	Tags        map[string]string `json:"tags,omitempty"` // Additional metadata tags
}
