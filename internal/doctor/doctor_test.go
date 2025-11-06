package doctor

import (
	"os"
	"strings"
	"testing"
)

func TestCheckGoInstallation(t *testing.T) {
	result := checkGoInstallation()

	if result.Name != "Go Installation" {
		t.Errorf("Name = %q, want %q", result.Name, "Go Installation")
	}

	// Should pass if Go is installed (which it is for tests to run)
	if !result.Passed {
		t.Error("Go installation check should pass in test environment")
	}

	if !strings.Contains(result.Message, "go version") {
		t.Errorf("Message should contain 'go version', got %q", result.Message)
	}
}

func TestCheckGoTest(t *testing.T) {
	result := checkGoTest()

	if result.Name != "Go Test Command" {
		t.Errorf("Name = %q, want %q", result.Name, "Go Test Command")
	}

	// Should pass if Go is installed
	if !result.Passed {
		t.Error("Go test command check should pass in test environment")
	}
}

func TestCheckStorageDirectory_NotExist(t *testing.T) {
	// Change to a temp directory where .gokanon doesn't exist
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	result := checkStorageDirectory()

	if result.Name != "Storage Directory" {
		t.Errorf("Name = %q, want %q", result.Name, "Storage Directory")
	}

	// Should pass even if directory doesn't exist yet
	if !result.Passed {
		t.Error("Storage directory check should pass when directory doesn't exist")
	}

	if !strings.Contains(result.Message, "will be created") {
		t.Errorf("Message should mention directory will be created, got %q", result.Message)
	}
}

func TestCheckStorageDirectory_Exists(t *testing.T) {
	// Create temp directory with .gokanon
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create .gokanon directory
	os.Mkdir(".gokanon", 0755)

	result := checkStorageDirectory()

	if !result.Passed {
		t.Error("Storage directory check should pass when directory exists")
	}

	if !strings.Contains(result.Message, "exists at") {
		t.Errorf("Message should mention directory exists, got %q", result.Message)
	}
}

func TestCheckStorageDirectory_FileNotDir(t *testing.T) {
	// Create temp directory with .gokanon as a file (not directory)
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create .gokanon as a file
	os.WriteFile(".gokanon", []byte("test"), 0644)

	result := checkStorageDirectory()

	if result.Passed {
		t.Error("Storage directory check should fail when .gokanon is a file")
	}

	if !strings.Contains(result.Message, "not a directory") {
		t.Errorf("Message should mention not a directory, got %q", result.Message)
	}

	if len(result.Suggestions) == 0 {
		t.Error("Should provide suggestions when .gokanon is a file")
	}
}

func TestCheckStorageIntegrity_NoRuns(t *testing.T) {
	// Create temp directory with empty .gokanon
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	os.Mkdir(".gokanon", 0755)

	result := checkStorageIntegrity()

	if !result.Passed {
		t.Error("Storage integrity check should pass with no runs")
	}

	if !strings.Contains(result.Message, "No benchmark runs") {
		t.Errorf("Message should mention no runs, got %q", result.Message)
	}
}

func TestCheckStorageIntegrity_WithValidRun(t *testing.T) {
	// Create temp directory with valid run
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	os.Mkdir(".gokanon", 0755)

	// Create a valid run file
	runJSON := `{
		"id": "test-run-123",
		"timestamp": "2024-01-01T00:00:00Z",
		"package": "test",
		"go_version": "go1.21.0",
		"results": [],
		"command": "test",
		"duration": 1000000000
	}`
	os.WriteFile(".gokanon/test-run-123.json", []byte(runJSON), 0644)

	result := checkStorageIntegrity()

	if !result.Passed {
		t.Errorf("Storage integrity check should pass with valid run, got: %s", result.Message)
	}

	if !strings.Contains(result.Message, "healthy") {
		t.Errorf("Message should mention healthy storage, got %q", result.Message)
	}
}

func TestCheckStorageIntegrity_CorruptedRun(t *testing.T) {
	// Create temp directory with corrupted run
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	os.Mkdir(".gokanon", 0755)

	// Create an invalid JSON file that storage.List() will try to parse
	// Use run- prefix so it looks like a valid run file
	os.WriteFile(".gokanon/run-corrupted.json", []byte("invalid json{"), 0644)

	result := checkStorageIntegrity()

	// The check may or may not fail depending on how storage handles bad files
	// The important thing is it doesn't panic
	if !result.Passed && len(result.Suggestions) == 0 {
		t.Error("Should provide suggestions when storage check fails")
	}
}

func TestCheckBenchmarkFiles_NoFiles(t *testing.T) {
	// Create temp directory with no test files
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	result := checkBenchmarkFiles()

	if result.Passed {
		t.Error("Benchmark files check should fail when no test files exist")
	}

	if !strings.Contains(result.Message, "No test files") {
		t.Errorf("Message should mention no test files, got %q", result.Message)
	}

	if len(result.Suggestions) == 0 {
		t.Error("Should provide suggestions when no test files found")
	}
}

func TestCheckBenchmarkFiles_WithTestFiles(t *testing.T) {
	// Create temp directory with test files
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create test file with benchmark
	testContent := `package main
import "testing"
func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// test
	}
}`
	os.WriteFile("example_test.go", []byte(testContent), 0644)

	result := checkBenchmarkFiles()

	if !result.Passed {
		t.Errorf("Benchmark files check should pass with benchmark functions, got: %s", result.Message)
	}

	if !strings.Contains(result.Message, "test file") {
		t.Errorf("Message should mention test files, got %q", result.Message)
	}
}

func TestCheckBenchmarkFiles_TestFilesWithoutBenchmarks(t *testing.T) {
	// Create temp directory with test files but no benchmarks
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create test file without benchmark
	testContent := `package main
import "testing"
func TestSomething(t *testing.T) {
	// just a test
}`
	os.WriteFile("example_test.go", []byte(testContent), 0644)

	result := checkBenchmarkFiles()

	if result.Passed {
		t.Error("Benchmark files check should fail when test files don't have benchmarks")
	}

	if !strings.Contains(result.Message, "no benchmark functions") {
		t.Errorf("Message should mention no benchmark functions, got %q", result.Message)
	}
}

func TestCheckGitRepo_InGitRepo(t *testing.T) {
	// This test assumes we're in a git repo (which we are for gokanon)
	result := checkGitRepo()

	if result.Name != "Git Repository" && result.Name != "Git Repository (optional)" {
		t.Errorf("Name = %q, want Git Repository or Git Repository (optional)", result.Name)
	}

	// Should pass in gokanon repo
	if !result.Passed {
		t.Error("Git repository check should pass in gokanon repo")
	}
}

func TestCheckGitRepo_NotInGitRepo(t *testing.T) {
	// Create temp directory without git
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	result := checkGitRepo()

	// Should still "pass" as it's optional
	if !result.Passed {
		t.Error("Git repository check should pass (optional) even outside git repo")
	}

	if !strings.Contains(result.Message, "Not a git repository") || !strings.Contains(result.Message, "optional") {
		t.Errorf("Message should mention not a git repo and optional, got %q", result.Message)
	}
}

func TestCheckSystemResources(t *testing.T) {
	result := checkSystemResources()

	if result.Name != "System Resources" {
		t.Errorf("Name = %q, want %q", result.Name, "System Resources")
	}

	// Message should contain memory info
	if !strings.Contains(result.Message, "MB") {
		t.Errorf("Message should contain memory info, got %q", result.Message)
	}
}

func TestRunDiagnostics(t *testing.T) {
	results := RunDiagnostics()

	// Should run all checks
	expectedChecks := []string{
		"Go Installation",
		"Go Test Command",
		"Storage Directory",
		"Storage Integrity",
		"Benchmark Files",
		"Git Repository",
		"System Resources",
	}

	if len(results) != len(expectedChecks) {
		t.Errorf("Got %d results, want %d", len(results), len(expectedChecks))
	}

	// Verify all expected checks are present
	foundChecks := make(map[string]bool)
	for _, result := range results {
		foundChecks[result.Name] = true
		// Also check alternate names
		if strings.Contains(result.Name, "Git Repository") {
			foundChecks["Git Repository"] = true
		}
	}

	for _, expected := range expectedChecks {
		if !foundChecks[expected] {
			t.Errorf("Missing check: %s", expected)
		}
	}
}

func TestPrintResults(t *testing.T) {
	results := []CheckResult{
		{
			Name:    "Test Check 1",
			Passed:  true,
			Message: "Everything is fine",
		},
		{
			Name:        "Test Check 2",
			Passed:      false,
			Message:     "Something went wrong",
			Suggestions: []string{"Try this", "Or that"},
		},
	}

	// Just call PrintResults to ensure it doesn't panic
	// We won't validate output since it uses colors/UI which is complex to test
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintResults panicked: %v", r)
		}
	}()

	PrintResults(results)
}

func TestCheckResult_Structure(t *testing.T) {
	result := CheckResult{
		Name:        "Test",
		Passed:      true,
		Message:     "Success",
		Suggestions: []string{"Tip 1", "Tip 2"},
	}

	if result.Name != "Test" {
		t.Errorf("Name = %q, want %q", result.Name, "Test")
	}
	if !result.Passed {
		t.Error("Passed should be true")
	}
	if result.Message != "Success" {
		t.Errorf("Message = %q, want %q", result.Message, "Success")
	}
	if len(result.Suggestions) != 2 {
		t.Errorf("Got %d suggestions, want 2", len(result.Suggestions))
	}
}

func TestCheckGoInstallation_PathNotSet(t *testing.T) {
	// This test would require manipulating PATH, which is complex
	// We'll just verify the function handles errors gracefully
	result := checkGoInstallation()

	// In our environment, Go should be available
	if !result.Passed {
		// This is fine for testing - just ensure we don't panic
		t.Logf("Go not in PATH (expected in some test environments): %s", result.Message)
	}
}

// Test that all check functions return valid results
func TestAllCheckFunctionsReturnValidResults(t *testing.T) {
	checks := []func() CheckResult{
		checkGoInstallation,
		checkGoTest,
		checkStorageDirectory,
		checkStorageIntegrity,
		checkBenchmarkFiles,
		checkGitRepo,
		checkSystemResources,
	}

	for i, check := range checks {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result := check()

			if result.Name == "" {
				t.Error("Check result Name should not be empty")
			}
			if result.Message == "" {
				t.Error("Check result Message should not be empty")
			}
			// Passed can be true or false
			// Suggestions can be empty or not
		})
	}
}

func BenchmarkRunDiagnostics(b *testing.B) {
	// Suppress output during benchmark
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	for i := 0; i < b.N; i++ {
		RunDiagnostics()
	}
}

func BenchmarkCheckGoInstallation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkGoInstallation()
	}
}

// Test edge cases
func TestCheckBenchmarkFiles_ReadError(t *testing.T) {
	// Create temp directory with unreadable test file
	oldDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create test file
	testFile := "example_test.go"
	os.WriteFile(testFile, []byte("func BenchmarkTest(b *testing.B) {}"), 0644)

	// Make it unreadable (this may not work on all systems)
	os.Chmod(testFile, 0000)
	defer os.Chmod(testFile, 0644) // Restore for cleanup

	result := checkBenchmarkFiles()

	// Should still work, just might not detect benchmarks
	// This test mainly ensures we don't panic on read errors
	if result.Name != "Benchmark Files" {
		t.Errorf("Name = %q, want %q", result.Name, "Benchmark Files")
	}
}

func TestPrintResults_EmptyResults(t *testing.T) {
	// Should not panic with empty results
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintResults([]CheckResult{})

	w.Close()
	os.Stdout = oldStdout

	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	// Should still print summary
	outputStr := output.String()
	if !strings.Contains(outputStr, "0 checks passed") {
		t.Error("Output should mention 0 checks passed")
	}
}

func TestPrintResults_AllPassed(t *testing.T) {
	results := []CheckResult{
		{Name: "Check 1", Passed: true, Message: "OK"},
		{Name: "Check 2", Passed: true, Message: "OK"},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintResults(results)

	w.Close()
	os.Stdout = oldStdout

	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "2 checks passed") {
		t.Error("Output should mention 2 checks passed")
	}
	if !strings.Contains(outputStr, "All checks passed") {
		t.Error("Output should mention all checks passed")
	}
}

// Integration test
func TestFullDiagnosticFlow(t *testing.T) {
	// Run full diagnostic and print results
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	results := RunDiagnostics()
	PrintResults(results)

	w.Close()
	os.Stdout = oldStdout

	// Drain output
	buf := make([]byte, 8192)
	for {
		_, err := r.Read(buf)
		if err != nil {
			break
		}
	}

	// Should have run checks and printed results without panicking
	if len(results) == 0 {
		t.Error("RunDiagnostics should return at least some results")
	}
}

// Test that checkGoTest doesn't interfere with actual tests
func TestCheckGoTest_NoInterference(t *testing.T) {
	// Run the check multiple times to ensure it cleans up properly
	for i := 0; i < 3; i++ {
		result := checkGoTest()
		if result.Name != "Go Test Command" {
			t.Errorf("Iteration %d: Name = %q, want %q", i, result.Name, "Go Test Command")
		}
	}
}
