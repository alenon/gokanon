package models

import (
	"testing"
	"time"
)

func TestBenchmarkResult(t *testing.T) {
	result := BenchmarkResult{
		Name:        "BenchmarkTest",
		Iterations:  1000,
		NsPerOp:     123.45,
		BytesPerOp:  64,
		AllocsPerOp: 2,
		MBPerSec:    10.5,
	}

	if result.Name != "BenchmarkTest" {
		t.Errorf("Expected name BenchmarkTest, got %s", result.Name)
	}
	if result.Iterations != 1000 {
		t.Errorf("Expected iterations 1000, got %d", result.Iterations)
	}
	if result.NsPerOp != 123.45 {
		t.Errorf("Expected NsPerOp 123.45, got %f", result.NsPerOp)
	}
	if result.BytesPerOp != 64 {
		t.Errorf("Expected BytesPerOp 64, got %d", result.BytesPerOp)
	}
	if result.AllocsPerOp != 2 {
		t.Errorf("Expected AllocsPerOp 2, got %d", result.AllocsPerOp)
	}
	if result.MBPerSec != 10.5 {
		t.Errorf("Expected MBPerSec 10.5, got %f", result.MBPerSec)
	}
}

func TestBenchmarkRun(t *testing.T) {
	now := time.Now()
	duration := 5 * time.Second

	results := []BenchmarkResult{
		{Name: "BenchmarkA", Iterations: 100, NsPerOp: 100.0},
		{Name: "BenchmarkB", Iterations: 200, NsPerOp: 200.0},
	}

	run := BenchmarkRun{
		ID:        "run-123",
		Timestamp: now,
		Package:   "./examples",
		GoVersion: "go1.21.0",
		Results:   results,
		Command:   "go test -bench=.",
		Duration:  duration,
	}

	if run.ID != "run-123" {
		t.Errorf("Expected ID run-123, got %s", run.ID)
	}
	if !run.Timestamp.Equal(now) {
		t.Errorf("Expected timestamp %v, got %v", now, run.Timestamp)
	}
	if run.Package != "./examples" {
		t.Errorf("Expected package ./examples, got %s", run.Package)
	}
	if run.GoVersion != "go1.21.0" {
		t.Errorf("Expected GoVersion go1.21.0, got %s", run.GoVersion)
	}
	if len(run.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(run.Results))
	}
	if run.Command != "go test -bench=." {
		t.Errorf("Expected command 'go test -bench=.', got %s", run.Command)
	}
	if run.Duration != duration {
		t.Errorf("Expected duration %v, got %v", duration, run.Duration)
	}
}

func TestComparison(t *testing.T) {
	comp := Comparison{
		Name:         "BenchmarkTest",
		OldNsPerOp:   100.0,
		NewNsPerOp:   90.0,
		Delta:        -10.0,
		DeltaPercent: -10.0,
		Status:       "improved",
	}

	if comp.Name != "BenchmarkTest" {
		t.Errorf("Expected name BenchmarkTest, got %s", comp.Name)
	}
	if comp.OldNsPerOp != 100.0 {
		t.Errorf("Expected OldNsPerOp 100.0, got %f", comp.OldNsPerOp)
	}
	if comp.NewNsPerOp != 90.0 {
		t.Errorf("Expected NewNsPerOp 90.0, got %f", comp.NewNsPerOp)
	}
	if comp.Delta != -10.0 {
		t.Errorf("Expected Delta -10.0, got %f", comp.Delta)
	}
	if comp.DeltaPercent != -10.0 {
		t.Errorf("Expected DeltaPercent -10.0, got %f", comp.DeltaPercent)
	}
	if comp.Status != "improved" {
		t.Errorf("Expected status improved, got %s", comp.Status)
	}
}

func TestComparisonStatuses(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"improved", "improved"},
		{"degraded", "degraded"},
		{"same", "same"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := Comparison{Status: tt.expected}
			if comp.Status != tt.expected {
				t.Errorf("Expected status %s, got %s", tt.expected, comp.Status)
			}
		})
	}
}
