# GoKanon Architecture & Integration Analysis

## Executive Summary

**GoKanon** is a powerful CLI benchmarking tool for Go programs (~13,236 lines of Go code) that captures, stores, compares, and visualizes benchmark results. It features profiling, statistical analysis, trend tracking, and a web dashboard. The codebase is well-structured with clear separation of concerns, making it highly extensible for AI analyzer integration.

---

## 1. MAIN ENTRY POINT AND CLI STRUCTURE

### Entry Point: `/home/user/gokanon/main.go`
```go
func main() {
    if err := cli.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### CLI Router: `/home/user/gokanon/internal/cli/cli.go`
The `cli.Execute()` function implements a command-based router with 14 main commands:

| Command | Purpose | Output |
|---------|---------|--------|
| **run** | Execute benchmarks with optional profiling | Saves BenchmarkRun as JSON |
| **list** | Display all saved runs | Tabular listing |
| **compare** | Compare two runs side-by-side | Comparison metrics |
| **export** | Export comparisons (HTML/CSV/Markdown) | Files in specified format |
| **stats** | Statistical analysis across multiple runs | Mean, median, stddev, CV |
| **trend** | Linear regression analysis over time | Direction, slope, confidence |
| **check** | CI/CD threshold validation | Pass/fail with exit codes |
| **flamegraph** | View CPU/memory flame graphs | Web UI on port 8080 |
| **serve** | Start interactive web dashboard | Web server on port 8080 |
| **delete** | Remove benchmark runs | Confirmation message |
| **doctor** | Run diagnostics | System checks |
| **interactive** | Interactive mode with autocomplete | REPL-style interface |
| **completion** | Shell completion scripts | Bash/Zsh/Fish |

---

## 2. DATA MODELS & BENCHMARK RESULTS

### Core Data Structures: `/home/user/gokanon/internal/models/benchmark.go`

#### BenchmarkRun (Complete Result Set)
```go
type BenchmarkRun struct {
    ID             string                          // Unique ID (run-{unix_timestamp})
    Timestamp      time.Time                       // When run executed
    Package        string                          // Package path tested
    GoVersion      string                          // Go compiler version
    Results        []BenchmarkResult               // Individual benchmark results
    Command        string                          // The "go test" command executed
    Duration       time.Duration                   // Total execution time
    CPUProfile     string                          // Path to CPU profile file
    MemoryProfile  string                          // Path to memory profile file
    ProfileSummary *ProfileSummary                 // AI-analyzable profile data
}
```

#### BenchmarkResult (Individual Benchmark)
```go
type BenchmarkResult struct {
    Name        string          // Benchmark name (e.g., "BenchmarkFoo")
    Iterations  int64           // Number of iterations run
    NsPerOp     float64         // Nanoseconds per operation (PRIMARY METRIC)
    BytesPerOp  int64           // Bytes allocated per operation
    AllocsPerOp int64           // Number of allocations per operation
    MBPerSec    float64         // Throughput in MB/s (for I/O benchmarks)
}
```

#### Comparison (Delta Between Runs)
```go
type Comparison struct {
    Name         string          // Benchmark name
    OldNsPerOp   float64         // Previous performance
    NewNsPerOp   float64         // Current performance
    Delta        float64         // Absolute difference (ns)
    DeltaPercent float64         // Percentage change
    Status       string          // "improved" | "degraded" | "same" (>5% threshold)
}
```

#### ProfileSummary (AI-Ready Data)
```go
type ProfileSummary struct {
    CPUTopFunctions    []FunctionProfile       // Top CPU consumers
    MemoryTopFunctions []FunctionProfile       // Top memory allocators
    MemoryLeaks        []MemoryLeak            // Potential memory issues
    HotPaths           []HotPath               // Critical execution paths
    Suggestions        []Suggestion            // Optimization recommendations
    TotalCPUSamples    int64                   // Total CPU samples collected
    TotalMemoryBytes   int64                   // Total memory allocated
}

type FunctionProfile struct {
    Name        string              // Function name (e.g., "runtime.mallocgc")
    FlatPercent float64             // Time in function itself
    CumPercent  float64             // Time in function + callees
    FlatValue   int64               // Actual samples or bytes
    CumValue    int64               // Cumulative value
}

type Suggestion struct {
    Type       string              // "cpu" | "memory" | "algorithm"
    Severity   string              // "low" | "medium" | "high"
    Function   string              // Affected function
    Issue      string              // What's the problem
    Suggestion string              // How to fix it
    Impact     string              // Expected improvement
}
```

---

## 3. HOW BENCHMARKS ARE RUN & RESULTS PROCESSED

### Benchmark Execution Flow

```
CLI Command: gokanon run -pkg=./examples -profile=cpu,mem
    ↓
RunCommand() in cli.go
    ├─ Parse flags (-bench, -pkg, -storage, -profile)
    ├─ Parse profiling options (CPU, Memory)
    ├─ Create Runner instance
    └─ Enable profiling if specified
        ↓
Runner.Run() in runner/runner.go
    ├─ Get Go version: exec.Command("go", "version")
    ├─ Generate unique ID: "run-{unix_timestamp}"
    ├─ Create temp directory for profiles
    ├─ Build "go test" command:
    │   - go test -bench {filter} -benchmem
    │   - (optional) -cpuprofile {path}
    │   - (optional) -memprofile {path}
    │   - {package_path}
    ├─ Execute: cmd.Run()
    ├─ Parse output with regex:
    │   BenchmarkFoo-8  1000000  1234 ns/op  512 B/op  10 allocs/op
    ├─ Create BenchmarkRun with results
    └─ If profiling enabled:
        ├─ Save profile files
        ├─ Load profiles into Analyzer
        ├─ Generate ProfileSummary
        └─ Identify hot functions & memory issues
            ↓
Storage.Save() in storage/storage.go
    └─ Write to: .gokanon/{run_id}.json
       Stores complete BenchmarkRun as JSON
```

### Parsing Benchmark Output

The regex pattern in `runner.go` (line 135):
```
^Benchmark(\S+)\s+(\d+)\s+([\d.]+)\s+ns/op(?:\s+([\d.]+)\s+MB/s)?(?:\s+(\d+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?
```

Extracts:
1. Benchmark name (after "Benchmark" prefix)
2. Iterations count
3. ns/op value (mandatory)
4. MB/s value (optional)
5. B/op value (optional)
6. allocs/op value (optional)

---

## 4. RESULT STORAGE & DATA FLOW

### Storage Structure: `/home/user/gokanon/internal/storage/storage.go`

```
.gokanon/
├── run-1699123456.json          # JSON file with all benchmark data
├── run-1699123123.json
├── run-1699123000.json
└── profiles/
    ├── run-1699123456/
    │   ├── cpu.prof              # CPU profile (pprof format)
    │   └── mem.prof              # Memory profile (pprof format)
    └── run-1699123123/
        ├── cpu.prof
        └── mem.prof
```

### Storage API
```go
type Storage struct {
    dir string  // Directory path (default: ".gokanon")
}

// Key methods:
Save(run *BenchmarkRun)          // Save single run as JSON
Load(id string)                  // Load run by ID
List()                           // Get all runs (sorted by timestamp, newest first)
GetLatest()                       // Get most recent run
SaveProfile(id, type, reader)    // Save profile file
GetCPUProfilePath(id)             // Get path to CPU profile
GetMemoryProfilePath(id)          // Get path to memory profile
Delete(id)                        // Remove run directory
```

### Data Persistence
- **Format**: JSON (human-readable, version-controllable)
- **Location**: `.gokanon/` directory (customizable via `-storage` flag)
- **Sorting**: By timestamp, newest first
- **Profiles**: Stored as binary pprof format alongside JSON metadata

---

## 5. OUTPUT FORMATS & REPORTING

### Display Mechanisms

#### 1. Terminal Output (CLI)
**Compare Command**:
```
Comparing: run-123 (2024-11-04 15:25:23) vs run-456 (2024-11-04 15:30:56)

✓ StringBuilder                      12345.67 ns/op →    11234.56 ns/op (-9.00%)
✗ StringConcatenation                98765.43 ns/op →   102345.67 ns/op (+3.63%)
~ StringJoin                         45678.90 ns/op →    45912.34 ns/op (+0.51%)

Summary: 1 improved, 1 degraded, 1 unchanged
```

**Stats Command**:
```
StringBuilder      Count:   5 | Mean:  362.45 ns/op | Median:  363.20 ns/op | StdDev: 4.12 (±1.1%) | Range: [358.30 - 367.50] ✓ Stable
```

**Trend Command**:
```
Benchmark: StringBuilder
  🟢 Trend: improving ↓ (slope: -2.34 ns/op per run)
  Confidence: 87.3% (R²)
  Data points: 370.25 → 365.12 (-1.4%) → 362.45 (-0.7%) → 359.87 (-0.7%) ...
```

#### 2. Export Formats

**HTML Report** (export/export.go):
- Styled comparison tables
- Color-coded status (✓ ✗ ~)
- Summary statistics
- Beautiful CSS styling

**CSV Format**:
```csv
Benchmark,Old (ns/op),New (ns/op),Delta (ns/op),Delta (%),Status
StringBuilder,12345.67,11234.56,-1111.11,-9.00,improved
```

**Markdown Format**:
```markdown
# Benchmark Comparison
Comparing: `run-123` vs `run-456`

| Status | Benchmark | Old (ns/op) | New (ns/op) | Delta | Delta (%) |
|--------|-----------|-------------|-------------|-------|-----------|
| 🟢 | StringBuilder | 12345.67 | 11234.56 | -1111.11 | -9.00% |
```

#### 3. Interactive Web Dashboard
**File**: `/home/user/gokanon/internal/dashboard/`

**Features**:
- Real-time visualization with Chart.js
- 5 main tabs: Overview, Trends, History, Compare, Profile Analysis
- Dark/Light mode toggle
- Search functionality
- Responsive design
- Embed mode for documentation

**API Endpoints**:
```
GET  /api/runs           - List all benchmark runs
GET  /api/runs/{id}      - Get run details
GET  /api/trends         - Get trend analysis data
GET  /api/stats          - Get statistical analysis
GET  /api/search         - Search benchmarks
GET  /                   - HTML dashboard frontend
```

#### 4. Flame Graph Viewer
**File**: `/home/user/gokanon/internal/webserver/server.go`

- Interactive visualization of CPU/Memory profiles
- Uses pprof Go library for parsing
- Web UI for exploring call stacks

---

## 6. ANALYSIS MODULES

### Statistical Analysis: `/home/user/gokanon/internal/stats/stats.go`

```go
type Stats struct {
    Name     string      // Benchmark name
    Count    int         // Number of runs
    Mean     float64     // Average performance
    Median   float64     // Middle value
    Min      float64     // Best performance
    Max      float64     // Worst performance
    StdDev   float64     // Standard deviation
    Variance float64     // Variance
    CV       float64     // Coefficient of Variation (StdDev/Mean * 100)
}

// Trend Analysis
type TrendAnalysis struct {
    BenchmarkName string      // Which benchmark
    Direction     string      // "improving" | "degrading" | "stable"
    TrendLine     float64     // Slope (ns/op per run)
    Confidence    float64     // R² value (0-1)
}
```

**Algorithms**:
- Linear regression for trend detection
- Coefficient of Variation for stability assessment
- Min/Max tracking for variance analysis

### Profile Analysis: `/home/user/gokanon/internal/profiler/profiler.go`

```go
type Analyzer struct {
    cpuProfile    *profile.Profile    // Parsed CPU profile
    memoryProfile *profile.Profile    // Parsed memory profile
}

// Methods:
LoadCPUProfile(data []byte)         // Parse CPU profile
LoadMemoryProfile(data []byte)      // Parse memory profile
Analyze() *ProfileSummary           // Generate complete analysis
analyzeCPUProfile()                 // Extract top functions
analyzeMemoryProfile()              // Find memory allocators
identifyHotPaths()                  // Find critical call chains
detectMemoryLeaks()                 // Identify allocation issues
generateSuggestions()               // Create optimization recommendations
```

### Comparison Logic: `/home/user/gokanon/internal/compare/compare.go`

```go
type Comparer struct {
    threshold float64     // Default 5% for "same" classification
}

// Classifies changes as:
// - Improved: DeltaPercent < -5%
// - Degraded: DeltaPercent > 5%
// - Same: -5% ≤ DeltaPercent ≤ 5%
```

### Threshold Checking: `/home/user/gokanon/internal/threshold/threshold.go`

```go
type Checker struct {
    maxDegradation float64     // Max allowed degradation % (for CI/CD)
}

// Used in CI/CD pipelines:
// Exit code 0: All benchmarks within threshold
// Exit code 1: Any benchmark exceeded threshold
```

---

## 7. ARCHITECTURE DIAGRAM

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI INTERFACE                            │
│  (cli.go: Execute() routes to 14 commands)                      │
└────────────┬────────────────────────────────────────────────────┘
             │
    ┌────────┴──────────┬──────────────┬──────────────┬──────────┐
    │                   │              │              │          │
    ▼                   ▼              ▼              ▼          ▼
┌───────────┐      ┌─────────┐   ┌──────────┐  ┌────────┐  ┌───────┐
│  RUN      │      │ COMPARE │   │ EXPORT   │  │ STATS  │  │PROFILE│
│  COMMAND  │      │ COMMAND │   │ COMMAND  │  │COMMAND │  │ANALYSIS│
└──────┬────┘      └────┬────┘   └────┬─────┘  └───┬────┘  └───┬───┘
       │                │              │           │            │
       ▼                ▼              ▼           ▼            ▼
   ┌───────────────────────────────────────────────────────────────┐
   │              CORE DATA PROCESSING LAYER                        │
   ├───────────────────────────────────────────────────────────────┤
   │ runner/       compare/      export/      stats/     profiler/ │
   │ runner.go     compare.go    export.go    stats.go   profiler  │
   │              • Go test      • HTML        • Multi-run• Pprof   │
   │              • Parsing      • CSV         • Trends   • Profile │
   │              • Profiling    • Markdown    • Linear   • Hot func│
   │              • Threshold    • Tables      • Regression• Leaks │
   └────────┬──────────────┬─────────────────┬──────────┬─────────┘
            │              │                 │          │
            └──────────────┼─────────────────┼──────────┘
                           │                 │
                           ▼                 ▼
                    ┌──────────────────────────────┐
                    │   STORAGE LAYER              │
                    │   storage/storage.go         │
                    ├──────────────────────────────┤
                    │ .gokanon/                    │
                    │  ├─ run-{id}.json (data)    │
                    │  └─ profiles/{id}/          │
                    │     ├─ cpu.prof             │
                    │     └─ mem.prof             │
                    └──────────────────────────────┘
                           │
            ┌──────────────┴──────────────┐
            ▼                             ▼
    ┌────────────────┐         ┌──────────────────┐
    │ CLI Display    │         │ Web Interfaces   │
    ├────────────────┤         ├──────────────────┤
    │• Tables        │         │ • Dashboard      │
    │• Status marks  │         │   (serve)        │
    │• Suggestions   │         │ • Flame Graphs   │
    │• Trends        │         │   (flamegraph)   │
    └────────────────┘         │ • API endpoints  │
                               └──────────────────┘
```

---

## 8. CURRENT DISPLAY MECHANISMS

### Where Results Are Currently Shown

1. **Terminal Output** (immediate feedback)
   - Location: `cli.go` command handlers
   - Shows: Formatted tables with status symbols
   - Format: Fixed-width text with emoji indicators

2. **JSON Files** (persistent storage)
   - Location: `.gokanon/{run-id}.json`
   - Contains: Complete benchmark metadata and ProfileSummary
   - Use: Loading historical data, comparisons

3. **Web Dashboard** (interactive)
   - Location: `dashboard/` package
   - Endpoint: `http://localhost:8080` (via `gokanon serve`)
   - Features: Charts, trends, search, filtering

4. **Exported Reports** (shareable)
   - Formats: HTML, CSV, Markdown
   - Location: User-specified file
   - Use: Sharing results, archiving, team reviews

5. **Profiling Reports** (performance analysis)
   - CLI: Terminal output with tables and suggestions
   - Web: Flame graphs at `http://localhost:8080` (via `gokanon flamegraph`)

---

## 9. KEY INTEGRATION POINTS FOR AI ANALYZER

### A. ProfileSummary Injection Point (HIGHEST PRIORITY)
**Location**: `runner/runner.go:handleProfiles()` → `profiler/profiler.go:Analyze()`

**Current Flow**:
```
Profile files → Analyzer.Analyze() → ProfileSummary 
              (CPU/Memory top functions, hot paths, memory leaks)
              → Stored in BenchmarkRun.ProfileSummary
              → Displayed in CLI via displayProfileSummary()
```

**Integration Point**: After `profiler.Analyze()` generates ProfileSummary, **inject AI analysis**:
```go
// In runner.go:handleProfiles() around line 250-270
profileSummary, err := analyzer.Analyze()
if err == nil {
    // NEW: AI Enhancement Point
    aiAnalysis := aianalyzer.EnhanceProfileAnalysis(profileSummary)
    profileSummary.Suggestions = aiAnalysis.EnhancedSuggestions
    run.ProfileSummary = profileSummary
}
```

**Why this is ideal**:
- Has access to parsed profile data
- ProfileSummary is already a structured format
- Output (Suggestions) is already expected and displayed
- No need to modify storage or CLI

### B. Comparison Analysis Point (SECONDARY)
**Location**: `cli.go:compareCommand()` → `compare/compare.go:Compare()`

**Enhancement Opportunity**:
```go
// In compareCommand() after line 331
comparisons := comparer.Compare(oldRun, newRun)
aiAnalysis := aianalyzer.AnalyzeRegression(oldRun, newRun, comparisons)
fmt.Println(aiAnalysis.RegressionInsights)
```

**Data Available**:
- Two complete BenchmarkRun objects
- Comparison results with delta percentages
- Historical context (timestamp, Go version, etc.)
- Profile summaries if profiling was enabled

### C. Trend Analysis Enhancement (SECONDARY)
**Location**: `stats/stats.go:AnalyzeTrend()`

**Enhancement**:
```go
trend := analyzer.AnalyzeTrend(runs, name)
aiInsights := aianalyzer.PredictFuturePerformance(trend, runs, name)
fmt.Printf("AI Prediction: %s\n", aiInsights)
```

**Context Available**:
- Historical trend data
- Linear regression results (slope, R² confidence)
- All historical values

### D. Web Dashboard Integration (TERTIARY)
**Location**: `dashboard/server.go` API endpoints

**Enhancement**:
```go
// Add new endpoint: /api/ai-analysis/{run_id}
mux.HandleFunc("/api/ai-analysis/", s.handleAIAnalysis)

func (s *Server) handleAIAnalysis(w http.ResponseWriter, r *http.Request) {
    runID := extractRunID(r.URL.Path)
    run, _ := s.storage.Load(runID)
    analysis := aianalyzer.AnalyzeRun(run)
    json.NewEncoder(w).Encode(analysis)
}
```

**Frontend Changes**:
- New "AI Analysis" tab in dashboard
- Display enhanced suggestions
- Show predictions and insights
- Trend forecasting visualizations

### E. Export Format Enhancement (TERTIARY)
**Location**: `export/export.go`

**New Export Option**:
```go
// Add new method
func (e *Exporter) ToAIReport(run *BenchmarkRun, filename string) error {
    analysis := aianalyzer.GenerateDetailedReport(run)
    // Render to markdown/HTML with AI insights
    return writeReport(filename, analysis)
}
```

---

## 10. DATA FLOW FOR AI INTEGRATION

### Complete Data Journey

```
┌─────────────────────────────────────┐
│  Benchmark Execution                │
│  go test -bench -benchmem -cpuprofile
└──────────────────┬──────────────────┘
                   │
                   ▼
        ┌─────────────────────┐
        │ Parse Output        │
        │ Extract metrics     │
        └──────────┬──────────┘
                   │
                   ▼
        ┌──────────────────────────┐
        │ BenchmarkRun Created     │
        │ (results populated)      │
        └──────────┬───────────────┘
                   │
                   ├─────────────────┬──────────────────┐
                   │                 │                  │
        ┌──────────▼────────┐  ┌─────▼───────┐  ┌─────▼─────┐
        │ Parse Profiles    │  │ Save to     │  │ Display   │
        │ (CPU, Memory)     │  │ Storage     │  │ in CLI    │
        │                   │  │ .gokanon/   │  │           │
        └─────────┬─────────┘  └─────────────┘  └───────────┘
                  │
                  ▼
        ┌─────────────────────┐
        │ Profiler.Analyze()  │
        │ Extract:            │
        │ - Top CPU functions │
        │ - Memory allocators │
        │ - Hot paths         │
        │ - Memory leaks      │
        └──────────┬──────────┘
                   │
                   ▼
        ┌──────────────────────────────┐
        │ ProfileSummary Generated     │
        │ (Suggestions: basic rules)   │
        │                              │
        │ 🔴 [CPU] foo                 │
        │    Issue: High CPU usage     │
        │    Suggestion: Profile...    │
        └──────────┬───────────────────┘
                   │
    ┌──────────────▼──────────────────┐
    │                                  │
    │  *** AI ANALYZER INTEGRATION ***  │
    │                                  │
    │  Enhanced Analysis:               │
    │  - Semantic understanding        │
    │  - Cross-reference suggestions   │
    │  - Pattern recognition           │
    │  - Prioritization ranking        │
    │  - Detailed explanations         │
    │                                  │
    └──────────────┬──────────────────┘
                   │
                   ▼
        ┌──────────────────────────────┐
        │ Enhanced ProfileSummary      │
        │ Suggestions: AI-improved     │
        │                              │
        │ 🔴 [CPU] foo                 │
        │    Issue: High CPU usage     │
        │    Root Cause: Detected      │
        │    Pattern: Similar to...    │
        │    Suggestion: Specific...   │
        │    Confidence: 94%           │
        └──────────┬───────────────────┘
                   │
         ┌─────────┼─────────┐
         │         │         │
         ▼         ▼         ▼
    ┌────────┐┌────────┐┌─────────┐
    │ Save   ││Display ││Dashboard│
    │to JSON ││in CLI  ││Export   │
    │storage ││with    ││Reports  │
    │        ││colors  ││& API    │
    └────────┘└────────┘└─────────┘
```

---

## 11. SUGGESTED AI ANALYZER FEATURES

Based on the architecture analysis, here are high-impact AI features:

### Phase 1: Low-Hanging Fruit
1. **Enhanced Suggestion Generation**
   - Semantic enrichment of basic suggestions
   - Context-aware recommendations
   - Confidence scoring

2. **Root Cause Analysis**
   - Pattern matching across profile data
   - Correlation with code patterns
   - Historical comparison

### Phase 2: Medium Effort
3. **Regression Detection**
   - Anomaly detection in trends
   - Statistical significance testing
   - Automated alerts

4. **Predictive Analysis**
   - Performance forecasting
   - Trend direction prediction
   - Degradation warnings

### Phase 3: Advanced
5. **Code Pattern Recognition**
   - Common bottleneck patterns
   - Optimization templates
   - Best practice recommendations

6. **Multi-Run Intelligence**
   - Cross-run pattern matching
   - Macro trend analysis
   - System-wide optimization recommendations

---

## 12. IMPLEMENTATION CHECKLIST FOR AI ANALYZER

### Prerequisites
- [ ] Understand ProfileSummary structure (see models.go)
- [ ] Review comparison.go for delta logic
- [ ] Study stats.go for trend calculations
- [ ] Test with sample benchmark runs

### Implementation Phases

**Phase 1: Data Ingestion**
- [ ] Create `internal/analyzer/` package
- [ ] Define AIAnalysis struct
- [ ] Implement parsing of BenchmarkRun
- [ ] Validate against sample data

**Phase 2: Integration Points**
- [ ] Hook into profiler.Analyze() output
- [ ] Add to compare command
- [ ] Extend stats module
- [ ] API endpoint in dashboard

**Phase 3: Analysis Features**
- [ ] Profile analysis enhancement
- [ ] Comparison insights
- [ ] Trend prediction
- [ ] Report generation

**Phase 4: Testing & Refinement**
- [ ] Unit tests for analyzer
- [ ] Integration tests with CLI
- [ ] Performance benchmarks
- [ ] User acceptance testing

---

## 13. KEY FILES REFERENCE

| File | Lines | Purpose |
|------|-------|---------|
| `main.go` | 16 | Entry point |
| `internal/cli/cli.go` | 1,061 | CLI router & commands |
| `internal/models/benchmark.go` | 85 | Data structures |
| `internal/runner/runner.go` | 300+ | Benchmark execution |
| `internal/storage/storage.go` | 200+ | JSON persistence |
| `internal/compare/compare.go` | 109 | Comparison logic |
| `internal/profiler/profiler.go` | 250+ | Profile analysis |
| `internal/stats/stats.go` | 204 | Statistical analysis |
| `internal/export/export.go` | 250+ | Export formats |
| `internal/dashboard/` | 500+ | Web UI & APIs |
| `internal/threshold/threshold.go` | 86 | CI/CD validation |

---

## 14. SUMMARY TABLE

| Aspect | Current Implementation | AI Integration Opportunity |
|--------|---|---|
| **Data Collection** | Benchmark running & parsing | N/A |
| **Storage** | JSON files | Use as input source |
| **Basic Analysis** | Compare, stats, trends | Enhance with semantics |
| **Profiling** | Pprof parsing, rule-based suggestions | AI pattern recognition |
| **Reporting** | HTML, CSV, Markdown | AI-powered insights |
| **Display** | CLI tables, web dashboard | AI visualization |
| **Prediction** | Linear regression only | ML-based forecasting |

