package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/alenon/gokanon/internal/models"
)

const (
	defaultDir = ".gokanon"
)

// Storage handles saving and loading benchmark results
type Storage struct {
	dir string
}

// NewStorage creates a new storage instance
func NewStorage(dir string) *Storage {
	if dir == "" {
		dir = defaultDir
	}
	return &Storage{dir: dir}
}

// Save saves a benchmark run to storage
func (s *Storage) Save(run *models.BenchmarkRun) error {
	// Ensure directory exists
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Create filename based on ID
	filename := filepath.Join(s.dir, run.ID+".json")

	// Marshal to JSON
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark run: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write benchmark run: %w", err)
	}

	return nil
}

// Load loads a benchmark run from storage by ID
func (s *Storage) Load(id string) (*models.BenchmarkRun, error) {
	filename := filepath.Join(s.dir, id+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read benchmark run: %w", err)
	}

	var run models.BenchmarkRun
	if err := json.Unmarshal(data, &run); err != nil {
		return nil, fmt.Errorf("failed to unmarshal benchmark run: %w", err)
	}

	return &run, nil
}

// List returns all available benchmark run IDs, sorted by timestamp (newest first)
func (s *Storage) List() ([]models.BenchmarkRun, error) {
	// Check if directory exists
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		return []models.BenchmarkRun{}, nil
	}

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}

	var runs []models.BenchmarkRun
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5] // Remove .json
		run, err := s.Load(id)
		if err != nil {
			continue // Skip invalid files
		}
		runs = append(runs, *run)
	}

	// Sort by timestamp, newest first
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].Timestamp.After(runs[j].Timestamp)
	})

	return runs, nil
}

// Delete removes a benchmark run from storage, including profile files
func (s *Storage) Delete(id string) error {
	filename := filepath.Join(s.dir, id+".json")
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete benchmark run: %w", err)
	}

	// Also delete profile directory if it exists
	profileDir := s.GetProfileDir(id)
	if _, err := os.Stat(profileDir); err == nil {
		if err := os.RemoveAll(profileDir); err != nil {
			// Log but don't fail if profile cleanup fails
			fmt.Fprintf(os.Stderr, "Warning: failed to delete profile directory: %v\n", err)
		}
	}

	return nil
}

// GetLatest returns the most recent benchmark run
func (s *Storage) GetLatest() (*models.BenchmarkRun, error) {
	runs, err := s.List()
	if err != nil {
		return nil, err
	}

	if len(runs) == 0 {
		return nil, fmt.Errorf("no benchmark runs found")
	}

	return &runs[0], nil
}

// GetProfileDir returns the profile directory for a given run ID
func (s *Storage) GetProfileDir(runID string) string {
	return filepath.Join(s.dir, "profiles", runID)
}

// GetCPUProfilePath returns the path to the CPU profile for a run
func (s *Storage) GetCPUProfilePath(runID string) string {
	return filepath.Join(s.GetProfileDir(runID), "cpu.prof")
}

// GetMemoryProfilePath returns the path to the memory profile for a run
func (s *Storage) GetMemoryProfilePath(runID string) string {
	return filepath.Join(s.GetProfileDir(runID), "mem.prof")
}

// SaveProfile saves a profile file to the storage
func (s *Storage) SaveProfile(runID, profileType string, data io.Reader) error {
	profileDir := s.GetProfileDir(runID)

	// Create profile directory
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Determine filename
	var filename string
	switch profileType {
	case "cpu":
		filename = s.GetCPUProfilePath(runID)
	case "memory", "mem":
		filename = s.GetMemoryProfilePath(runID)
	default:
		return fmt.Errorf("unknown profile type: %s", profileType)
	}

	// Create file
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create profile file: %w", err)
	}
	defer out.Close()

	// Copy data
	if _, err := io.Copy(out, data); err != nil {
		return fmt.Errorf("failed to write profile data: %w", err)
	}

	return nil
}

// LoadProfile loads a profile file from storage
func (s *Storage) LoadProfile(runID, profileType string) ([]byte, error) {
	var filename string
	switch profileType {
	case "cpu":
		filename = s.GetCPUProfilePath(runID)
	case "memory", "mem":
		filename = s.GetMemoryProfilePath(runID)
	default:
		return nil, fmt.Errorf("unknown profile type: %s", profileType)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile file: %w", err)
	}

	return data, nil
}

// HasProfile checks if a profile exists for a run
func (s *Storage) HasProfile(runID, profileType string) bool {
	var filename string
	switch profileType {
	case "cpu":
		filename = s.GetCPUProfilePath(runID)
	case "memory", "mem":
		filename = s.GetMemoryProfilePath(runID)
	default:
		return false
	}

	_, err := os.Stat(filename)
	return err == nil
}

// GetBaselineDir returns the baselines directory
func (s *Storage) GetBaselineDir() string {
	return filepath.Join(s.dir, "baselines")
}

// SaveBaseline saves a benchmark run as a baseline with the given name
func (s *Storage) SaveBaseline(name, runID, description string, tags map[string]string) (*models.Baseline, error) {
	// Load the run
	run, err := s.Load(runID)
	if err != nil {
		return nil, fmt.Errorf("failed to load run %s: %w", runID, err)
	}

	// Create baseline
	baseline := &models.Baseline{
		Name:        name,
		RunID:       runID,
		CreatedAt:   time.Now(),
		Description: description,
		Run:         run,
		Tags:        tags,
	}

	// Ensure baselines directory exists
	baselineDir := s.GetBaselineDir()
	if err := os.MkdirAll(baselineDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create baselines directory: %w", err)
	}

	// Save baseline
	filename := filepath.Join(baselineDir, name+".json")
	data, err := json.MarshalIndent(baseline, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal baseline: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write baseline: %w", err)
	}

	return baseline, nil
}

// LoadBaseline loads a baseline by name
func (s *Storage) LoadBaseline(name string) (*models.Baseline, error) {
	filename := filepath.Join(s.GetBaselineDir(), name+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read baseline %s: %w", name, err)
	}

	var baseline models.Baseline
	if err := json.Unmarshal(data, &baseline); err != nil {
		return nil, fmt.Errorf("failed to unmarshal baseline: %w", err)
	}

	return &baseline, nil
}

// ListBaselines returns all available baselines
func (s *Storage) ListBaselines() ([]models.Baseline, error) {
	baselineDir := s.GetBaselineDir()

	// Check if directory exists
	if _, err := os.Stat(baselineDir); os.IsNotExist(err) {
		return []models.Baseline{}, nil
	}

	entries, err := os.ReadDir(baselineDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read baselines directory: %w", err)
	}

	var baselines []models.Baseline
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5] // Remove .json
		baseline, err := s.LoadBaseline(name)
		if err != nil {
			continue // Skip invalid files
		}
		baselines = append(baselines, *baseline)
	}

	// Sort by creation time, newest first
	sort.Slice(baselines, func(i, j int) bool {
		return baselines[i].CreatedAt.After(baselines[j].CreatedAt)
	})

	return baselines, nil
}

// DeleteBaseline removes a baseline from storage
func (s *Storage) DeleteBaseline(name string) error {
	filename := filepath.Join(s.GetBaselineDir(), name+".json")
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete baseline %s: %w", name, err)
	}
	return nil
}

// HasBaseline checks if a baseline with the given name exists
func (s *Storage) HasBaseline(name string) bool {
	filename := filepath.Join(s.GetBaselineDir(), name+".json")
	_, err := os.Stat(filename)
	return err == nil
}
