# GoKanon AI Analyzer Integration - Quick Start Guide

## Overview
This document provides a quick reference for implementing an AI analyzer feature in GoKanon.

## Core Data Structures to Work With

### 1. ProfileSummary (Most Important)
Located in `internal/models/benchmark.go`

```go
type ProfileSummary struct {
    CPUTopFunctions    []FunctionProfile  // Top CPU-consuming functions
    MemoryTopFunctions []FunctionProfile  // Top memory allocators
    MemoryLeaks        []MemoryLeak       // Potential memory issues
    HotPaths           []HotPath          // Critical execution paths
    Suggestions        []Suggestion       // Optimization recommendations
    TotalCPUSamples    int64
    TotalMemoryBytes   int64
}

type Suggestion struct {
    Type       string  // "cpu", "memory", "algorithm"
    Severity   string  // "low", "medium", "high"
    Function   string  // Function name
    Issue      string  // Problem description
    Suggestion string  // Solution description
    Impact     string  // Expected performance improvement
}
```

### 2. BenchmarkRun
```go
type BenchmarkRun struct {
    ID             string
    Timestamp      time.Time
    Package        string
    GoVersion      string
    Results        []BenchmarkResult
    Command        string
    Duration       time.Duration
    CPUProfile     string
    MemoryProfile  string
    ProfileSummary *ProfileSummary  // AI target data
}
```

### 3. Comparison
```go
type Comparison struct {
    Name         string
    OldNsPerOp   float64
    NewNsPerOp   float64
    Delta        float64
    DeltaPercent float64
    Status       string  // "improved", "degraded", "same"
}
```

## Primary Integration Point

### Location: `internal/runner/runner.go` - `handleProfiles()` method

**Current code (lines ~250-270)**:
```go
profileSummary, err := analyzer.Analyze()
if err != nil {
    fmt.Fprintf(os.Stderr, "Warning: failed to analyze profiles: %v\n", err)
}
// Stored in run.ProfileSummary
```

**AI Integration Pattern**:
```go
// After analyzer.Analyze()
if profileSummary != nil {
    // NEW: Call AI analyzer
    enhancedSuggestions := aianalyzer.EnhanceSuggestions(profileSummary)
    profileSummary.Suggestions = enhancedSuggestions
}
```

## Secondary Integration Points

### 1. Compare Command
**File**: `internal/cli/cli.go` - `compareCommand()` function

```go
comparisons := comparer.Compare(oldRun, newRun)
// NEW: Add AI analysis
aiReport := aianalyzer.AnalyzeRegression(oldRun, newRun, comparisons)
fmt.Println(aiReport)
```

### 2. Trend Analysis
**File**: `internal/stats/stats.go` - `AnalyzeTrend()` method

```go
trend := analyzer.AnalyzeTrend(runs, benchmarkName)
// NEW: AI prediction
prediction := aianalyzer.PredictTrend(trend, runs, benchmarkName)
```

### 3. Web Dashboard API
**File**: `internal/dashboard/server.go`

```go
// NEW: Add endpoint
mux.HandleFunc("/api/ai-analysis/", s.handleAIAnalysis)

func (s *Server) handleAIAnalysis(w http.ResponseWriter, r *http.Request) {
    // Return AI analysis as JSON
}
```

## Implementation Strategy

### Step 1: Create AI Analyzer Package
```bash
mkdir -p internal/analyzer
touch internal/analyzer/analyzer.go
```

### Step 2: Define AI Analyzer Interface
```go
package analyzer

import "github.com/alenon/gokanon/internal/models"

type AIAnalyzer interface {
    EnhanceSuggestions(*models.ProfileSummary) []models.Suggestion
    AnalyzeRegression(*models.BenchmarkRun, *models.BenchmarkRun, []models.Comparison) string
    PredictTrend(*models.TrendAnalysis, []models.BenchmarkRun, string) string
}
```

### Step 3: Implement Enhancement Logic
```go
func EnhanceSuggestions(summary *models.ProfileSummary) []models.Suggestion {
    enhanced := make([]models.Suggestion, 0)
    
    for _, sug := range summary.Suggestions {
        // Enhance each suggestion with AI analysis
        enhanced = append(enhanced, models.Suggestion{
            Type:       sug.Type,
            Severity:   sug.Severity,
            Function:   sug.Function,
            Issue:      sug.Issue,
            Suggestion: aiEnhance(sug.Suggestion),
            Impact:     aiPredictImpact(sug),
        })
    }
    
    return enhanced
}
```

### Step 4: Hook into Profiler
**File**: `internal/runner/runner.go` - Around line 253

```go
// After profiler.Analyze()
profileSummary, err := analyzer.Analyze()
if err == nil && profileSummary != nil {
    // NEW: Enhance with AI
    aianalyzer := analyzer.NewAIAnalyzer()
    profileSummary.Suggestions = aianalyzer.EnhanceSuggestions(profileSummary)
}
```

## Data Flow Diagram for AI Integration

```
BenchmarkRun
    ├─ Results[] (metrics)
    ├─ CPUProfile (binary data)
    ├─ MemoryProfile (binary data)
    └─ ProfileSummary
        ├─ CPUTopFunctions[]
        ├─ MemoryTopFunctions[]
        ├─ MemoryLeaks[]
        ├─ HotPaths[]
        └─ Suggestions[]
                ↓
        [AI ANALYZER HERE]
                ↓
        Enhanced Suggestions[]
            ├─ Better descriptions
            ├─ Confidence scores
            ├─ Pattern matching
            ├─ Root cause analysis
            └─ Expected impact estimates
                ↓
        Displayed in:
        ├─ CLI output
        ├─ Web dashboard
        ├─ Exported reports
        └─ JSON storage
```

## Key Information to Access from Benchmarks

### CPU Analysis
- Which functions consume the most CPU time
- Call chains (hot paths)
- Flat vs cumulative time
- Total samples collected

### Memory Analysis
- Which functions allocate the most memory
- Allocation patterns (consistent vs sporadic)
- Total memory allocated
- Potential memory leaks

### Regression Detection
- Percentage changes in performance
- Stability across multiple runs
- Trend direction and slope
- Statistical confidence (R² value)

## Testing the Integration

### Create Test Benchmarks
```bash
gokanon run -pkg=./examples -profile=cpu,mem
```

### Access Generated Data
```bash
# List all runs
gokanon list

# Load a run programmatically
store := storage.NewStorage(".gokanon")
run, _ := store.Load("run-<timestamp>")

// run.ProfileSummary contains the data to analyze
fmt.Printf("%+v\n", run.ProfileSummary.Suggestions)
```

## CLI Output Examples

### Before AI Enhancement
```
PROFILE ANALYSIS
================
🔥 CPU Hot Functions (Total: 12543)
────────────────────────────────────
Function                  Flat%    Cum%
runtime.mallocgc          15.2%    45.8%
mypackage.processData     12.7%    32.1%

💡 Optimization Suggestions
1. 🔴 [CPU] runtime.mallocgc
   Issue: Function consumes 15.2% of CPU time
   Suggestion: Consider profiling in isolation
   Potential Impact: Could improve performance
```

### After AI Enhancement (Target)
```
PROFILE ANALYSIS (AI-ENHANCED)
==============================
🔥 CPU Hot Functions (Total: 12543)
────────────────────────────────────
Function                  Flat%    Cum%
runtime.mallocgc          15.2%    45.8%
mypackage.processData     12.7%    32.1%

💡 AI-Powered Optimization Insights
1. 🔴 [CPU] mypackage.processData
   Issue: High GC pressure detected
   Root Cause: Excessive allocations in loop at line 245
   Similar to: Pattern seen in 3 other functions
   Recommendation: Use sync.Pool for temporary buffers
   Expected Impact: ~8-12% CPU improvement
   Confidence: 94% (based on similar patterns)
   Severity: HIGH
```

## Testing Integration

### Unit Test Template
```go
package analyzer

import (
    "testing"
    "github.com/alenon/gokanon/internal/models"
)

func TestEnhanceSuggestions(t *testing.T) {
    // Create test suggestion
    sug := models.Suggestion{
        Type:       "cpu",
        Severity:   "high",
        Function:   "mypackage.processData",
        Issue:      "Function consumes 12.7% of CPU time",
        Suggestion: "Profile in isolation",
        Impact:     "",
    }
    
    summary := &models.ProfileSummary{
        Suggestions: []models.Suggestion{sug},
    }
    
    // Run AI enhancement
    enhanced := EnhanceSuggestions(summary)
    
    // Verify enhancement
    if len(enhanced) == 0 {
        t.Fatal("Expected enhanced suggestions")
    }
    
    if enhanced[0].Impact == "" {
        t.Error("Expected Impact to be populated")
    }
}
```

## Performance Considerations

1. **Caching**: Cache AI analysis results for repeated runs
2. **Async**: Run expensive analysis asynchronously if needed
3. **Batching**: Process multiple benchmarks efficiently
4. **Storage**: Don't increase JSON size too much

## Minimal Implementation Example

```go
// internal/analyzer/analyzer.go
package analyzer

import (
    "fmt"
    "github.com/alenon/gokanon/internal/models"
)

type AIAnalyzer struct {}

func NewAIAnalyzer() *AIAnalyzer {
    return &AIAnalyzer{}
}

// Enhance suggestions with AI analysis
func (a *AIAnalyzer) EnhanceSuggestions(summary *models.ProfileSummary) []models.Suggestion {
    enhanced := make([]models.Suggestion, 0)
    
    for _, sug := range summary.Suggestions {
        // Add AI-powered analysis
        enhanced = append(enhanced, models.Suggestion{
            Type:       sug.Type,
            Severity:   prioritizeSeverity(sug),
            Function:   sug.Function,
            Issue:      sug.Issue,
            Suggestion: improveRecommendation(sug),
            Impact:     estimateImpact(sug),
        })
    }
    
    return enhanced
}

func prioritizeSeverity(sug models.Suggestion) string {
    // AI logic to rank severity
    return sug.Severity
}

func improveRecommendation(sug models.Suggestion) string {
    // AI logic to enhance recommendation
    return fmt.Sprintf("[AI-Enhanced] %s", sug.Suggestion)
}

func estimateImpact(sug models.Suggestion) string {
    // AI logic to predict impact
    return "Estimated 5-15% improvement"
}
```

## Documentation Links

- **Main Architecture**: See `ARCHITECTURE_AI_INTEGRATION.md`
- **Data Models**: `/home/user/gokanon/internal/models/benchmark.go`
- **Profiler Code**: `/home/user/gokanon/internal/profiler/profiler.go`
- **CLI Commands**: `/home/user/gokanon/internal/cli/cli.go`
- **Storage Layer**: `/home/user/gokanon/internal/storage/storage.go`

## Next Steps

1. Create `internal/analyzer/` package
2. Define AIAnalyzer interface
3. Implement minimal enhancement
4. Hook into profiler.Analyze()
5. Test with sample benchmarks
6. Extend to comparison and trend analysis
7. Add web dashboard API endpoint
8. Create report export with AI insights

