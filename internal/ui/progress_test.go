package ui

import (
	"io"
	"testing"
	"time"

	"github.com/schollz/progressbar/v3"
)

func TestNewProgressBar(t *testing.T) {
	pb := NewProgressBar(100, "Test progress")
	if pb == nil {
		t.Fatal("NewProgressBar() returned nil")
	}
	if pb.bar == nil {
		t.Error("ProgressBar.bar is nil")
	}
}

func TestProgressBarAdd(t *testing.T) {
	pb := NewProgressBar(100, "Test progress")

	err := pb.Add(10)
	if err != nil {
		t.Errorf("Add() error = %v, want nil", err)
	}

	err = pb.Add(20)
	if err != nil {
		t.Errorf("Add() error = %v, want nil", err)
	}

	pb.Finish()
}

func TestProgressBarSet(t *testing.T) {
	pb := NewProgressBar(100, "Test progress")

	err := pb.Set(50)
	if err != nil {
		t.Errorf("Set() error = %v, want nil", err)
	}

	err = pb.Set(75)
	if err != nil {
		t.Errorf("Set() error = %v, want nil", err)
	}

	pb.Finish()
}

func TestProgressBarFinish(t *testing.T) {
	pb := NewProgressBar(100, "Test progress")

	err := pb.Set(100)
	if err != nil {
		t.Errorf("Set() error = %v, want nil", err)
	}

	err = pb.Finish()
	if err != nil {
		t.Errorf("Finish() error = %v, want nil", err)
	}
}

func TestProgressBarClear(t *testing.T) {
	pb := NewProgressBar(100, "Test progress")

	err := pb.Set(50)
	if err != nil {
		t.Errorf("Set() error = %v, want nil", err)
	}

	err = pb.Clear()
	if err != nil {
		t.Errorf("Clear() error = %v, want nil", err)
	}
}

func TestProgressBarDescribe(t *testing.T) {
	pb := NewProgressBar(100, "Initial description")

	// Change description
	pb.Describe("Updated description")

	// Should not panic
	pb.Add(10)
	pb.Finish()
}

func TestNewIndeterminateSpinner(t *testing.T) {
	spinner := NewIndeterminateSpinner("Loading")
	if spinner == nil {
		t.Fatal("NewIndeterminateSpinner() returned nil")
	}
	if spinner.bar == nil {
		t.Error("IndeterminateSpinner.bar is nil")
	}

	// Add some progress
	spinner.Add(1)
	time.Sleep(100 * time.Millisecond)
	spinner.Finish()
}

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner("Processing")
	if spinner == nil {
		t.Fatal("NewSpinner() returned nil")
	}
	if spinner.message == "" {
		t.Error("Spinner message is empty")
	}
	if spinner.message != "Processing" {
		t.Errorf("Spinner message = %q, want %q", spinner.message, "Processing")
	}
	if len(spinner.spinChars) == 0 {
		t.Error("Spinner has no spin characters")
	}
}

func TestSpinnerStartStop(t *testing.T) {
	spinner := NewSpinner("Test operation")

	// Start spinner
	spinner.Start()
	if !spinner.isRunning {
		t.Error("Spinner should be running after Start()")
	}

	// Let it spin for a bit
	time.Sleep(250 * time.Millisecond)

	// Stop spinner
	spinner.Stop()
	if spinner.isRunning {
		t.Error("Spinner should not be running after Stop()")
	}
}

func TestSpinnerDoubleStart(t *testing.T) {
	spinner := NewSpinner("Test operation")

	spinner.Start()
	// Starting again should be a no-op
	spinner.Start()

	time.Sleep(100 * time.Millisecond)
	spinner.Stop()
}

func TestSpinnerDoubleStop(t *testing.T) {
	spinner := NewSpinner("Test operation")

	spinner.Start()
	time.Sleep(100 * time.Millisecond)

	spinner.Stop()
	// Stopping again should be a no-op
	spinner.Stop()
}

func TestSpinnerUpdateMessage(t *testing.T) {
	spinner := NewSpinner("Initial message")

	if spinner.message != "Initial message" {
		t.Errorf("Initial message = %q, want %q", spinner.message, "Initial message")
	}

	spinner.UpdateMessage("Updated message")

	if spinner.message != "Updated message" {
		t.Errorf("Updated message = %q, want %q", spinner.message, "Updated message")
	}
}

func TestSpinnerWithStartStop(t *testing.T) {
	spinner := NewSpinner("Long running operation")

	spinner.Start()
	if !spinner.isRunning {
		t.Error("Spinner should be running")
	}

	// Simulate some work
	time.Sleep(200 * time.Millisecond)

	// Update message while running
	spinner.UpdateMessage("Still processing")

	time.Sleep(200 * time.Millisecond)

	spinner.Stop()
	if spinner.isRunning {
		t.Error("Spinner should not be running")
	}
}

func TestSpinnerStopWithoutStart(t *testing.T) {
	spinner := NewSpinner("Test")

	// Stopping without starting should not panic
	spinner.Stop()

	if spinner.isRunning {
		t.Error("Spinner should not be running")
	}
}

func TestSpinnerCharacters(t *testing.T) {
	spinner := NewSpinner("Test")

	expectedChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	if len(spinner.spinChars) != len(expectedChars) {
		t.Errorf("Spinner has %d characters, want %d", len(spinner.spinChars), len(expectedChars))
	}

	for i, char := range expectedChars {
		if i < len(spinner.spinChars) && spinner.spinChars[i] != char {
			t.Errorf("spinChars[%d] = %q, want %q", i, spinner.spinChars[i], char)
		}
	}
}

func TestProgressBarFullCycle(t *testing.T) {
	pb := NewProgressBar(10, "Test")

	for i := 0; i < 10; i++ {
		if err := pb.Add(1); err != nil {
			t.Errorf("Add() at iteration %d: %v", i, err)
		}
	}

	if err := pb.Finish(); err != nil {
		t.Errorf("Finish() error: %v", err)
	}
}

func TestMultipleProgressBars(t *testing.T) {
	pb1 := NewProgressBar(100, "Task 1")
	pb2 := NewProgressBar(100, "Task 2")

	pb1.Set(50)
	pb2.Set(75)

	pb1.Finish()
	pb2.Finish()
}

func BenchmarkProgressBarAdd(b *testing.B) {
	// Create progress bar with io.Discard to suppress output during benchmark
	bar := progressbar.NewOptions(b.N,
		progressbar.OptionSetWriter(io.Discard),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(40),
	)
	pb := &ProgressBar{bar: bar}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pb.Add(1)
	}
	pb.Finish()
}

func BenchmarkSpinnerStartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spinner := &Spinner{
			writer:    io.Discard,
			message:   "Benchmark",
			stopChan:  make(chan bool),
			spinChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		}
		spinner.Start()
		spinner.Stop()
	}
}

func TestSpinnerConcurrency(t *testing.T) {
	// Test that spinner handles concurrent Start/Stop gracefully
	spinner := NewSpinner("Concurrent test")

	done := make(chan bool)

	go func() {
		spinner.Start()
		time.Sleep(100 * time.Millisecond)
		spinner.Stop()
		done <- true
	}()

	<-done
}

func TestProgressBarDescribeChanges(t *testing.T) {
	pb := NewProgressBar(100, "Step 1")

	pb.Set(25)
	pb.Describe("Step 2")

	pb.Set(50)
	pb.Describe("Step 3")

	pb.Set(75)
	pb.Describe("Step 4")

	pb.Set(100)
	pb.Finish()
}

// Test edge cases
func TestProgressBarZeroMax(t *testing.T) {
	// Should not panic with zero max
	pb := NewProgressBar(0, "Empty task")
	if pb == nil {
		t.Fatal("NewProgressBar(0) returned nil")
	}
	pb.Finish()
}

func TestProgressBarNegativeValue(t *testing.T) {
	pb := NewProgressBar(100, "Test")

	// The underlying library may handle this, but we shouldn't panic
	pb.Set(-10)
	pb.Finish()
}

func TestSpinnerEmptyMessage(t *testing.T) {
	spinner := NewSpinner("")

	spinner.Start()
	time.Sleep(100 * time.Millisecond)
	spinner.Stop()

	// Should not panic with empty message
}

func TestSpinnerSpecialCharacters(t *testing.T) {
	specialMessages := []string{
		"Processing 你好",
		"Loading... 🚀",
		"Working on it™",
		"Test\nNewline",
		"Tab\there",
	}

	for _, msg := range specialMessages {
		t.Run(msg, func(t *testing.T) {
			spinner := NewSpinner(msg)
			spinner.Start()
			time.Sleep(50 * time.Millisecond)
			spinner.Stop()
			// Should handle special characters without crashing
		})
	}
}

func ExampleNewProgressBar() {
	pb := NewProgressBar(100, "Downloading")

	for i := 0; i < 100; i++ {
		pb.Add(1)
		time.Sleep(10 * time.Millisecond)
	}

	pb.Finish()
}

func ExampleNewSpinner() {
	spinner := NewSpinner("Processing data")
	spinner.Start()

	// Do some work
	time.Sleep(2 * time.Second)

	spinner.Stop()
}

// Test that repeatChar is accessible through the package
func TestRepeatCharHelper(t *testing.T) {
	result := repeatChar("-", 5)
	if result != "-----" {
		t.Errorf("repeatChar('-', 5) = %q, want '-----'", result)
	}
}

// Validate spinner currentChar wraps around correctly
func TestSpinnerCharacterRotation(t *testing.T) {
	spinner := NewSpinner("Test")

	// Initially should be 0
	if spinner.currentChar != 0 {
		t.Errorf("Initial currentChar = %d, want 0", spinner.currentChar)
	}

	spinner.Start()
	time.Sleep(1200 * time.Millisecond) // Should cycle through all chars multiple times

	spinner.Stop()

	// currentChar should still be valid (0 to len-1)
	if spinner.currentChar < 0 || spinner.currentChar >= len(spinner.spinChars) {
		t.Errorf("currentChar = %d, want 0-%d", spinner.currentChar, len(spinner.spinChars)-1)
	}
}
