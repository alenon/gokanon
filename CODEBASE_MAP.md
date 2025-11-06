# GoKanon Codebase Map

## Directory Structure & File Purposes

```
gokanon/
│
├── main.go (16 lines)
│   └─> Entry point: calls cli.Execute()
│
├── internal/
│   │
│   ├── cli/ (1,061 lines)
│   │   └── cli.go
│   │       ├─> Execute(): Main command router
│   │       ├─> runCommand(): Execute benchmarks
│   │       ├─> compareCommand(): Compare two runs
│   │       ├─> exportCommand(): Export to HTML/CSV/Markdown
│   │       ├─> statsCommand(): Multi-run statistics
│   │       ├─> trendCommand(): Trend analysis
│   │       ├─> checkCommand(): CI/CD threshold checking
│   │       ├─> flamegraphCommand(): Profile visualization
│   │       ├─> serveCommand(): Web dashboard
│   │       ├─> deleteCommand(): Remove runs
│   │       ├─> doctorCommand(): System diagnostics
│   │       └─> interactiveCommand(): REPL mode
│   │
│   ├── models/ (85 lines) **[AI DATA SOURCE]**
│   │   └── benchmark.go
│   │       ├─> BenchmarkRun: Complete benchmark results
│   │       ├─> BenchmarkResult: Individual benchmark
│   │       ├─> Comparison: Delta between runs
│   │       ├─> ProfileSummary: Profile analysis data
│   │       ├─> FunctionProfile: Function metrics
│   │       ├─> MemoryLeak: Memory issue
│   │       ├─> HotPath: Critical call chain
│   │       └─> Suggestion: Optimization recommendation
│   │
│   ├── runner/ (300+ lines) **[BENCHMARK EXECUTOR]**
│   │   └── runner.go
│   │       ├─> Runner: Benchmark execution engine
│   │       ├─> Run(): Execute go test with profiling
│   │       ├─> parseOutput(): Parse benchmark results
│   │       ├─> handleProfiles(): **[PRIMARY AI HOOK]**
│   │       │   ├─> Profiler.Analyze(): Generate ProfileSummary
│   │       │   ├─> Save profiles to storage
│   │       │   └─> **Where to inject AI enhancement**
│   │       └─> getGoVersion(): Get Go version
│   │
│   ├── profiler/ (250+ lines) **[PROFILE ANALYZER]**
│   │   └── profiler.go
│   │       ├─> Analyzer: Pprof profile processor
│   │       ├─> LoadCPUProfile(): Parse CPU profile
│   │       ├─> LoadMemoryProfile(): Parse memory profile
│   │       ├─> Analyze(): Generate ProfileSummary with:
│   │       │   ├─> CPUTopFunctions
│   │       │   ├─> MemoryTopFunctions
│   │       │   ├─> MemoryLeaks
│   │       │   ├─> HotPaths
│   │       │   └─> Basic Suggestions (rule-based)
│   │       ├─> analyzeCPUProfile()
│   │       ├─> analyzeMemoryProfile()
│   │       ├─> identifyHotPaths()
│   │       ├─> detectMemoryLeaks()
│   │       └─> generateSuggestions(): Basic suggestions
│   │
│   ├── storage/ (200+ lines) **[DATA PERSISTENCE]**
│   │   └── storage.go
│   │       ├─> Storage: File-based persistence
│   │       ├─> Save(): Write BenchmarkRun to JSON
│   │       ├─> Load(): Load run by ID
│   │       ├─> List(): Get all runs (sorted)
│   │       ├─> GetLatest(): Most recent run
│   │       ├─> SaveProfile(): Save profile file
│   │       ├─> Delete(): Remove run
│   │       └─> Data format: .gokanon/{run-id}.json
│   │
│   ├── compare/ (109 lines) **[COMPARISON LOGIC]**
│   │   └── compare.go
│   │       ├─> Comparer: Benchmark comparison
│   │       ├─> Compare(): Compare two runs
│   │       ├─> compareResults(): Compare individual benchmarks
│   │       ├─> Threshold: 5% for "same" classification
│   │       ├─> Status: "improved", "degraded", "same"
│   │       ├─> FormatComparison(): Format for display
│   │       └─> Summary(): Generate summary text
│   │
│   ├── stats/ (204 lines) **[STATISTICAL ANALYSIS]**
│   │   └── stats.go
│   │       ├─> Analyzer: Multi-run statistics
│   │       ├─> AnalyzeMultiple(): Calculate stats across runs
│   │       ├─> Stats: Mean, median, stddev, CV
│   │       ├─> IsStable(): Check stability (CV threshold)
│   │       ├─> AnalyzeTrend(): Linear regression analysis
│   │       ├─> TrendAnalysis: Direction, slope, R²
│   │       ├─> linearRegression(): Calculate trend line
│   │       └─> FormatStats(): Format for display
│   │
│   ├── export/ (250+ lines) **[REPORT GENERATION]**
│   │   └── export.go
│   │       ├─> Exporter: Multi-format export
│   │       ├─> ToHTML(): Generate HTML report
│   │       ├─> ToCSV(): Generate CSV report
│   │       ├─> ToMarkdown(): Generate Markdown
│   │       ├─> HTML: Styled comparison tables
│   │       ├─> CSV: Spreadsheet-compatible
│   │       └─> Markdown: Documentation-friendly
│   │
│   ├── threshold/ (86 lines) **[CI/CD VALIDATION]**
│   │   └── threshold.go
│   │       ├─> Checker: Threshold validation
│   │       ├─> Check(): Validate against threshold
│   │       ├─> Result: Pass/fail with failures list
│   │       ├─> Failure: Individual benchmark failure
│   │       ├─> FormatResult(): Format for display
│   │       └─> ExitCode(): Exit code for CI/CD
│   │
│   ├── dashboard/ (500+ lines) **[WEB UI]**
│   │   ├── server.go: HTTP server & API
│   │   │   ├─> Server: Dashboard server
│   │   │   ├─> Start(): Start HTTP server
│   │   │   ├─> handleRuns(): List runs API
│   │   │   ├─> handleRunDetail(): Run details API
│   │   │   ├─> handleTrends(): Trends API
│   │   │   ├─> handleStats(): Stats API
│   │   │   ├─> handleSearch(): Search API
│   │   │   ├─> handleIndex(): Serve HTML
│   │   │   └─> handleStatic(): Serve static files
│   │   ├── app.go: JavaScript frontend (embedded)
│   │   ├── frontend.go: HTML template
│   │   └── Features:
│   │       ├─ Overview tab
│   │       ├─ Trends tab
│   │       ├─ History tab
│   │       ├─ Compare tab
│   │       ├─ Dark/Light mode
│   │       ├─ Search
│   │       └─ Responsive design
│   │
│   ├── webserver/ (200+ lines) **[FLAME GRAPH VIEWER]**
│   │   └── server.go
│   │       ├─> Server: Flame graph web server
│   │       ├─> Start(): Start with pprof
│   │       ├─> handleIndex(): HTML page
│   │       └─> Interactive profile visualization
│   │
│   ├── ui/ (150+ lines) **[TERMINAL UI]**
│   │   ├── colors.go: Terminal colors
│   │   ├── errors.go: Error formatting
│   │   ├── progress.go: Progress bar
│   │   └── Helper functions for CLI output
│   │
│   ├── doctor/ (100+ lines) **[DIAGNOSTICS]**
│   │   └── doctor.go
│   │       ├─> RunDiagnostics(): Check system setup
│   │       └─> PrintResults(): Display results
│   │
│   ├── interactive/ (150+ lines) **[REPL MODE]**
│   │   └── interactive.go
│   │       ├─> Session: Interactive mode
│   │       ├─> RegisterCommand(): Register handlers
│   │       └─> Run(): Start REPL loop
│   │
│   └── threshold/ (see above)
│
└── examples/
    ├── string_test.go: String benchmark examples
    └── slice_test.go: Slice benchmark examples
```

## Data Flow Relationships

### 1. Benchmark Execution Flow
```
main.go
  └─> cli.Execute()
      └─> runCommand()
          ├─> runner.Run()
          │   ├─> exec go test
          │   ├─> parseOutput()
          │   ├─> handleProfiles()
          │   │   └─> profiler.Analyze()
          │   │       ├─> analyzeCPUProfile()
          │   │       ├─> analyzeMemoryProfile()
          │   │       ├─> identifyHotPaths()
          │   │       ├─> detectMemoryLeaks()
          │   │       └─> generateSuggestions()
          │   └─> return BenchmarkRun
          │
          ├─> storage.Save(run)
          │   └─> Write .gokanon/{run-id}.json
          │
          └─> displayProfileSummary()
              └─> Print to CLI
```

### 2. Comparison Flow
```
compareCommand()
  ├─> storage.Load(oldID)
  ├─> storage.Load(newID)
  ├─> comparer.Compare()
  │   └─> return Comparison[]
  └─> Print comparisons
```

### 3. Analysis Flow
```
statsCommand() or trendCommand()
  ├─> storage.List()
  ├─> analyzer.AnalyzeMultiple() or AnalyzeTrend()
  └─> Print results
```

### 4. Export Flow
```
exportCommand()
  ├─> Load two runs
  ├─> comparer.Compare()
  ├─> exporter.ToHTML/CSV/Markdown()
  └─> Write file
```

### 5. Dashboard Flow
```
serveCommand()
  ├─> dashboard.NewServer()
  ├─> Server.Start() (HTTP server)
  └─> API endpoints:
      ├─> /api/runs -> storage.List()
      ├─> /api/runs/{id} -> storage.Load()
      ├─> /api/trends -> stats.AnalyzeTrend()
      └─> /api/stats -> stats.AnalyzeMultiple()
```

## Key Integration Points

### Highest Priority: Profile Enhancement
**Location**: `runner/runner.go:handleProfiles()`
**Current**: Calls `profiler.Analyzer.Analyze()` → generates `ProfileSummary`
**Opportunity**: Enhance Suggestions after analysis

```
profiler.Analyze() → ProfileSummary
                 ↓
              [AI HERE] ← INJECT ENHANCEMENT
                 ↓
          Enhanced Suggestions
```

### Secondary: Comparison Analysis
**Location**: `cli/cli.go:compareCommand()`
**Current**: Displays comparison results
**Opportunity**: Add AI insights about regression

```
comparer.Compare() → Comparison[]
                ↓
            [AI HERE] ← ANALYZE REGRESSION
                ↓
        AI Insights + CLI output
```

### Tertiary: Trend Prediction
**Location**: `stats/stats.go:AnalyzeTrend()`
**Current**: Linear regression only
**Opportunity**: ML-based forecasting

```
runs[] → AnalyzeTrend() → TrendAnalysis
                       ↓
                  [AI HERE] ← PREDICT FUTURE
                       ↓
            Prediction + Confidence
```

## File Dependencies

```
cli.go
├── runner.go (import)
├── storage.go (import)
├── compare.go (import)
├── export.go (import)
├── stats.go (import)
├── dashboard.go (import)
├── threshold.go (import)
└── profiler.go (indirectly via runner)

runner.go
├── models.go (import)
├── storage.go (import)
└── profiler.go (import)

profiler.go
├── models.go (import)
└── pprof library (external)

storage.go
└── models.go (import)

compare.go
└── models.go (import)

stats.go
└── models.go (import)

export.go
└── models.go (import)

dashboard.go
├── storage.go (import)
└── models.go (indirectly)
```

## Data Structure Relationships

```
BenchmarkRun
├─ ID: string
├─ Timestamp: time.Time
├─ Results: []BenchmarkResult
│   └─ Name, Iterations, NsPerOp, BytesPerOp, AllocsPerOp
├─ CPUProfile: string (path)
├─ MemoryProfile: string (path)
└─ ProfileSummary: *ProfileSummary
    ├─ CPUTopFunctions: []FunctionProfile
    ├─ MemoryTopFunctions: []FunctionProfile
    ├─ MemoryLeaks: []MemoryLeak
    ├─ HotPaths: []HotPath
    └─ Suggestions: []Suggestion ← **AI ENHANCEMENT TARGET**

Comparison
├─ Name: string
├─ OldNsPerOp: float64
├─ NewNsPerOp: float64
├─ Delta: float64
├─ DeltaPercent: float64
└─ Status: string

TrendAnalysis
├─ BenchmarkName: string
├─ Direction: string
├─ TrendLine: float64
└─ Confidence: float64
```

## Testing Entry Points

### Run Benchmarks with Profiling
```bash
gokanon run -pkg=./examples -profile=cpu,mem
```

### Access Generated Data
```go
store := storage.NewStorage(".gokanon")
run, _ := store.Load("run-<id>")
suggestions := run.ProfileSummary.Suggestions
```

### Test Integration Points

1. **Profile Analysis**: `internal/profiler/profiler_test.go`
2. **Comparison Logic**: `internal/compare/compare_test.go`
3. **Stats Analysis**: `internal/stats/stats_test.go`
4. **Storage**: `internal/storage/storage_test.go`

## Performance Characteristics

- **Benchmark Parsing**: ~1-2ms (regex-based)
- **Profile Analysis**: ~100-500ms (depends on profile size)
- **Storage I/O**: ~10-50ms per run
- **Dashboard**: Real-time queries from JSON files
- **Statistics**: O(n) for n runs

## Key Metrics to Leverage

### CPU Profile
- Function name
- Flat percentage (time in function)
- Cumulative percentage (time + callees)
- Sample count

### Memory Profile
- Function name
- Allocation count
- Bytes allocated
- Severity assessment

### Benchmark Results
- ns/op (primary metric)
- B/op (memory allocation)
- allocs/op (allocation count)
- Trend slope (ns/op per run)

## Extension Points for AI

1. **Suggestion Enhancement**: Improve text of suggestions
2. **Pattern Recognition**: Identify common bottleneck patterns
3. **Severity Ranking**: Rerank suggestions by impact
4. **Correlation Analysis**: Link CPU and memory issues
5. **Prediction**: Forecast future performance
6. **Report Generation**: Create AI-powered reports

