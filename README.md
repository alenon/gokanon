# gokanon

A powerful CLI tool for running and comparing Go benchmark tests. gokanon makes it easy to track performance over time by saving benchmark results, providing detailed comparisons, statistical analysis, and CI/CD integration.

## Features

- **Interactive Web Dashboard**: Real-time visualization with charts, trends, and performance insights
- **Run Benchmarks**: Execute Go benchmarks with custom filters and package paths
- **Save Results**: Automatically save benchmark results with metadata (timestamp, Go version, duration)
- **CPU/Memory Profiling**: Generate and analyze CPU and memory profiles with flame graphs
- **Profile Analysis**: Automatic detection of hot functions, memory leaks, and optimization opportunities
- **Flame Graph Visualization**: Interactive web UI for viewing CPU and memory profiles
- **Optimization Suggestions**: AI-powered recommendations based on profile analysis
- **Compare Results**: Compare two benchmark runs to see performance improvements or degradations
- **Export Reports**: Export comparisons to HTML, CSV, or Markdown formats
- **Statistical Analysis**: Analyze multiple runs with mean, median, standard deviation, and stability metrics
- **Trend Analysis**: Track performance trends over time with linear regression
- **CI/CD Integration**: Automated threshold checking for continuous integration pipelines
- **Track History**: List all saved benchmark results with timestamps
- **Easy Management**: Delete old benchmark results to keep your workspace clean

## Installation

### Quick Install (Recommended)

**macOS & Linux** - One-line install without Go:
```bash
curl -sSL https://raw.githubusercontent.com/alenon/gokanon/main/install.sh | bash
```

**macOS with Homebrew:**
```bash
brew install alenon/tap/gokanon
```

### Manual Installation (Pre-built Binaries)

Download pre-built binaries from [GitHub Releases](https://github.com/alenon/gokanon/releases/latest):

**macOS (Apple Silicon M1/M2/M3):**
```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-darwin-arm64.tar.gz | tar xz
sudo mv gokanon-darwin-arm64 /usr/local/bin/gokanon
```

**macOS (Intel):**
```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-darwin-amd64.tar.gz | tar xz
sudo mv gokanon-darwin-amd64 /usr/local/bin/gokanon
```

**Linux (x86_64):**
```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-linux-amd64.tar.gz | tar xz
sudo mv gokanon-linux-amd64 /usr/local/bin/gokanon
```

**Linux (ARM64):**
```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-linux-arm64.tar.gz | tar xz
sudo mv gokanon-linux-arm64 /usr/local/bin/gokanon
```

**Windows:**
Download `gokanon-windows-amd64.exe.zip` from the [releases page](https://github.com/alenon/gokanon/releases/latest), extract, and add to PATH.

### Install with Go

If you have Go installed:
```bash
go install github.com/alenon/gokanon@latest
```

### Build from Source

```bash
git clone https://github.com/alenon/gokanon.git
cd gokanon
make build
# Binary will be in ./bin/gokanon
```

## GitHub Action

The easiest way to integrate gokanon into your CI/CD pipeline is to use the official GitHub Action:

```yaml
- name: Run benchmarks
  uses: alenon/gokanon@v1
  with:
    packages: './...'
    threshold-percent: 10
    enable-profiling: 'cpu,mem'
    trend-analysis-runs: 10
    export-format: 'html'
```

Features:
- **Zero configuration** - Works out of the box
- **Automatic comparison** - Compares against baseline
- **Beautiful HTML reports** - Interactive charts and graphs with Chart.js
- **CPU/Memory profiling** - Identify performance bottlenecks
- **Trend analysis** - Track performance over time
- **Statistical analysis** - Stability metrics across runs
- **PR integration** - Posts results as PR comments
- **Multiple formats** - Export to HTML, CSV, Markdown
- **Artifact uploads** - Automatically uploads reports and profiles

For detailed documentation, examples, and advanced usage, see [ACTION.md](ACTION.md).

Example workflows are available in [examples/workflows/](examples/workflows/).

## Quick Start

```bash
# Run benchmarks
gokanon run -pkg=./...

# Start interactive web dashboard
gokanon serve

# Run benchmarks with CPU and memory profiling
gokanon run --profile=cpu,mem

# View flame graphs in browser
gokanon flamegraph --latest

# Compare last two runs
gokanon compare --latest

# Export to HTML
gokanon export --latest -format=html

# Check threshold (for CI/CD)
gokanon check --latest -threshold=10
```

## Usage

### Running Benchmarks

Run all benchmarks in the current package:

```bash
gokanon run
```

Run benchmarks with a specific filter:

```bash
gokanon run -bench=BenchmarkStringBuilder
```

Run benchmarks in a specific package or all packages:

```bash
gokanon run -pkg=./examples
gokanon run -pkg=./...
```

Custom storage directory:

```bash
gokanon run -storage=./my-results
```

### CPU and Memory Profiling

gokanon includes powerful profiling capabilities that go beyond basic benchmarking. Generate CPU and memory profiles to identify performance bottlenecks, memory leaks, and optimization opportunities.

#### Running Benchmarks with Profiling

Enable CPU profiling:

```bash
gokanon run --profile=cpu
```

Enable memory profiling:

```bash
gokanon run --profile=mem
```

Enable both CPU and memory profiling:

```bash
gokanon run --profile=cpu,mem
```

When profiling is enabled, gokanon will:
1. Generate pprof profile files during benchmark execution
2. Automatically analyze the profiles
3. Identify hot functions and memory allocation patterns
4. Detect potential memory leaks
5. Provide actionable optimization suggestions

#### Viewing Flame Graphs

After running benchmarks with profiling enabled, you can view interactive visualizations:

```bash
# View profiles for a specific run
gokanon flamegraph run-123

# View profiles for the latest run
gokanon flamegraph --latest
```

This starts a web server (default port 8080) that provides:
- **CPU Flame Graphs**: Visualize where your code spends time
- **Memory Flame Graphs**: See memory allocation patterns
- **Side-by-side Comparison**: Compare CPU and memory profiles
- **Download Options**: Download raw profile files for use with `go tool pprof`

Access the visualization at `http://localhost:8080`

#### Profile Analysis Output

When you run benchmarks with profiling enabled, you'll see detailed analysis:

```
PROFILE ANALYSIS
================================================================================

🔥 CPU Hot Functions (Total samples: 12543)
--------------------------------------------------------------------------------
Function                              Flat%    Cum%
runtime.mallocgc                      15.2%    45.8%
mypackage.processData                 12.7%    32.1%
encoding/json.Marshal                  8.3%    18.9%

💾 Memory Hot Functions (Total: 245.3 MB)
--------------------------------------------------------------------------------
Function                              Flat%    Bytes
mypackage.allocateBuffer              28.5%    69.9 MB
strings.Builder.Grow                  15.2%    37.3 MB

🎯 Hot Execution Paths
--------------------------------------------------------------------------------
1. 32.5% of execution time (4078 samples)
   Critical path consuming 32.5% of execution time
   Path: main.main → processRequests → handleRequest → parseJSON

⚠️  Potential Memory Issues
--------------------------------------------------------------------------------
🔴 mypackage.cacheData (high)
   Allocations: 15234 (125.6 MB)
   Allocated 125.6 MB but much less in use - potential leak

💡 Optimization Suggestions
================================================================================

1. 🔴 [CPU] mypackage.processData
   Issue: Function consumes 12.7% of CPU time
   Suggestion: Consider optimizing this hot function - profile it in isolation
   Potential Impact: Could improve overall performance by up to 8.9%

2. 🔴 [MEMORY] mypackage.allocateBuffer
   Issue: Function allocates 28.5% of total memory
   Suggestion: Consider using sync.Pool or pre-allocate with appropriate capacity
   Potential Impact: Could significantly reduce allocation pressure and GC overhead
```

#### Profile Features

**CPU Profile Analysis:**
- Top CPU-consuming functions (by flat and cumulative time)
- Hot execution paths (critical call chains)
- CPU time distribution

**Memory Profile Analysis:**
- Top memory-allocating functions
- Total memory allocated
- Potential memory leak detection

**Optimization Suggestions:**
- High-severity issues requiring immediate attention
- Medium-severity improvements
- Low-severity optimizations
- Expected impact of each optimization

**Flame Graph Visualization:**
- Interactive web-based viewer
- CPU and memory flame graphs
- Side-by-side comparison view
- Download raw profiles for advanced analysis

#### Advanced Profiling Workflow

For comprehensive performance analysis:

```bash
# 1. Run benchmarks with profiling
gokanon run --profile=cpu,mem

# 2. View the profile summary in terminal (automatic)
# The analysis is displayed immediately after benchmarks complete

# 3. Open flame graphs in browser
gokanon flamegraph --latest

# 4. For advanced analysis, download profiles and use go tool pprof
# The profile paths are shown in the output
go tool pprof -http=:8080 .gokanon/profiles/run-123/cpu.prof
```

#### Integration with go tool pprof

All profiles are stored in `.gokanon/profiles/<run-id>/` and are compatible with standard Go profiling tools:

```bash
# Interactive terminal UI
go tool pprof .gokanon/profiles/run-123/cpu.prof

# Generate flame graph
go tool pprof -http=:8080 .gokanon/profiles/run-123/cpu.prof

# Generate call graph
go tool pprof -pdf .gokanon/profiles/run-123/cpu.prof > callgraph.pdf

# Top functions
go tool pprof -top .gokanon/profiles/run-123/cpu.prof
```

### AI-Powered Analysis

gokanon includes an AI analyzer that provides intelligent insights and optimization suggestions for your benchmark results. The AI analyzer uses free AI services to understand your profiling data and comparison results, offering actionable recommendations.

#### Quick Start

Enable AI analysis with a single environment variable:

```bash
# Enable AI analysis
export GOKANON_AI_ENABLED=true

# Run benchmarks with profiling
gokanon run --profile=cpu,mem

# Compare runs with AI insights
gokanon compare --latest
```

#### Supported AI Providers

**Ollama (Recommended for local use)**
- Completely free, no API keys required
- Private - your data stays on your machine
- Install: https://ollama.ai/

```bash
# Setup Ollama
ollama pull llama3.2
ollama serve &

# Enable in gokanon
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=ollama
```

**OpenAI (GPT-4o, GPT-4-turbo)**
- State-of-the-art performance
- Requires API key from https://platform.openai.com/api-keys

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai
export GOKANON_AI_API_KEY=sk-your-openai-key
export GOKANON_AI_MODEL=gpt-4o
```

**Anthropic Claude (Sonnet 4.5, Haiku 4.5)**
- Excellent reasoning and analysis
- Requires API key from https://console.anthropic.com

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=anthropic
export GOKANON_AI_API_KEY=sk-ant-your-key
export GOKANON_AI_MODEL=claude-sonnet-4-5-20250929
```

**Google Gemini (Gemini 2.5 Flash, 2.0 Flash)**
- Very affordable and fast
- Free tier available through Google AI Studio
- Requires API key from https://aistudio.google.com/apikey

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=gemini
export GOKANON_AI_API_KEY=your-google-api-key
export GOKANON_AI_MODEL=gemini-2.5-flash
```

**Groq (Fast cloud inference)**
- Very fast inference with free tier
- Requires API key from https://console.groq.com

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=groq
export GOKANON_AI_API_KEY=your-api-key
```

**OpenAI-Compatible (LM Studio, LocalAI, Cursor, etc.)**
- Works with any OpenAI-compatible API
- Supports local and custom deployments

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai-compatible
export GOKANON_AI_BASE_URL=http://localhost:1234/v1
export GOKANON_AI_MODEL=local-model
```

#### What the AI Analyzer Does

- **Enhanced Profile Analysis**: AI examines CPU hotspots, memory allocations, and hot paths to provide deeper insights than rule-based analysis
- **Intelligent Suggestions**: Get specific, actionable optimization recommendations with context about why they matter
- **Comparison Insights**: Understand performance changes between runs with explanations of likely causes
- **Pattern Recognition**: AI identifies common performance anti-patterns and suggests proven solutions

For detailed configuration and usage, see [AI_ANALYZER.md](AI_ANALYZER.md).

### Interactive Web Dashboard

The interactive web dashboard provides a comprehensive, real-time visualization interface for all your benchmark data. It's the easiest way to analyze trends, compare runs, and share results with your team.

#### Starting the Dashboard

Start the dashboard server:

```bash
# Start on default port (8080)
gokanon serve

# Start on custom port
gokanon serve -port=9000

# Allow remote access
gokanon serve -addr=0.0.0.0
```

Then open your browser to `http://localhost:8080`

#### Dashboard Features

**Overview Tab:**
- Real-time statistics (total runs, tests, benchmarks)
- Recent performance trends chart
- Quick access to recent benchmark runs

**Trends Tab:**
- Historical performance graphs with Chart.js
- Multi-line charts showing performance over time
- Statistical analysis (mean, median, std dev, CV)
- Automatic trend detection (improving/degrading/stable)
- Filter by specific benchmarks
- Customizable time range (10, 25, 50, or 100 runs)

**History Tab:**
- Complete benchmark run history
- Sortable and filterable table
- Search by package name or run ID
- Click any run to view details

**Compare Tab:**
- Side-by-side comparison of any two runs
- Performance deltas with color-coded indicators
- Percentage improvements and degradations

**Additional Features:**
- 🌙 Dark mode with persistent theme preference
- 🔍 Global search across all runs and benchmarks
- 🔗 Shareable URLs for specific runs
- 📱 Responsive design for mobile and tablet
- 🎨 Beautiful, modern UI with smooth animations

#### Embed Mode

Embed the dashboard in your documentation:

```html
<iframe
  src="http://localhost:8080?embed=true&tab=trends"
  width="100%"
  height="600px">
</iframe>
```

For detailed documentation, see [docs/DASHBOARD.md](docs/DASHBOARD.md).

### Listing Results

List all saved benchmark results:

```bash
gokanon list
```

Output example:
```
ID              Timestamp            Benchmarks  Duration      Package
--              ---------            ----------  --------      -------
run-1699123456  2024-11-04 15:30:56  10          2.5s         ./examples
run-1699123123  2024-11-04 15:25:23  10          2.3s         ./examples
```

### Comparing Results

Compare two specific benchmark runs:

```bash
gokanon compare run-1699123123 run-1699123456
```

Compare the last two runs:

```bash
gokanon compare --latest
```

Output example:
```
Comparing: run-1699123123 (2024-11-04 15:25:23) vs run-1699123456 (2024-11-04 15:30:56)

✓ StringBuilder                              12345.67 ns/op →    11234.56 ns/op (-9.00%)
✗ StringConcatenation                        98765.43 ns/op →   102345.67 ns/op (+3.63%)
~ StringJoin                                 45678.90 ns/op →    45912.34 ns/op (+0.51%)

Summary: 1 improved, 1 degraded, 1 unchanged
```

Status indicators:
- `✓` Improved performance (lower is better)
- `✗` Degraded performance
- `~` No significant change (within 5% threshold)

### Exporting Results

Export comparisons to various formats:

**HTML Report** (with beautiful styling):
```bash
gokanon export --latest -format=html -output=report.html
```

**CSV** (for spreadsheet analysis):
```bash
gokanon export --latest -format=csv -output=results.csv
```

**Markdown** (for documentation):
```bash
gokanon export --latest -format=markdown -output=comparison.md
```

Compare specific runs:
```bash
gokanon export run-123 run-456 -format=html
```

### Statistical Analysis

Analyze multiple benchmark runs to understand stability and variation:

```bash
# Analyze all runs
gokanon stats

# Analyze last 5 runs
gokanon stats -last=5

# Custom stability threshold
gokanon stats -last=10 -cv-threshold=15
```

Output example:
```
Statistical Analysis (5 runs)
Runs: 2024-11-04 10:00:00 to 2024-11-04 15:30:00

Benchmark Statistics:
----------------------------------------------------------------------------------------------------
StringBuilder          Count:   5 | Mean:    362.45 ns/op | Median:    363.20 ns/op | StdDev:     4.12 (±1.1%) | Range: [358.30 - 367.50] ✓ Stable
StringConcatenation    Count:   5 | Mean:   5234.67 ns/op | Median:   5198.40 ns/op | StdDev:   234.89 (±4.5%) | Range: [4986.20 - 5543.10] ✓ Stable
```

Metrics:
- **Count**: Number of benchmark runs analyzed
- **Mean**: Average performance across runs
- **Median**: Middle value (less affected by outliers)
- **StdDev**: Standard deviation (variation measure)
- **CV**: Coefficient of variation (% - lower is more stable)
- **Range**: Min and max values observed

### Trend Analysis

Track performance trends over time:

```bash
# Analyze last 10 runs
gokanon trend -last=10

# Analyze specific benchmark
gokanon trend -benchmark=BenchmarkStringBuilder -last=20
```

Output example:
```
Performance Trend Analysis (10 runs)
Period: 2024-11-01 10:00:00 to 2024-11-04 15:30:00

Benchmark: StringBuilder
  🟢 Trend: improving ↓ (slope: -2.34 ns/op per run)
  Confidence: 87.3% (R²)
  Data points: 370.25 → 365.12 (-1.4%) → 362.45 (-0.7%) → 359.87 (-0.7%) ...

Benchmark: StringConcatenation
  🔴 Trend: degrading ↑ (slope: +12.45 ns/op per run)
  Confidence: 92.1% (R²)
  Data points: 5123.45 → 5198.23 (+1.5%) → 5267.89 (+1.3%) → 5334.12 (+1.3%) ...

Benchmark: StringJoin
  ⚪ Trend: stable → (slope: -0.34 ns/op per run)
  Confidence: 45.2% (R²)
  Data points: 879.12 → 881.34 (+0.3%) → 878.56 (-0.3%) → 880.23 (+0.2%) ...
```

Indicators:
- **🟢 improving**: Performance getting better over time
- **🔴 degrading**: Performance getting worse over time
- **⚪ stable**: No significant trend
- **Confidence**: R² value (higher = more reliable trend)

### Threshold Checking (CI/CD)

Check if performance degraded beyond acceptable limits:

```bash
# Check with 10% threshold
gokanon check --latest -threshold=10

# Check specific runs
gokanon check run-123 run-456 -threshold=5

# Strict threshold for critical code
gokanon check --latest -threshold=2
```

Output example:
```
Threshold Check (max degradation: 10.0%)
Comparing: run-123 vs run-456

✗ 2/10 benchmarks failed the threshold check:

  • StringConcatenation: Performance degraded by 15.30% (threshold: 10.00%)
  • MapIteration: Performance degraded by 12.45% (threshold: 10.00%)
```

**Exit codes**:
- `0`: All benchmarks passed threshold
- `1`: One or more benchmarks failed threshold (useful for CI/CD)

### Deleting Results

Delete a specific benchmark run:

```bash
gokanon delete run-1699123456
```

## Example Benchmarks

The repository includes example benchmarks in the `examples/` directory that demonstrate common Go performance patterns:

**String Operations** (`examples/string_test.go`):
- String concatenation methods comparison
- String formatting benchmarks

**Slice and Map Operations** (`examples/slice_test.go`):
- Slice append with and without pre-allocation
- Slice copying
- Map access and iteration

Run the example benchmarks:

```bash
gokanon run -pkg=./examples
```

## CI/CD Integration

gokanon is designed for seamless CI/CD integration.

### GitHub Actions (Recommended)

Use the official GitHub Action for the easiest setup:

```yaml
name: Benchmarks

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5

      - name: Run benchmarks
        uses: alenon/gokanon@v1
        with:
          packages: './...'
          threshold-percent: 10
          export-format: 'html,markdown'
```

See [ACTION.md](ACTION.md) for complete documentation and [examples/workflows/](examples/workflows/) for ready-to-use workflow examples.

### Manual GitHub Actions Setup

If you prefer to use the CLI directly:

```yaml
name: Benchmark Check

on:
  pull_request:
    branches: [ main ]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v5
    - uses: actions/setup-go@v6
      with:
        go-version: '1.21'

    - name: Install gokanon
      run: go install github.com/alenon/gokanon@latest

    - name: Run benchmarks
      run: gokanon run -pkg=./...

    - name: Check threshold
      run: gokanon check --latest -threshold=10

    - name: Export HTML report
      if: always()
      run: gokanon export --latest -format=html -output=report.html

    - name: Upload report
      if: always()
      uses: actions/upload-artifact@v5
      with:
        name: benchmark-report
        path: report.html
```

See `.github/workflows/benchmark.yml` for a complete example with baseline comparison.

### Supported CI Platforms

- GitHub Actions
- GitLab CI
- Jenkins
- CircleCI
- Any CI platform with Go support

## Storage

By default, benchmark results are stored in the `.gokanon` directory in JSON format. Each run is saved with a unique ID based on the timestamp.

You can specify a custom storage directory using the `-storage` flag with any command:

```bash
gokanon run -storage=./benchmark-results
gokanon list -storage=./benchmark-results
gokanon compare --latest -storage=./benchmark-results
```

## Benchmark Result Structure

Each saved benchmark run includes:

- **ID**: Unique identifier (e.g., `run-1699123456`)
- **Timestamp**: When the benchmark was run
- **Go Version**: Version of Go used to run the benchmark
- **Duration**: Total time taken to run all benchmarks
- **Command**: The exact command used to run the benchmarks
- **Results**: Array of individual benchmark results containing:
  - Name
  - Iterations
  - Nanoseconds per operation (ns/op)
  - Bytes per operation (B/op)
  - Allocations per operation (allocs/op)
  - MB/s (if applicable)

## Use Cases

### Performance Regression Testing

Run benchmarks before and after code changes to ensure performance doesn't degrade:

```bash
# Before changes
gokanon run

# Make your code changes...

# After changes
gokanon run

# Compare
gokanon compare --latest

# Check if degradation is acceptable
gokanon check --latest -threshold=5
```

### Tracking Performance Over Time

Build a history of benchmark results to track performance trends:

```bash
# Run benchmarks regularly (e.g., daily or per commit)
gokanon run

# Review statistics
gokanon stats -last=30

# Analyze trends
gokanon trend -last=30

# Review all results
gokanon list
```

### A/B Testing Optimizations

Compare different optimization approaches:

```bash
# Test approach A
gokanon run
# Note the run ID (e.g., run-123)

# Change to approach B
# Modify code...
gokanon run
# Note the run ID (e.g., run-456)

# Compare
gokanon compare run-123 run-456

# Export detailed report
gokanon export run-123 run-456 -format=html
```

### Continuous Integration

Automatically fail builds if performance degrades:

```bash
# In CI pipeline
gokanon run -pkg=./...

# Fail if degradation > 10%
gokanon check --latest -threshold=10 || exit 1

# Export report for artifacts
gokanon export --latest -format=html -output=ci-report.html
```

## Advanced Features

### Multiple Run Analysis

For more reliable results, run benchmarks multiple times and analyze:

```bash
# Run benchmarks 5 times
for i in {1..5}; do
  gokanon run -pkg=./...
  sleep 10  # Cool down period
done

# Analyze stability
gokanon stats -last=5

# Check trends
gokanon trend -last=5
```

### Export to Multiple Formats

Generate reports in different formats for different audiences:

```bash
# For management (HTML with visualizations)
gokanon export --latest -format=html -output=report.html

# For analysis (CSV for Excel/Google Sheets)
gokanon export --latest -format=csv -output=data.csv

# For documentation (Markdown for GitHub)
gokanon export --latest -format=markdown -output=BENCHMARKS.md
```

### Benchmark-Specific Analysis

Focus on specific benchmarks:

```bash
# Run specific benchmark
gokanon run -bench=BenchmarkCriticalPath

# Analyze trend for specific benchmark
gokanon trend -benchmark=BenchmarkCriticalPath -last=20
```

## Tips

1. **Consistent Environment**: Run benchmarks in consistent environments (same hardware, similar system load) for accurate comparisons
2. **Multiple Runs**: Consider running benchmarks multiple times and comparing averages for more reliable results
3. **Baseline**: Keep a baseline benchmark run to compare all future optimizations against
4. **CI Integration**: Integrate gokanon into your CI pipeline to catch performance regressions early
5. **Thresholds**: Set appropriate thresholds based on your application's performance requirements
6. **Cool Down**: When running multiple benchmarks in sequence, add cool-down periods to prevent thermal throttling
7. **Statistics**: Use `gokanon stats` to understand benchmark stability before making optimization decisions

## Troubleshooting

**No benchmarks found**:
- Ensure your test files are named `*_test.go`
- Ensure benchmark functions start with `Benchmark`
- Check that you're in the correct directory or using the right `-pkg` flag

**High variation in results**:
- Use `gokanon stats` to check coefficient of variation
- Increase benchmark time: `go test -bench=. -benchtime=10s`
- Ensure system is not under heavy load
- Close unnecessary applications

**Permission denied**:
- The tool needs write access to the storage directory (default: `.gokanon`)
- Use a different storage directory with `-storage` flag if needed

**Comparison shows no results**:
- The benchmark names must match exactly between runs
- Ensure both run IDs exist (use `gokanon list` to verify)

**Threshold check always fails**:
- Check if your threshold is too strict for your benchmark variability
- Use `gokanon stats` to understand typical variation
- Consider running benchmarks multiple times and using median values

## Development & Testing

### Running Tests

GoKanon has comprehensive test coverage with beautiful, easy-to-understand test reports.

**Quick test run:**
```bash
go test ./...
```

**Generate beautiful test report:**
```bash
./scripts/test-report.sh
```

**Save report to file:**
```bash
./scripts/test-report.sh test-report.txt
```

The test report provides:
- 🎨 **Color-coded output** with visual indicators
- 📊 **Progress bars** for coverage percentages
- 🎯 **A-F grading** for overall coverage
- 📦 **Package breakdown** with detailed metrics
- 🔍 **Command analysis** for each GoKanon command
- 💡 **Actionable recommendations** for improvement

**Test coverage:**
- Overall: ~69% and growing
- Commands: 62.7% (34.5% → 62.7% improvement)
- Core packages: 80%+ coverage

See `TEST_COVERAGE_REPORT.md` for detailed coverage analysis.

### CI/CD Testing

Tests run automatically on all pushes and pull requests via GitHub Actions.

The workflow:
- Runs all tests
- Generates coverage reports
- Posts results as PR comments
- Uploads test artifacts
- Enforces quality gates

See `.github/workflows/test-report.yml` for configuration.

### End-to-End Verification

Run the complete verification suite:

```bash
./test-verify-all-commands.sh
```

This script tests all 14 commands end-to-end with real benchmarks.

### Building

```bash
# Build binary
go build -o gokanon .

# Run locally
./gokanon --help
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

Before submitting a PR:
1. Run tests: `./scripts/test-report.sh`
2. Ensure coverage doesn't decrease
3. Add tests for new features
4. Update documentation

## License

See LICENSE file for details.

## Documentation

- [GitHub Action Documentation](ACTION.md) - Complete guide for the GitHub Action
- [Example Workflows](examples/workflows/) - Ready-to-use workflow examples
- [CI/CD Integration Guide](docs/CI_INTEGRATION.md) - Detailed CI/CD setup instructions
- [GitHub Actions Example](.github/workflows/benchmark.yml) - Complete workflow example
- [Homebrew Setup Guide](docs/HOMEBREW.md) - How to set up and maintain the Homebrew tap
- [Release Process](docs/RELEASE.md) - Complete guide for creating releases

## Roadmap

Future enhancements being considered:
- ✅ Export to HTML, CSV, Markdown
- ✅ Statistical analysis
- ✅ Trend analysis
- ✅ CI/CD integration
- ✅ CPU/Memory profiling integration
- ✅ Flame graph visualization
- ✅ Memory leak detection
- ✅ Hot path identification
- ✅ Profile-guided optimization suggestions
- 📋 Benchmark comparison graphs
- 📋 Profile comparison between runs
- 📋 Slack/Discord notifications for regressions
- 📋 Integration with benchstat
- 📋 Enhanced web UI for result visualization

Have a feature request? [Open an issue](https://github.com/alenon/gokanon/issues)!
