package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/storage"
)

// TestCalculateBasicStats tests the calculateBasicStats function
func TestCalculateBasicStats(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected map[string]float64
	}{
		{
			name:   "empty values",
			values: []float64{},
			expected: map[string]float64{
				"mean":   0,
				"median": 0,
				"stdDev": 0,
				"cv":     0,
				"min":    0,
				"max":    0,
			},
		},
		{
			name:   "single value",
			values: []float64{100.0},
			expected: map[string]float64{
				"mean":   100.0,
				"median": 100.0,
				"stdDev": 0,
				"cv":     0,
				"min":    100.0,
				"max":    100.0,
			},
		},
		{
			name:   "multiple values - even count",
			values: []float64{10.0, 20.0, 30.0, 40.0},
			expected: map[string]float64{
				"mean":   25.0,
				"median": 25.0,
				"min":    10.0,
				"max":    40.0,
			},
		},
		{
			name:   "multiple values - odd count",
			values: []float64{10.0, 20.0, 30.0},
			expected: map[string]float64{
				"mean":   20.0,
				"median": 20.0,
				"min":    10.0,
				"max":    30.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateBasicStats(tt.values)

			if result["mean"] != tt.expected["mean"] {
				t.Errorf("mean = %v, want %v", result["mean"], tt.expected["mean"])
			}
			if result["median"] != tt.expected["median"] {
				t.Errorf("median = %v, want %v", result["median"], tt.expected["median"])
			}
			if result["min"] != tt.expected["min"] {
				t.Errorf("min = %v, want %v", result["min"], tt.expected["min"])
			}
			if result["max"] != tt.expected["max"] {
				t.Errorf("max = %v, want %v", result["max"], tt.expected["max"])
			}
		})
	}
}

// TestCalculateSlope tests the calculateSlope function
func TestCalculateSlope(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
		delta    float64 // allowed difference for float comparison
	}{
		{
			name:     "empty values",
			values:   []float64{},
			expected: 0,
			delta:    0,
		},
		{
			name:     "single value",
			values:   []float64{100.0},
			expected: 0,
			delta:    0,
		},
		{
			name:     "increasing trend",
			values:   []float64{10.0, 20.0, 30.0, 40.0, 50.0},
			expected: 10.0,
			delta:    0.1,
		},
		{
			name:     "decreasing trend",
			values:   []float64{50.0, 40.0, 30.0, 20.0, 10.0},
			expected: -10.0,
			delta:    0.1,
		},
		{
			name:     "flat trend",
			values:   []float64{100.0, 100.0, 100.0, 100.0},
			expected: 0,
			delta:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSlope(tt.values)

			diff := result - tt.expected
			if diff < 0 {
				diff = -diff
			}

			if diff > tt.delta {
				t.Errorf("slope = %v, want %v (±%v)", result, tt.expected, tt.delta)
			}
		})
	}
}

// TestGetTrendDirection tests the getTrendDirection function
func TestGetTrendDirection(t *testing.T) {
	tests := []struct {
		name     string
		slope    float64
		expected string
	}{
		{
			name:     "strong degradation",
			slope:    10.0,
			expected: "degrading",
		},
		{
			name:     "weak degradation",
			slope:    3.0,
			expected: "stable",
		},
		{
			name:     "stable",
			slope:    0.0,
			expected: "stable",
		},
		{
			name:     "weak improvement",
			slope:    -3.0,
			expected: "stable",
		},
		{
			name:     "strong improvement",
			slope:    -10.0,
			expected: "improving",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTrendDirection(tt.slope)
			if result != tt.expected {
				t.Errorf("getTrendDirection(%v) = %v, want %v", tt.slope, result, tt.expected)
			}
		})
	}
}

// TestHandleRuns tests the /api/runs endpoint
func TestHandleRuns(t *testing.T) {
	// Create temporary storage
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test data
	run1 := &models.BenchmarkRun{
		ID:        "test-run-1",
		Timestamp: time.Now().Add(-1 * time.Hour),
		Package:   "test/package1",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkTest1", NsPerOp: 100.0, BytesPerOp: 64, AllocsPerOp: 1},
			{Name: "BenchmarkTest2", NsPerOp: 200.0, BytesPerOp: 128, AllocsPerOp: 2},
		},
	}

	run2 := &models.BenchmarkRun{
		ID:        "test-run-2",
		Timestamp: time.Now(),
		Package:   "test/package2",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkTest3", NsPerOp: 150.0, BytesPerOp: 96, AllocsPerOp: 1},
		},
	}

	if err := store.Save(run1); err != nil {
		t.Fatalf("failed to save test run 1: %v", err)
	}
	if err := store.Save(run2); err != nil {
		t.Fatalf("failed to save test run 2: %v", err)
	}

	// Create server
	server := NewServer(store, "localhost", 8080)

	// Test GET request
	req := httptest.NewRequest(http.MethodGet, "/api/runs", nil)
	w := httptest.NewRecorder()

	server.handleRuns(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	// Parse response
	var runs []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&runs); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(runs) != 2 {
		t.Errorf("got %d runs, want 2", len(runs))
	}

	// Verify response structure
	for _, run := range runs {
		if _, ok := run["id"]; !ok {
			t.Error("response missing 'id' field")
		}
		if _, ok := run["timestamp"]; !ok {
			t.Error("response missing 'timestamp' field")
		}
		if _, ok := run["package"]; !ok {
			t.Error("response missing 'package' field")
		}
		if _, ok := run["numTests"]; !ok {
			t.Error("response missing 'numTests' field")
		}
	}
}

// TestHandleRunsMethodNotAllowed tests method validation
func TestHandleRunsMethodNotAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodPost, "/api/runs", nil)
	w := httptest.NewRecorder()

	server.handleRuns(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusMethodNotAllowed)
	}
}

// TestHandleRunDetail tests the /api/runs/:id endpoint
func TestHandleRunDetail(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test run
	run := &models.BenchmarkRun{
		ID:        "test-run-detail",
		Timestamp: time.Now(),
		Package:   "test/package",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkTest", NsPerOp: 100.0},
		},
	}

	if err := store.Save(run); err != nil {
		t.Fatalf("failed to save test run: %v", err)
	}

	server := NewServer(store, "localhost", 8080)

	// Test valid run ID
	req := httptest.NewRequest(http.MethodGet, "/api/runs/test-run-detail", nil)
	w := httptest.NewRecorder()

	server.handleRunDetail(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var result models.BenchmarkRun
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ID != "test-run-detail" {
		t.Errorf("got ID %v, want test-run-detail", result.ID)
	}
}

// TestHandleRunDetailNotFound tests 404 handling
func TestHandleRunDetailNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/api/runs/nonexistent", nil)
	w := httptest.NewRecorder()

	server.handleRunDetail(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusNotFound)
	}
}

// TestHandleStats tests the /api/stats endpoint
func TestHandleStats(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test runs
	for i := 0; i < 3; i++ {
		run := &models.BenchmarkRun{
			ID:        fmt.Sprintf("test-run-%d", i),
			Timestamp: time.Now().Add(-time.Duration(i) * time.Hour),
			Package:   "test/package",
			GoVersion: "go1.21.0",
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkTest1", NsPerOp: 100.0},
				{Name: "BenchmarkTest2", NsPerOp: 200.0},
			},
		}
		if err := store.Save(run); err != nil {
			t.Fatalf("failed to save test run %d: %v", i, err)
		}
	}

	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	server.handleStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	if totalRuns, ok := stats["totalRuns"].(float64); !ok || totalRuns != 3 {
		t.Errorf("totalRuns = %v, want 3", stats["totalRuns"])
	}

	if totalTests, ok := stats["totalTests"].(float64); !ok || totalTests != 6 {
		t.Errorf("totalTests = %v, want 6", stats["totalTests"])
	}

	if benchmarks, ok := stats["benchmarks"].([]interface{}); !ok || len(benchmarks) != 2 {
		t.Errorf("benchmarks count = %v, want 2", len(benchmarks))
	}
}

// TestHandleSearch tests the /api/search endpoint
func TestHandleSearch(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test runs
	run1 := &models.BenchmarkRun{
		ID:        "test-run-search-1",
		Timestamp: time.Now(),
		Package:   "github.com/test/package",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkStringBuilder", NsPerOp: 100.0},
			{Name: "BenchmarkStringConcat", NsPerOp: 200.0},
		},
	}

	if err := store.Save(run1); err != nil {
		t.Fatalf("failed to save test run: %v", err)
	}

	server := NewServer(store, "localhost", 8080)

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "search by package",
			query:         "github",
			expectedCount: 1,
		},
		{
			name:          "search by benchmark name",
			query:         "String",
			expectedCount: 2,
		},
		{
			name:          "search with no results",
			query:         "nonexistent",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/search?q="+tt.query, nil)
			w := httptest.NewRecorder()

			server.handleSearch(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			count := int(result["count"].(float64))
			if count != tt.expectedCount {
				t.Errorf("result count = %v, want %v", count, tt.expectedCount)
			}
		})
	}
}

// TestHandleSearchMissingQuery tests validation
func TestHandleSearchMissingQuery(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()

	server.handleSearch(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

// TestHandleTrends tests the /api/trends endpoint
func TestHandleTrends(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test runs with trend data
	for i := 0; i < 5; i++ {
		run := &models.BenchmarkRun{
			ID:        fmt.Sprintf("test-run-%d", i),
			Timestamp: time.Now().Add(-time.Duration(5-i) * time.Hour),
			Package:   "test/package",
			GoVersion: "go1.21.0",
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkTest", NsPerOp: 100.0 + float64(i)*10.0},
			},
		}
		if err := store.Save(run); err != nil {
			t.Fatalf("failed to save test run %d: %v", i, err)
		}
	}

	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/api/trends?limit=5", nil)
	w := httptest.NewRecorder()

	server.handleTrends(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	trends, ok := result["trends"].(map[string]interface{})
	if !ok {
		t.Fatal("response missing 'trends' field")
	}

	benchData, ok := trends["BenchmarkTest"].([]interface{})
	if !ok {
		t.Fatal("trends missing 'BenchmarkTest' data")
	}

	if len(benchData) != 5 {
		t.Errorf("trend data length = %v, want 5", len(benchData))
	}

	// Verify statistics
	statistics, ok := result["statistics"].(map[string]interface{})
	if !ok {
		t.Fatal("response missing 'statistics' field")
	}

	benchStats, ok := statistics["BenchmarkTest"].(map[string]interface{})
	if !ok {
		t.Fatal("statistics missing 'BenchmarkTest' data")
	}

	if _, ok := benchStats["mean"]; !ok {
		t.Error("statistics missing 'mean' field")
	}
	if _, ok := benchStats["trend"]; !ok {
		t.Error("statistics missing 'trend' field")
	}
}

// TestHandleIndex tests the index HTML endpoint
func TestHandleIndex(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("content type = %v, want text/html; charset=utf-8", contentType)
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("response body is empty")
	}

	// Verify HTML contains expected elements
	expectedElements := []string{
		"<!DOCTYPE html>",
		"GoKanon Dashboard",
		"chart.js",
	}

	for _, elem := range expectedElements {
		if !contains(body, elem) {
			t.Errorf("response body missing expected element: %s", elem)
		}
	}
}

// TestHandleStatic tests static file serving
func TestHandleStatic(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	tests := []struct {
		name        string
		path        string
		wantCode    int
		wantContent string
	}{
		{
			name:        "CSS file",
			path:        "/static/styles.css",
			wantCode:    http.StatusOK,
			wantContent: "text/css",
		},
		{
			name:        "JS file",
			path:        "/static/app.js",
			wantCode:    http.StatusOK,
			wantContent: "application/javascript",
		},
		{
			name:        "unknown file",
			path:        "/static/unknown.txt",
			wantCode:    http.StatusNotFound,
			wantContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			server.handleStatic(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantCode)
			}

			if tt.wantContent != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.wantContent {
					t.Errorf("content type = %v, want %v", contentType, tt.wantContent)
				}
			}
		})
	}
}

// TestNewServer tests server creation
func TestNewServer(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	server := NewServer(store, "localhost", 8080)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}

	if server.storage != store {
		t.Error("server storage not set correctly")
	}

	if server.addr != "localhost" {
		t.Errorf("server addr = %v, want localhost", server.addr)
	}

	if server.port != 8080 {
		t.Errorf("server port = %v, want 8080", server.port)
	}
}

// TestHandleIndexNotFound tests 404 for non-root paths
func TestHandleIndexNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/invalid-path", nil)
	w := httptest.NewRecorder()

	server.handleIndex(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusNotFound)
	}
}

// TestHandleStatsEmptyStorage tests stats with no data
func TestHandleStatsEmptyStorage(t *testing.T) {
	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)
	server := NewServer(store, "localhost", 8080)

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	server.handleStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if totalRuns, ok := stats["totalRuns"].(float64); !ok || totalRuns != 0 {
		t.Errorf("totalRuns = %v, want 0", stats["totalRuns"])
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkCalculateBasicStats benchmarks the stats calculation
func BenchmarkCalculateBasicStats(b *testing.B) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateBasicStats(values)
	}
}

// BenchmarkCalculateSlope benchmarks slope calculation
func BenchmarkCalculateSlope(b *testing.B) {
	values := make([]float64, 100)
	for i := range values {
		values[i] = float64(i) * 1.5
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateSlope(values)
	}
}

// TestIntegration tests the full server integration
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test data
	run := &models.BenchmarkRun{
		ID:        "integration-test",
		Timestamp: time.Now(),
		Package:   "test/integration",
		GoVersion: "go1.21.0",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkIntegration", NsPerOp: 123.45},
		},
	}

	if err := store.Save(run); err != nil {
		t.Fatalf("failed to save test run: %v", err)
	}

	server := NewServer(store, "localhost", 0) // Use port 0 for automatic assignment

	// Test multiple endpoints in sequence
	endpoints := []struct {
		path       string
		wantStatus int
	}{
		{"/", http.StatusOK},
		{"/api/runs", http.StatusOK},
		{"/api/runs/integration-test", http.StatusOK},
		{"/api/stats", http.StatusOK},
		{"/api/search?q=integration", http.StatusOK},
		{"/api/trends?limit=10", http.StatusOK},
		{"/static/styles.css", http.StatusOK},
		{"/static/app.js", http.StatusOK},
	}

	for _, ep := range endpoints {
		t.Run(ep.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ep.path, nil)
			w := httptest.NewRecorder()

			// Route the request
			switch {
			case ep.path == "/":
				server.handleIndex(w, req)
			case ep.path == "/api/runs":
				server.handleRuns(w, req)
			case contains(ep.path, "/api/runs/"):
				server.handleRunDetail(w, req)
			case ep.path == "/api/stats":
				server.handleStats(w, req)
			case contains(ep.path, "/api/search"):
				server.handleSearch(w, req)
			case contains(ep.path, "/api/trends"):
				server.handleTrends(w, req)
			case contains(ep.path, "/static/"):
				server.handleStatic(w, req)
			}

			if w.Code != ep.wantStatus {
				t.Errorf("status code = %v, want %v", w.Code, ep.wantStatus)
			}
		})
	}
}

// TestConcurrentAccess tests concurrent access to endpoints
func TestConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concurrent test in short mode")
	}

	tmpDir := t.TempDir()
	store := storage.NewStorage(tmpDir)

	// Create test data
	for i := 0; i < 10; i++ {
		run := &models.BenchmarkRun{
			ID:        fmt.Sprintf("concurrent-test-%d", i),
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			Package:   "test/concurrent",
			GoVersion: "go1.21.0",
			Results: []models.BenchmarkResult{
				{Name: "BenchmarkConcurrent", NsPerOp: 100.0},
			},
		}
		if err := store.Save(run); err != nil {
			t.Fatalf("failed to save test run %d: %v", i, err)
		}
	}

	server := NewServer(store, "localhost", 8080)

	// Make concurrent requests
	const numRequests = 50
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/api/runs", nil)
			w := httptest.NewRecorder()
			server.handleRuns(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("concurrent request failed with status %v", w.Code)
			}
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}
}
