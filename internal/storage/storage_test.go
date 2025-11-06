package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alenon/gokanon/internal/models"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		expectedDir string
	}{
		{"default directory", "", ".gokanon"},
		{"custom directory", "/tmp/test", "/tmp/test"},
		{"relative directory", "./benchmarks", "./benchmarks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStorage(tt.dir)
			if s.dir != tt.expectedDir {
				t.Errorf("Expected dir %s, got %s", tt.expectedDir, s.dir)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create test data
	run := &models.BenchmarkRun{
		ID:        "test-run-123",
		Timestamp: time.Now(),
		Package:   "./examples",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkA", Iterations: 1000, NsPerOp: 100.0},
			{Name: "BenchmarkB", Iterations: 2000, NsPerOp: 200.0},
		},
		Command:  "go test -bench=.",
		Duration: 5 * time.Second,
	}

	// Test Save
	err := s.Save(run)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	filename := filepath.Join(tempDir, run.ID+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to exist", filename)
	}

	// Test Load
	loaded, err := s.Load(run.ID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify loaded data
	if loaded.ID != run.ID {
		t.Errorf("Expected ID %s, got %s", run.ID, loaded.ID)
	}
	if loaded.Package != run.Package {
		t.Errorf("Expected Package %s, got %s", run.Package, loaded.Package)
	}
	if loaded.GoVersion != run.GoVersion {
		t.Errorf("Expected GoVersion %s, got %s", run.GoVersion, loaded.GoVersion)
	}
	if len(loaded.Results) != len(run.Results) {
		t.Errorf("Expected %d results, got %d", len(run.Results), len(loaded.Results))
	}
}

func TestList(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Test empty list
	runs, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(runs) != 0 {
		t.Errorf("Expected empty list, got %d items", len(runs))
	}

	// Add test data
	now := time.Now()
	for i := 0; i < 3; i++ {
		run := &models.BenchmarkRun{
			ID:        string(rune('A'+i)) + "-run",
			Timestamp: now.Add(time.Duration(i) * time.Second),
			Package:   "./test",
			GoVersion: "go1.21.0",
			Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
			Command:   "go test -bench=.",
			Duration:  time.Second,
		}
		if err := s.Save(run); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	// Test list with items
	runs, err = s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(runs) != 3 {
		t.Errorf("Expected 3 runs, got %d", len(runs))
	}

	// Verify sorting (newest first)
	for i := 0; i < len(runs)-1; i++ {
		if runs[i].Timestamp.Before(runs[i+1].Timestamp) {
			t.Errorf("Expected runs to be sorted newest first")
		}
	}
}

func TestDelete(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create test data
	run := &models.BenchmarkRun{
		ID:        "delete-test",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}

	if err := s.Save(run); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	filename := filepath.Join(tempDir, run.ID+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist before delete")
	}

	// Delete
	err := s.Delete(run.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("Expected file to not exist after delete")
	}

	// Test deleting non-existent file
	err = s.Delete("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent file")
	}
}

func TestGetLatest(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Test with no runs
	_, err := s.GetLatest()
	if err == nil {
		t.Error("Expected error when getting latest from empty storage")
	}

	// Add test data with different timestamps
	now := time.Now()
	runs := []*models.BenchmarkRun{
		{
			ID:        "old-run",
			Timestamp: now.Add(-2 * time.Hour),
			Package:   "./test",
			GoVersion: "go1.21.0",
			Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
			Command:   "go test -bench=.",
			Duration:  time.Second,
		},
		{
			ID:        "middle-run",
			Timestamp: now.Add(-1 * time.Hour),
			Package:   "./test",
			GoVersion: "go1.21.0",
			Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
			Command:   "go test -bench=.",
			Duration:  time.Second,
		},
		{
			ID:        "latest-run",
			Timestamp: now,
			Package:   "./test",
			GoVersion: "go1.21.0",
			Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
			Command:   "go test -bench=.",
			Duration:  time.Second,
		},
	}

	for _, run := range runs {
		if err := s.Save(run); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	// Get latest
	latest, err := s.GetLatest()
	if err != nil {
		t.Fatalf("GetLatest failed: %v", err)
	}

	if latest.ID != "latest-run" {
		t.Errorf("Expected latest run to have ID 'latest-run', got %s", latest.ID)
	}
}

func TestLoadNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	_, err := s.Load("non-existent")
	if err == nil {
		t.Error("Expected error when loading non-existent run")
	}
}

func TestListWithInvalidFiles(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create a directory in the storage dir (should be ignored)
	os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)

	// Create a non-JSON file (should be ignored)
	os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0644)

	// Create an invalid JSON file (should be ignored)
	os.WriteFile(filepath.Join(tempDir, "invalid.json"), []byte("not valid json"), 0644)

	// Create valid run
	run := &models.BenchmarkRun{
		ID:        "valid-run",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	s.Save(run)

	// List should only return the valid run
	runs, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(runs) != 1 {
		t.Errorf("Expected 1 run, got %d", len(runs))
	}
	if runs[0].ID != "valid-run" {
		t.Errorf("Expected ID 'valid-run', got %s", runs[0].ID)
	}
}

func TestGetProfilePaths(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "test-run-123"

	profileDir := s.GetProfileDir(runID)
	expectedProfileDir := filepath.Join(tempDir, "profiles", runID)
	if profileDir != expectedProfileDir {
		t.Errorf("Expected profile dir %s, got %s", expectedProfileDir, profileDir)
	}

	cpuPath := s.GetCPUProfilePath(runID)
	expectedCPUPath := filepath.Join(tempDir, "profiles", runID, "cpu.prof")
	if cpuPath != expectedCPUPath {
		t.Errorf("Expected CPU profile path %s, got %s", expectedCPUPath, cpuPath)
	}

	memPath := s.GetMemoryProfilePath(runID)
	expectedMemPath := filepath.Join(tempDir, "profiles", runID, "mem.prof")
	if memPath != expectedMemPath {
		t.Errorf("Expected memory profile path %s, got %s", expectedMemPath, memPath)
	}
}

func TestSaveAndLoadProfile(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "test-run-123"

	tests := []struct {
		name        string
		profileType string
		data        string
	}{
		{"cpu profile", "cpu", "cpu profile data"},
		{"memory profile", "memory", "memory profile data"},
		{"mem profile alias", "mem", "mem profile data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save profile
			reader := os.NewFile(0, tt.data)
			defer reader.Close()

			// Use a pipe to simulate io.Reader
			r, w, _ := os.Pipe()
			go func() {
				w.Write([]byte(tt.data))
				w.Close()
			}()

			err := s.SaveProfile(runID, tt.profileType, r)
			if err != nil {
				t.Fatalf("SaveProfile failed: %v", err)
			}

			// Load profile
			data, err := s.LoadProfile(runID, tt.profileType)
			if err != nil {
				t.Fatalf("LoadProfile failed: %v", err)
			}

			if string(data) != tt.data {
				t.Errorf("Expected profile data %s, got %s", tt.data, string(data))
			}
		})
	}
}

func TestSaveProfileInvalidType(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "test-run-123"

	r, w, _ := os.Pipe()
	go func() {
		w.Write([]byte("test data"))
		w.Close()
	}()

	err := s.SaveProfile(runID, "invalid", r)
	if err == nil {
		t.Error("Expected error when saving profile with invalid type")
	}
}

func TestLoadProfileInvalidType(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "test-run-123"

	_, err := s.LoadProfile(runID, "invalid")
	if err == nil {
		t.Error("Expected error when loading profile with invalid type")
	}
}

func TestLoadProfileNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "non-existent"

	_, err := s.LoadProfile(runID, "cpu")
	if err == nil {
		t.Error("Expected error when loading non-existent profile")
	}
}

func TestHasProfile(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "test-run-123"

	// Test non-existent profile
	if s.HasProfile(runID, "cpu") {
		t.Error("Expected HasProfile to return false for non-existent profile")
	}

	// Save a profile
	r, w, _ := os.Pipe()
	go func() {
		w.Write([]byte("cpu data"))
		w.Close()
	}()
	s.SaveProfile(runID, "cpu", r)

	// Test existing profile
	if !s.HasProfile(runID, "cpu") {
		t.Error("Expected HasProfile to return true for existing profile")
	}

	// Test invalid type
	if s.HasProfile(runID, "invalid") {
		t.Error("Expected HasProfile to return false for invalid type")
	}
}

func TestBaselineOperations(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create and save a run first
	run := &models.BenchmarkRun{
		ID:        "test-run-123",
		Timestamp: time.Now(),
		Package:   "./examples",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkA", Iterations: 1000, NsPerOp: 100.0},
		},
		Command:  "go test -bench=.",
		Duration: 5 * time.Second,
	}
	if err := s.Save(run); err != nil {
		t.Fatalf("Save run failed: %v", err)
	}

	// Test SaveBaseline
	tags := map[string]string{"version": "1.0", "env": "prod"}
	baseline, err := s.SaveBaseline("test-baseline", run.ID, "Test baseline description", tags)
	if err != nil {
		t.Fatalf("SaveBaseline failed: %v", err)
	}

	if baseline.Name != "test-baseline" {
		t.Errorf("Expected baseline name 'test-baseline', got %s", baseline.Name)
	}
	if baseline.RunID != run.ID {
		t.Errorf("Expected run ID %s, got %s", run.ID, baseline.RunID)
	}
	if baseline.Description != "Test baseline description" {
		t.Errorf("Expected description 'Test baseline description', got %s", baseline.Description)
	}
	if len(baseline.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(baseline.Tags))
	}

	// Test LoadBaseline
	loaded, err := s.LoadBaseline("test-baseline")
	if err != nil {
		t.Fatalf("LoadBaseline failed: %v", err)
	}
	if loaded.Name != baseline.Name {
		t.Errorf("Expected baseline name %s, got %s", baseline.Name, loaded.Name)
	}

	// Test HasBaseline
	if !s.HasBaseline("test-baseline") {
		t.Error("Expected HasBaseline to return true")
	}
	if s.HasBaseline("non-existent") {
		t.Error("Expected HasBaseline to return false for non-existent baseline")
	}

	// Test ListBaselines
	baselines, err := s.ListBaselines()
	if err != nil {
		t.Fatalf("ListBaselines failed: %v", err)
	}
	if len(baselines) != 1 {
		t.Errorf("Expected 1 baseline, got %d", len(baselines))
	}

	// Test DeleteBaseline
	err = s.DeleteBaseline("test-baseline")
	if err != nil {
		t.Fatalf("DeleteBaseline failed: %v", err)
	}
	if s.HasBaseline("test-baseline") {
		t.Error("Expected baseline to be deleted")
	}
}

func TestGetBaselineDir(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	baselineDir := s.GetBaselineDir()
	expected := filepath.Join(tempDir, "baselines")
	if baselineDir != expected {
		t.Errorf("Expected baseline dir %s, got %s", expected, baselineDir)
	}
}

func TestSaveBaselineWithNonExistentRun(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	_, err := s.SaveBaseline("test", "non-existent-run", "description", nil)
	if err == nil {
		t.Error("Expected error when saving baseline with non-existent run")
	}
}

func TestLoadBaselineNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	_, err := s.LoadBaseline("non-existent")
	if err == nil {
		t.Error("Expected error when loading non-existent baseline")
	}
}

func TestListBaselinesEmpty(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	baselines, err := s.ListBaselines()
	if err != nil {
		t.Fatalf("ListBaselines failed: %v", err)
	}
	if len(baselines) != 0 {
		t.Errorf("Expected empty list, got %d baselines", len(baselines))
	}
}

func TestListBaselinesWithInvalidFiles(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create baselines directory
	baselineDir := s.GetBaselineDir()
	os.MkdirAll(baselineDir, 0755)

	// Create invalid files that should be ignored
	os.Mkdir(filepath.Join(baselineDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(baselineDir, "test.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(baselineDir, "invalid.json"), []byte("not valid json"), 0644)

	// Create valid baseline
	run := &models.BenchmarkRun{
		ID:        "test-run",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	s.Save(run)
	s.SaveBaseline("valid-baseline", run.ID, "test", nil)

	baselines, err := s.ListBaselines()
	if err != nil {
		t.Fatalf("ListBaselines failed: %v", err)
	}
	if len(baselines) != 1 {
		t.Errorf("Expected 1 baseline, got %d", len(baselines))
	}
}

func TestDeleteBaselineNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	err := s.DeleteBaseline("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent baseline")
	}
}

func TestDeleteWithProfileDirectory(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)
	runID := "test-run-with-profiles"

	// Create and save a run
	run := &models.BenchmarkRun{
		ID:        runID,
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	s.Save(run)

	// Save a profile
	r, w, _ := os.Pipe()
	go func() {
		w.Write([]byte("cpu data"))
		w.Close()
	}()
	s.SaveProfile(runID, "cpu", r)

	// Verify profile exists
	if !s.HasProfile(runID, "cpu") {
		t.Fatal("Profile should exist before delete")
	}

	// Delete run
	err := s.Delete(runID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify profile directory is deleted
	profileDir := s.GetProfileDir(runID)
	if _, err := os.Stat(profileDir); !os.IsNotExist(err) {
		t.Error("Expected profile directory to be deleted")
	}
}

func TestListBaselinesOrdering(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create test runs
	now := time.Now()
	for i := 0; i < 3; i++ {
		run := &models.BenchmarkRun{
			ID:        string(rune('A'+i)) + "-run",
			Timestamp: now.Add(time.Duration(i) * time.Second),
			Package:   "./test",
			GoVersion: "go1.21.0",
			Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
			Command:   "go test -bench=.",
			Duration:  time.Second,
		}
		s.Save(run)

		// Save baseline with increasing timestamps
		time.Sleep(time.Millisecond * 10) // Ensure different timestamps
		_, err := s.SaveBaseline(string(rune('A'+i))+"-baseline", run.ID, "test", nil)
		if err != nil {
			t.Fatalf("SaveBaseline failed: %v", err)
		}
	}

	// List baselines
	baselines, err := s.ListBaselines()
	if err != nil {
		t.Fatalf("ListBaselines failed: %v", err)
	}

	// Verify sorting (newest first)
	for i := 0; i < len(baselines)-1; i++ {
		if baselines[i].CreatedAt.Before(baselines[i+1].CreatedAt) {
			t.Error("Expected baselines to be sorted newest first")
		}
	}
}

func TestSaveProfileCreateError(t *testing.T) {
	// Use a path that will cause mkdir to fail
	s := NewStorage("/proc/invalid-path-for-test")

	r, w, _ := os.Pipe()
	go func() {
		w.Write([]byte("test data"))
		w.Close()
	}()

	err := s.SaveProfile("test-run", "cpu", r)
	if err == nil {
		t.Error("Expected error when saving profile to invalid directory")
	}
}

func TestSaveBaselineWithTags(t *testing.T) {
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	// Create and save a run
	run := &models.BenchmarkRun{
		ID:        "test-run",
		Timestamp: time.Now(),
		Package:   "./test",
		GoVersion: "go1.21.0",
		Results:   []models.BenchmarkResult{{Name: "Test", Iterations: 100, NsPerOp: 100.0}},
		Command:   "go test -bench=.",
		Duration:  time.Second,
	}
	s.Save(run)

	// Save baseline with nil tags
	baseline, err := s.SaveBaseline("test-baseline-nil-tags", run.ID, "test", nil)
	if err != nil {
		t.Fatalf("SaveBaseline with nil tags failed: %v", err)
	}
	if baseline.Tags != nil {
		t.Error("Expected nil tags to remain nil")
	}

	// Save baseline with empty tags
	emptyTags := map[string]string{}
	baseline2, err := s.SaveBaseline("test-baseline-empty-tags", run.ID, "test", emptyTags)
	if err != nil {
		t.Fatalf("SaveBaseline with empty tags failed: %v", err)
	}
	if len(baseline2.Tags) != 0 {
		t.Error("Expected empty tags map")
	}
}

func TestGetLatestWithEmptyListError(t *testing.T) {
	// This is already covered by TestGetLatest, but adding explicit error path test
	tempDir := t.TempDir()
	s := NewStorage(tempDir)

	_, err := s.GetLatest()
	if err == nil {
		t.Error("Expected error when getting latest from empty storage")
	}

	expectedMsg := "no benchmark runs found"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain %q, got %q", expectedMsg, err.Error())
	}
}

func TestHasProfileWithInvalidPath(t *testing.T) {
	s := NewStorage("/nonexistent/path")

	// HasProfile should return false for non-existent storage
	if s.HasProfile("test-run", "cpu") {
		t.Error("Expected HasProfile to return false for non-existent storage path")
	}
}
