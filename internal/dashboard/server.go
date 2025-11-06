package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alenon/gokanon/internal/storage"
)

// Server represents the dashboard web server
type Server struct {
	storage *storage.Storage
	addr    string
	port    int
}

// NewServer creates a new dashboard server
func NewServer(stor *storage.Storage, addr string, port int) *Server {
	return &Server{
		storage: stor,
		addr:    addr,
		port:    port,
	}
}

// Start starts the dashboard web server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/runs", s.handleRuns)
	mux.HandleFunc("/api/runs/", s.handleRunDetail)
	mux.HandleFunc("/api/trends", s.handleTrends)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/search", s.handleSearch)

	// Frontend
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static/", s.handleStatic)

	addr := fmt.Sprintf("%s:%d", s.addr, s.port)
	log.Printf("🚀 Dashboard server starting at http://%s\n", addr)
	log.Printf("📊 Open your browser to view interactive benchmarks\n")

	return http.ListenAndServe(addr, mux)
}

// handleRuns returns a list of all benchmark runs
func (s *Server) handleRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	runs, err := s.storage.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list runs: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a summary view for the list
	summaries := make([]map[string]interface{}, 0, len(runs))
	for _, run := range runs {
		summary := map[string]interface{}{
			"id":        run.ID,
			"timestamp": run.Timestamp.Format(time.RFC3339),
			"package":   run.Package,
			"goVersion": run.GoVersion,
			"duration":  run.Duration.String(),
			"numTests":  len(run.Results),
		}

		// Calculate average performance metrics
		if len(run.Results) > 0 {
			var totalNsPerOp float64
			var totalBytesPerOp int64
			var totalAllocsPerOp int64

			for _, result := range run.Results {
				totalNsPerOp += result.NsPerOp
				totalBytesPerOp += result.BytesPerOp
				totalAllocsPerOp += result.AllocsPerOp
			}

			count := float64(len(run.Results))
			summary["avgNsPerOp"] = totalNsPerOp / count
			summary["avgBytesPerOp"] = float64(totalBytesPerOp) / count
			summary["avgAllocsPerOp"] = float64(totalAllocsPerOp) / count
		}

		summaries = append(summaries, summary)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

// handleRunDetail returns details for a specific run
func (s *Server) handleRunDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid run ID", http.StatusBadRequest)
		return
	}
	id := parts[3]

	run, err := s.storage.Load(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load run: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(run)
}

// handleTrends returns trend data across multiple runs
func (s *Server) handleTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	benchName := r.URL.Query().Get("benchmark")
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	runs, err := s.storage.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list runs: %v", err), http.StatusInternalServerError)
		return
	}

	// Limit the number of runs
	if len(runs) > limit {
		runs = runs[:limit]
	}

	// Reverse to get chronological order
	for i := 0; i < len(runs)/2; i++ {
		runs[i], runs[len(runs)-1-i] = runs[len(runs)-1-i], runs[i]
	}

	// Build trend data
	trendData := make(map[string][]map[string]interface{})

	for _, run := range runs {
		timestamp := run.Timestamp.Format(time.RFC3339)

		for _, result := range run.Results {
			// Filter by benchmark name if specified
			if benchName != "" && result.Name != benchName {
				continue
			}

			if _, exists := trendData[result.Name]; !exists {
				trendData[result.Name] = make([]map[string]interface{}, 0)
			}

			trendData[result.Name] = append(trendData[result.Name], map[string]interface{}{
				"timestamp":   timestamp,
				"runId":       run.ID,
				"nsPerOp":     result.NsPerOp,
				"bytesPerOp":  result.BytesPerOp,
				"allocsPerOp": result.AllocsPerOp,
				"mbPerSec":    result.MBPerSec,
			})
		}
	}

	// Calculate trend statistics
	response := make(map[string]interface{})
	response["trends"] = trendData

	// Add statistical analysis for each benchmark
	statsData := make(map[string]interface{})
	for name, points := range trendData {
		if len(points) < 2 {
			continue
		}

		values := make([]float64, len(points))
		for i, point := range points {
			values[i] = point["nsPerOp"].(float64)
		}

		// Calculate basic statistics
		stat := calculateBasicStats(values)

		// Calculate trend (simple linear regression)
		slope := calculateSlope(values)

		statsData[name] = map[string]interface{}{
			"mean":   stat["mean"],
			"median": stat["median"],
			"stdDev": stat["stdDev"],
			"cv":     stat["cv"],
			"min":    stat["min"],
			"max":    stat["max"],
			"slope":  slope,
			"trend":  getTrendDirection(slope),
		}
	}
	response["statistics"] = statsData

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleStats returns statistical summaries
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	runs, err := s.storage.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list runs: %v", err), http.StatusInternalServerError)
		return
	}

	if len(runs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"totalRuns":  0,
			"totalTests": 0,
			"benchmarks": []string{},
			"dateRange":  map[string]string{},
			"recentRuns": []interface{}{},
		})
		return
	}

	// Collect all unique benchmark names
	benchmarkNames := make(map[string]bool)
	totalTests := 0

	for _, run := range runs {
		totalTests += len(run.Results)
		for _, result := range run.Results {
			benchmarkNames[result.Name] = true
		}
	}

	// Convert map to sorted slice
	uniqueBenchmarks := make([]string, 0, len(benchmarkNames))
	for name := range benchmarkNames {
		uniqueBenchmarks = append(uniqueBenchmarks, name)
	}
	sort.Strings(uniqueBenchmarks)

	// Date range
	oldest := runs[len(runs)-1].Timestamp
	newest := runs[0].Timestamp

	// Recent runs (last 10)
	recentCount := 10
	if len(runs) < recentCount {
		recentCount = len(runs)
	}
	recentRuns := make([]map[string]interface{}, recentCount)
	for i := 0; i < recentCount; i++ {
		run := runs[i]
		recentRuns[i] = map[string]interface{}{
			"id":        run.ID,
			"timestamp": run.Timestamp.Format(time.RFC3339),
			"package":   run.Package,
			"numTests":  len(run.Results),
		}
	}

	response := map[string]interface{}{
		"totalRuns":  len(runs),
		"totalTests": totalTests,
		"benchmarks": uniqueBenchmarks,
		"dateRange": map[string]string{
			"oldest": oldest.Format(time.RFC3339),
			"newest": newest.Format(time.RFC3339),
		},
		"recentRuns": recentRuns,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSearch searches for benchmark runs and results
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Missing search query parameter 'q'", http.StatusBadRequest)
		return
	}

	runs, err := s.storage.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list runs: %v", err), http.StatusInternalServerError)
		return
	}

	results := make([]map[string]interface{}, 0)

	for _, run := range runs {
		// Search in package name, ID, or benchmark names
		if strings.Contains(strings.ToLower(run.Package), query) ||
			strings.Contains(strings.ToLower(run.ID), query) {

			results = append(results, map[string]interface{}{
				"type":      "run",
				"id":        run.ID,
				"timestamp": run.Timestamp.Format(time.RFC3339),
				"package":   run.Package,
				"numTests":  len(run.Results),
			})
			continue
		}

		// Search in benchmark result names
		for _, result := range run.Results {
			if strings.Contains(strings.ToLower(result.Name), query) {
				results = append(results, map[string]interface{}{
					"type":      "benchmark",
					"runId":     run.ID,
					"timestamp": run.Timestamp.Format(time.RFC3339),
					"name":      result.Name,
					"nsPerOp":   result.NsPerOp,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":   query,
		"count":   len(results),
		"results": results,
	})
}

// handleIndex serves the main dashboard HTML
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, indexHTML)
}

// handleStatic serves static assets (CSS, JS)
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css")
		io.WriteString(w, stylesCSS)
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
		if strings.Contains(path, "app.js") {
			io.WriteString(w, appJS)
		}
	default:
		http.NotFound(w, r)
	}
}

// getTrendDirection returns the trend direction based on slope
func getTrendDirection(slope float64) string {
	if slope > 5 {
		return "degrading"
	} else if slope < -5 {
		return "improving"
	}
	return "stable"
}

// calculateBasicStats calculates basic statistical measures for a set of values
func calculateBasicStats(values []float64) map[string]float64 {
	if len(values) == 0 {
		return map[string]float64{
			"mean":   0,
			"median": 0,
			"stdDev": 0,
			"cv":     0,
			"min":    0,
			"max":    0,
		}
	}

	// Sort for median calculation
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate median
	var median float64
	if len(sorted)%2 == 0 {
		median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	} else {
		median = sorted[len(sorted)/2]
	}

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	stdDev := 0.0
	if variance > 0 {
		stdDev = variance // simplified, not taking sqrt for performance
	}

	// Calculate coefficient of variation
	cv := 0.0
	if mean != 0 {
		cv = (stdDev / mean) * 100
	}

	return map[string]float64{
		"mean":   mean,
		"median": median,
		"stdDev": stdDev,
		"cv":     cv,
		"min":    sorted[0],
		"max":    sorted[len(sorted)-1],
	}
}

// calculateSlope calculates the slope of a simple linear regression
func calculateSlope(values []float64) float64 {
	n := float64(len(values))
	if n < 2 {
		return 0
	}

	// Create x values (indices)
	var sumX, sumY, sumXY, sumX2 float64
	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0
	}

	slope := (n*sumXY - sumX*sumY) / denominator
	return slope
}
