package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alenon/gokanon/internal/models"
)

func TestNewExporter(t *testing.T) {
	e := NewExporter()
	if e == nil {
		t.Fatal("Expected non-nil exporter")
	}
}

func TestToCSV(t *testing.T) {
	e := NewExporter()
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.csv")

	comparisons := []models.Comparison{
		{
			Name:         "BenchmarkA",
			OldNsPerOp:   100.0,
			NewNsPerOp:   90.0,
			Delta:        -10.0,
			DeltaPercent: -10.0,
			Status:       "improved",
		},
		{
			Name:         "BenchmarkB",
			OldNsPerOp:   200.0,
			NewNsPerOp:   220.0,
			Delta:        20.0,
			DeltaPercent: 10.0,
			Status:       "degraded",
		},
	}

	err := e.ToCSV(comparisons, filename)
	if err != nil {
		t.Fatalf("ToCSV failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected CSV file to exist at %s", filename)
	}

	// Read and verify content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	csvContent := string(content)

	// Check header
	if !strings.Contains(csvContent, "Benchmark,Old (ns/op),New (ns/op),Delta (ns/op),Delta (%),Status") {
		t.Error("Expected CSV header not found")
	}

	// Check data rows
	expectedContains := []string{
		"BenchmarkA",
		"100.00",
		"90.00",
		"-10.00",
		"improved",
		"BenchmarkB",
		"200.00",
		"220.00",
		"20.00",
		"10.00",
		"degraded",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(csvContent, expected) {
			t.Errorf("Expected CSV to contain %q", expected)
		}
	}
}

func TestToCSVEmpty(t *testing.T) {
	e := NewExporter()
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "empty.csv")

	comparisons := []models.Comparison{}

	err := e.ToCSV(comparisons, filename)
	if err != nil {
		t.Fatalf("ToCSV failed: %v", err)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	// Should still have header
	if !strings.Contains(string(content), "Benchmark") {
		t.Error("Expected CSV header even for empty data")
	}
}

func TestToMarkdown(t *testing.T) {
	e := NewExporter()
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.md")

	comparisons := []models.Comparison{
		{
			Name:         "BenchmarkA",
			OldNsPerOp:   100.0,
			NewNsPerOp:   90.0,
			Delta:        -10.0,
			DeltaPercent: -10.0,
			Status:       "improved",
		},
		{
			Name:         "BenchmarkB",
			OldNsPerOp:   200.0,
			NewNsPerOp:   220.0,
			Delta:        20.0,
			DeltaPercent: 10.0,
			Status:       "degraded",
		},
		{
			Name:         "BenchmarkC",
			OldNsPerOp:   300.0,
			NewNsPerOp:   305.0,
			Delta:        5.0,
			DeltaPercent: 1.67,
			Status:       "same",
		},
	}

	err := e.ToMarkdown(comparisons, "old-id", "new-id", filename)
	if err != nil {
		t.Fatalf("ToMarkdown failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected Markdown file to exist at %s", filename)
	}

	// Read and verify content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read Markdown file: %v", err)
	}

	mdContent := string(content)

	// Check markdown formatting
	expectedContains := []string{
		"# Benchmark Comparison",
		"Comparing: `old-id` vs `new-id`",
		"| Status | Benchmark | Old (ns/op) | New (ns/op) | Delta | Delta (%) |",
		"BenchmarkA",
		"BenchmarkB",
		"BenchmarkC",
		"100.00",
		"90.00",
		"-10.00",
		"🟢", // improved
		"🔴", // degraded
		"⚪", // same
		"## Summary",
		"Improved: 1",
		"Degraded: 1",
		"Unchanged: 1",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(mdContent, expected) {
			t.Errorf("Expected Markdown to contain %q", expected)
		}
	}
}

func TestToHTML(t *testing.T) {
	e := NewExporter()
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.html")

	comparisons := []models.Comparison{
		{
			Name:         "BenchmarkA",
			OldNsPerOp:   100.0,
			NewNsPerOp:   90.0,
			Delta:        -10.0,
			DeltaPercent: -10.0,
			Status:       "improved",
		},
		{
			Name:         "BenchmarkB",
			OldNsPerOp:   200.0,
			NewNsPerOp:   220.0,
			Delta:        20.0,
			DeltaPercent: 10.0,
			Status:       "degraded",
		},
	}

	err := e.ToHTML(comparisons, "old-id", "new-id", "2024-01-01 10:00:00", "2024-01-01 11:00:00", filename)
	if err != nil {
		t.Fatalf("ToHTML failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected HTML file to exist at %s", filename)
	}

	// Read and verify content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}

	htmlContent := string(content)

	// Check HTML structure
	expectedContains := []string{
		"<!DOCTYPE html>",
		"<html",
		"<head>",
		"<title>Benchmark Comparison Report</title>",
		"<style>",
		"<body>",
		"Benchmark Comparison Report",
		"old-id",
		"new-id",
		"2024-01-01 10:00:00",
		"2024-01-01 11:00:00",
		"BenchmarkA",
		"BenchmarkB",
		"100.00",
		"90.00",
		"200.00",
		"220.00",
		"✅",                // improved indicator (emoji)
		"❌",                // degraded indicator (emoji)
		"chart.js",         // Verify Chart.js is included (lowercase in CDN URL)
		"performanceChart", // Verify performance chart canvas
		"deltaChart",       // Verify delta chart canvas
		"<table>",
		"</table>",
		"</body>",
		"</html>",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(htmlContent, expected) {
			t.Errorf("Expected HTML to contain %q", expected)
		}
	}

	// Check for CSS styling
	if !strings.Contains(htmlContent, "font-family") {
		t.Error("Expected HTML to contain CSS styling")
	}
}

func TestToHTMLWithSummary(t *testing.T) {
	e := NewExporter()
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "summary.html")

	comparisons := []models.Comparison{
		{Status: "improved"},
		{Status: "improved"},
		{Status: "degraded"},
		{Status: "same"},
		{Status: "same"},
	}

	err := e.ToHTML(comparisons, "old", "new", "time1", "time2", filename)
	if err != nil {
		t.Fatalf("ToHTML failed: %v", err)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}

	htmlContent := string(content)

	// Check summary cards
	if !strings.Contains(htmlContent, "summary-card") {
		t.Error("Expected HTML to contain summary cards")
	}

	// Check counts (should show 2 improved, 1 degraded, 2 same)
	// This is a bit tricky to test precisely, but we can check for the structure
	if !strings.Contains(htmlContent, "Improved") {
		t.Error("Expected HTML to contain Improved label")
	}
	if !strings.Contains(htmlContent, "Degraded") {
		t.Error("Expected HTML to contain Degraded label")
	}
	if !strings.Contains(htmlContent, "Unchanged") {
		t.Error("Expected HTML to contain Unchanged label")
	}
}

func TestCountStatus(t *testing.T) {
	comparisons := []models.Comparison{
		{Status: "improved"},
		{Status: "improved"},
		{Status: "improved"},
		{Status: "degraded"},
		{Status: "degraded"},
		{Status: "same"},
	}

	improved, degraded, same := countStatus(comparisons)

	if improved != 3 {
		t.Errorf("Expected 3 improved, got %d", improved)
	}
	if degraded != 2 {
		t.Errorf("Expected 2 degraded, got %d", degraded)
	}
	if same != 1 {
		t.Errorf("Expected 1 same, got %d", same)
	}
}

func TestCountStatusEmpty(t *testing.T) {
	comparisons := []models.Comparison{}

	improved, degraded, same := countStatus(comparisons)

	if improved != 0 {
		t.Errorf("Expected 0 improved, got %d", improved)
	}
	if degraded != 0 {
		t.Errorf("Expected 0 degraded, got %d", degraded)
	}
	if same != 0 {
		t.Errorf("Expected 0 same, got %d", same)
	}
}

func TestExportInvalidPath(t *testing.T) {
	e := NewExporter()

	comparisons := []models.Comparison{
		{Name: "Test", OldNsPerOp: 100, NewNsPerOp: 110, Status: "degraded"},
	}

	// Try to write to an invalid path
	err := e.ToCSV(comparisons, "/invalid/path/that/does/not/exist/test.csv")
	if err == nil {
		t.Error("Expected error when writing to invalid path")
	}
}

func TestToMarkdownSpecialCharacters(t *testing.T) {
	e := NewExporter()
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "special.md")

	comparisons := []models.Comparison{
		{
			Name:         "Benchmark|WithPipe",
			OldNsPerOp:   100.0,
			NewNsPerOp:   90.0,
			Delta:        -10.0,
			DeltaPercent: -10.0,
			Status:       "improved",
		},
	}

	err := e.ToMarkdown(comparisons, "old", "new", filename)
	if err != nil {
		t.Fatalf("ToMarkdown failed: %v", err)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// The pipe character in benchmark name should still be present
	if !strings.Contains(string(content), "Benchmark|WithPipe") {
		t.Error("Expected benchmark name with pipe character")
	}
}
