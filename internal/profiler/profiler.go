package profiler

import (
	"bytes"
	"fmt"
	"runtime/pprof"
	"sort"
	"strings"

	"github.com/alenon/gokanon/internal/models"
	"github.com/google/pprof/profile"
)

// Analyzer analyzes pprof profiles
type Analyzer struct {
	cpuProfile    *profile.Profile
	memoryProfile *profile.Profile
}

// NewAnalyzer creates a new profile analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// LoadCPUProfile loads a CPU profile from data
func (a *Analyzer) LoadCPUProfile(data []byte) error {
	prof, err := profile.Parse(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to parse CPU profile: %w", err)
	}
	a.cpuProfile = prof
	return nil
}

// LoadMemoryProfile loads a memory profile from data
func (a *Analyzer) LoadMemoryProfile(data []byte) error {
	prof, err := profile.Parse(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to parse memory profile: %w", err)
	}
	a.memoryProfile = prof
	return nil
}

// Analyze generates a complete profile summary
func (a *Analyzer) Analyze() (*models.ProfileSummary, error) {
	summary := &models.ProfileSummary{}

	// Analyze CPU profile if available
	if a.cpuProfile != nil {
		cpuFuncs, totalSamples, err := a.analyzeCPUProfile()
		if err != nil {
			return nil, fmt.Errorf("failed to analyze CPU profile: %w", err)
		}
		summary.CPUTopFunctions = cpuFuncs
		summary.TotalCPUSamples = totalSamples

		// Identify hot paths from CPU profile
		hotPaths := a.identifyHotPaths()
		summary.HotPaths = hotPaths
	}

	// Analyze memory profile if available
	if a.memoryProfile != nil {
		memFuncs, totalBytes, err := a.analyzeMemoryProfile()
		if err != nil {
			return nil, fmt.Errorf("failed to analyze memory profile: %w", err)
		}
		summary.MemoryTopFunctions = memFuncs
		summary.TotalMemoryBytes = totalBytes

		// Detect potential memory leaks
		leaks := a.detectMemoryLeaks()
		summary.MemoryLeaks = leaks
	}

	// Generate optimization suggestions
	suggestions := a.generateSuggestions(summary)
	summary.Suggestions = suggestions

	return summary, nil
}

// analyzeCPUProfile extracts top CPU-consuming functions
func (a *Analyzer) analyzeCPUProfile() ([]models.FunctionProfile, int64, error) {
	if a.cpuProfile == nil {
		return nil, 0, nil
	}

	// Get total samples
	var totalSamples int64
	for _, sample := range a.cpuProfile.Sample {
		totalSamples += sample.Value[0]
	}

	if totalSamples == 0 {
		return nil, 0, nil
	}

	// Aggregate by function
	funcStats := make(map[string]*funcStat)
	for _, sample := range a.cpuProfile.Sample {
		value := sample.Value[0]
		if len(sample.Location) == 0 {
			continue
		}

		// Get the leaf function (top of stack)
		loc := sample.Location[0]
		if len(loc.Line) == 0 {
			continue
		}

		fn := loc.Line[0].Function
		if fn == nil {
			continue
		}

		funcName := fn.Name
		if stat, exists := funcStats[funcName]; exists {
			stat.flat += value
			stat.cum += value
		} else {
			funcStats[funcName] = &funcStat{
				name: funcName,
				flat: value,
				cum:  value,
			}
		}

		// Add to callers as well (cumulative)
		for i := 1; i < len(sample.Location); i++ {
			if len(sample.Location[i].Line) == 0 {
				continue
			}
			callerFn := sample.Location[i].Line[0].Function
			if callerFn == nil {
				continue
			}
			callerName := callerFn.Name
			if stat, exists := funcStats[callerName]; exists {
				stat.cum += value
			} else {
				funcStats[callerName] = &funcStat{
					name: callerName,
					flat: 0,
					cum:  value,
				}
			}
		}
	}

	// Convert to slice and sort by flat time
	var stats []*funcStat
	for _, stat := range funcStats {
		stats = append(stats, stat)
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].flat > stats[j].flat
	})

	// Take top 10
	topCount := 10
	if len(stats) < topCount {
		topCount = len(stats)
	}

	var result []models.FunctionProfile
	for i := 0; i < topCount; i++ {
		stat := stats[i]
		result = append(result, models.FunctionProfile{
			Name:        cleanFunctionName(stat.name),
			FlatPercent: float64(stat.flat) / float64(totalSamples) * 100,
			CumPercent:  float64(stat.cum) / float64(totalSamples) * 100,
			FlatValue:   stat.flat,
			CumValue:    stat.cum,
		})
	}

	return result, totalSamples, nil
}

// analyzeMemoryProfile extracts top memory-allocating functions
func (a *Analyzer) analyzeMemoryProfile() ([]models.FunctionProfile, int64, error) {
	if a.memoryProfile == nil {
		return nil, 0, nil
	}

	// Find the alloc_space index (bytes allocated)
	allocSpaceIdx := -1
	for i, st := range a.memoryProfile.SampleType {
		if st.Type == "alloc_space" {
			allocSpaceIdx = i
			break
		}
	}

	if allocSpaceIdx == -1 {
		return nil, 0, fmt.Errorf("alloc_space not found in memory profile")
	}

	// Get total bytes
	var totalBytes int64
	for _, sample := range a.memoryProfile.Sample {
		totalBytes += sample.Value[allocSpaceIdx]
	}

	if totalBytes == 0 {
		return nil, 0, nil
	}

	// Aggregate by function
	funcStats := make(map[string]*funcStat)
	for _, sample := range a.memoryProfile.Sample {
		value := sample.Value[allocSpaceIdx]
		if len(sample.Location) == 0 {
			continue
		}

		// Get the leaf function
		loc := sample.Location[0]
		if len(loc.Line) == 0 {
			continue
		}

		fn := loc.Line[0].Function
		if fn == nil {
			continue
		}

		funcName := fn.Name
		if stat, exists := funcStats[funcName]; exists {
			stat.flat += value
			stat.cum += value
		} else {
			funcStats[funcName] = &funcStat{
				name: funcName,
				flat: value,
				cum:  value,
			}
		}
	}

	// Convert to slice and sort
	var stats []*funcStat
	for _, stat := range funcStats {
		stats = append(stats, stat)
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].flat > stats[j].flat
	})

	// Take top 10
	topCount := 10
	if len(stats) < topCount {
		topCount = len(stats)
	}

	var result []models.FunctionProfile
	for i := 0; i < topCount; i++ {
		stat := stats[i]
		result = append(result, models.FunctionProfile{
			Name:        cleanFunctionName(stat.name),
			FlatPercent: float64(stat.flat) / float64(totalBytes) * 100,
			CumPercent:  float64(stat.cum) / float64(totalBytes) * 100,
			FlatValue:   stat.flat,
			CumValue:    stat.cum,
		})
	}

	return result, totalBytes, nil
}

// identifyHotPaths identifies critical execution paths from CPU profile
func (a *Analyzer) identifyHotPaths() []models.HotPath {
	if a.cpuProfile == nil {
		return nil
	}

	// Track call stacks and their frequencies
	pathCounts := make(map[string]*pathStat)

	var totalSamples int64
	for _, sample := range a.cpuProfile.Sample {
		value := sample.Value[0]
		totalSamples += value

		// Build call stack path
		var path []string
		for i := len(sample.Location) - 1; i >= 0; i-- {
			if len(sample.Location[i].Line) == 0 {
				continue
			}
			fn := sample.Location[i].Line[0].Function
			if fn != nil {
				path = append(path, cleanFunctionName(fn.Name))
			}
		}

		if len(path) == 0 {
			continue
		}

		pathKey := strings.Join(path, " -> ")
		if stat, exists := pathCounts[pathKey]; exists {
			stat.count += value
		} else {
			pathCounts[pathKey] = &pathStat{
				path:  path,
				count: value,
			}
		}
	}

	// Convert to slice and sort
	var paths []*pathStat
	for _, stat := range pathCounts {
		paths = append(paths, stat)
	}
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].count > paths[j].count
	})

	// Take top 5 paths that consume > 5% of time
	var result []models.HotPath
	for _, p := range paths {
		percentage := float64(p.count) / float64(totalSamples) * 100
		if percentage < 5.0 {
			break
		}
		if len(result) >= 5 {
			break
		}

		result = append(result, models.HotPath{
			Path:        p.path,
			Percentage:  percentage,
			Occurrences: p.count,
			Description: fmt.Sprintf("Critical path consuming %.1f%% of execution time", percentage),
		})
	}

	return result
}

// detectMemoryLeaks identifies potential memory leak patterns
func (a *Analyzer) detectMemoryLeaks() []models.MemoryLeak {
	if a.memoryProfile == nil {
		return nil
	}

	// Find alloc_space and inuse_space indices
	allocSpaceIdx := -1
	inuseSpaceIdx := -1
	for i, st := range a.memoryProfile.SampleType {
		if st.Type == "alloc_space" {
			allocSpaceIdx = i
		}
		if st.Type == "inuse_space" {
			inuseSpaceIdx = i
		}
	}

	if allocSpaceIdx == -1 || inuseSpaceIdx == -1 {
		return nil
	}

	// Look for functions with high allocations
	funcLeaks := make(map[string]*leakStat)

	for _, sample := range a.memoryProfile.Sample {
		allocated := sample.Value[allocSpaceIdx]
		inuse := sample.Value[inuseSpaceIdx]

		if len(sample.Location) == 0 || allocated == 0 {
			continue
		}

		loc := sample.Location[0]
		if len(loc.Line) == 0 {
			continue
		}

		fn := loc.Line[0].Function
		if fn == nil {
			continue
		}

		funcName := fn.Name
		if stat, exists := funcLeaks[funcName]; exists {
			stat.allocated += allocated
			stat.inuse += inuse
			stat.count++
		} else {
			funcLeaks[funcName] = &leakStat{
				name:      funcName,
				allocated: allocated,
				inuse:     inuse,
				count:     1,
			}
		}
	}

	var result []models.MemoryLeak
	for _, stat := range funcLeaks {
		// If a function allocated much more than it's using, it might be leaking
		if stat.allocated > stat.inuse*2 && stat.allocated > 1024*1024 { // > 1MB
			severity := "low"
			if stat.allocated > 10*1024*1024 {
				severity = "high"
			} else if stat.allocated > 5*1024*1024 {
				severity = "medium"
			}

			result = append(result, models.MemoryLeak{
				Function:    cleanFunctionName(stat.name),
				Allocations: stat.count,
				Bytes:       stat.allocated,
				Severity:    severity,
				Description: fmt.Sprintf("Allocated %s but much less in use - potential leak", formatBytes(stat.allocated)),
			})
		}
	}

	// Sort by severity and bytes
	sort.Slice(result, func(i, j int) bool {
		if result[i].Severity != result[j].Severity {
			severityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
			return severityOrder[result[i].Severity] > severityOrder[result[j].Severity]
		}
		return result[i].Bytes > result[j].Bytes
	})

	// Limit to top 5
	if len(result) > 5 {
		result = result[:5]
	}

	return result
}

// generateSuggestions generates optimization suggestions based on profile data
func (a *Analyzer) generateSuggestions(summary *models.ProfileSummary) []models.Suggestion {
	var suggestions []models.Suggestion

	// CPU suggestions
	if len(summary.CPUTopFunctions) > 0 {
		top := summary.CPUTopFunctions[0]
		if top.FlatPercent > 30 {
			suggestions = append(suggestions, models.Suggestion{
				Type:       "cpu",
				Severity:   "high",
				Function:   top.Name,
				Issue:      fmt.Sprintf("Function consumes %.1f%% of CPU time", top.FlatPercent),
				Suggestion: "Consider optimizing this hot function - profile it in isolation, look for unnecessary allocations, consider algorithmic improvements",
				Impact:     fmt.Sprintf("Could improve overall performance by up to %.0f%%", top.FlatPercent*0.7),
			})
		}
	}

	// Memory suggestions
	if len(summary.MemoryTopFunctions) > 0 {
		top := summary.MemoryTopFunctions[0]
		if top.FlatPercent > 40 {
			suggestions = append(suggestions, models.Suggestion{
				Type:       "memory",
				Severity:   "high",
				Function:   top.Name,
				Issue:      fmt.Sprintf("Function allocates %.1f%% of total memory", top.FlatPercent),
				Suggestion: "Consider using sync.Pool for reusable objects, or pre-allocate slices/maps with appropriate capacity",
				Impact:     "Could significantly reduce allocation pressure and GC overhead",
			})
		}
	}

	// Memory leak suggestions
	for _, leak := range summary.MemoryLeaks {
		if leak.Severity == "high" {
			suggestions = append(suggestions, models.Suggestion{
				Type:       "memory",
				Severity:   "high",
				Function:   leak.Function,
				Issue:      "Potential memory leak detected",
				Suggestion: "Review this function for retained references, unclosed resources, or unbounded caches",
				Impact:     "Could prevent memory growth and improve stability",
			})
		}
	}

	// Hot path suggestions
	if len(summary.HotPaths) > 0 {
		for _, path := range summary.HotPaths {
			if path.Percentage > 25 {
				suggestions = append(suggestions, models.Suggestion{
					Type:       "cpu",
					Severity:   "medium",
					Function:   strings.Join(path.Path, " -> "),
					Issue:      fmt.Sprintf("Hot path consuming %.1f%% of execution", path.Percentage),
					Suggestion: "Analyze this call chain for optimization opportunities - consider caching, lazy evaluation, or algorithmic improvements",
					Impact:     fmt.Sprintf("Optimizing this path could improve performance by %.0f-%.0f%%", path.Percentage*0.5, path.Percentage*0.8),
				})
				break // Only suggest one hot path
			}
		}
	}

	return suggestions
}

// Helper types
type funcStat struct {
	name string
	flat int64
	cum  int64
}

type pathStat struct {
	path  []string
	count int64
}

type leakStat struct {
	name      string
	allocated int64
	inuse     int64
	count     int64
}

// cleanFunctionName removes package paths and simplifies function names
func cleanFunctionName(name string) string {
	// Remove package path, keep only last part
	parts := strings.Split(name, "/")
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}

	// Remove type parameters
	if idx := strings.Index(name, "["); idx != -1 {
		name = name[:idx]
	}

	return name
}

// formatBytes formats bytes in human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetProfileTypes returns available profile types from pprof
func GetProfileTypes() []string {
	var types []string
	for _, p := range pprof.Profiles() {
		types = append(types, p.Name())
	}
	return types
}
