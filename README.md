# gokanon

A powerful CLI tool for running and comparing Go benchmark tests with profiling, analysis, and CI/CD integration.

## Features

- **Run & Save Benchmarks** - Execute and automatically save benchmark results
- **CPU/Memory Profiling** - Generate profiles with flame graph visualization
- **Compare Results** - Compare benchmark runs with detailed performance analysis
- **Statistical Analysis** - Analyze multiple runs with mean, median, and stability metrics
- **Trend Analysis** - Track performance trends over time with regression detection
- **Interactive Dashboard** - Web-based visualization with charts and insights
- **Export Reports** - Generate reports in HTML, CSV, or Markdown
- **CI/CD Integration** - Automated threshold checking for pipelines
- **AI Analysis** - Get intelligent optimization suggestions (Ollama, OpenAI, Claude, Gemini)

## Installation

### Quick Install

**macOS & Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/alenon/gokanon/main/install.sh | bash
```

**macOS with Homebrew:**
```bash
brew install alenon/tap/gokanon
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/alenon/gokanon/releases/latest)

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-darwin-arm64.tar.gz | tar xz
sudo mv gokanon-darwin-arm64 /usr/local/bin/gokanon
```

**Linux (x86_64):**
```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-linux-amd64.tar.gz | tar xz
sudo mv gokanon-linux-amd64 /usr/local/bin/gokanon
```

### Install with Go

```bash
go install github.com/alenon/gokanon@latest
```

## Quick Start

```bash
# Run benchmarks
gokanon run -pkg=./...

# View interactive dashboard
gokanon serve

# Run with profiling
gokanon run --profile=cpu,mem

# View flame graphs
gokanon flamegraph --latest

# Compare results
gokanon compare --latest

# Export to HTML
gokanon export --latest -format=html

# Check performance threshold (CI/CD)
gokanon check --latest -threshold=10
```

## Usage

### Running Benchmarks

```bash
# All benchmarks in current package
gokanon run

# Specific benchmark pattern
gokanon run -bench=BenchmarkStringBuilder

# All packages
gokanon run -pkg=./...

# With profiling
gokanon run --profile=cpu,mem
```

### Profiling & Analysis

Generate CPU and memory profiles to identify bottlenecks:

```bash
# Enable profiling
gokanon run --profile=cpu,mem

# View flame graphs
gokanon flamegraph --latest
```

The profiler automatically:
- Identifies hot functions and memory allocation patterns
- Detects potential memory leaks
- Provides optimization suggestions with impact analysis

### AI-Powered Analysis

Enable AI analysis for intelligent insights:

```bash
# Using Ollama (free, local)
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=ollama

# Using OpenAI
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai
export GOKANON_AI_API_KEY=sk-your-key
export GOKANON_AI_MODEL=gpt-4o

# Using Anthropic Claude
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=anthropic
export GOKANON_AI_API_KEY=sk-ant-your-key

# Run with AI analysis
gokanon run --profile=cpu,mem
gokanon compare --latest
```

Supported providers: Ollama, OpenAI, Anthropic, Gemini, Groq, OpenAI-compatible APIs

### Interactive Dashboard

```bash
# Start dashboard
gokanon serve

# Custom port
gokanon serve -port=9000
```

Access at `http://localhost:8080` for:
- Real-time performance trends
- Historical data with charts
- Side-by-side comparisons
- Dark mode support

### Comparing Results

```bash
# Last two runs
gokanon compare --latest

# Specific runs
gokanon compare run-123 run-456

# Against baseline
gokanon compare --baseline=v1.0
```

### Statistical & Trend Analysis

```bash
# Analyze last 5 runs
gokanon stats -last=5

# Track performance trends
gokanon trend -last=10
```

### Exporting Reports

```bash
# HTML report
gokanon export --latest -format=html -output=report.html

# CSV for analysis
gokanon export --latest -format=csv -output=results.csv

# Markdown for docs
gokanon export --latest -format=markdown -output=comparison.md
```

### CI/CD Integration

```bash
# Fail if degradation > 10%
gokanon check --latest -threshold=10
```

**GitHub Action:**
```yaml
- name: Run benchmarks
  uses: alenon/gokanon@v1
  with:
    packages: './...'
    threshold-percent: 10
    enable-profiling: 'cpu,mem'
    export-format: 'html'
```

See `action.yml` for complete GitHub Action configuration.

### Managing Results

```bash
# List all results
gokanon list

# Delete a run
gokanon delete run-123

# Manage baselines
gokanon baseline save -name=v1.0
gokanon baseline list
gokanon baseline show -name=v1.0
```

## Commands

```bash
gokanon run          # Run benchmarks and save results
gokanon list         # List saved benchmark results
gokanon compare      # Compare benchmark results
gokanon export       # Export comparison to HTML/CSV/Markdown
gokanon stats        # Statistical analysis of multiple runs
gokanon trend        # Performance trend analysis
gokanon check        # Check thresholds (CI/CD)
gokanon flamegraph   # View CPU/memory flame graphs
gokanon serve        # Start interactive dashboard
gokanon delete       # Delete benchmark results
gokanon baseline     # Manage baselines (save, list, show, delete)
gokanon doctor       # Run diagnostics
gokanon interactive  # Interactive mode with auto-completion
gokanon completion   # Install shell completion
gokanon version      # Show version information
gokanon help         # Show help
```

## Storage

Results are stored in `.gokanon` directory by default. Use `-storage` flag to customize:

```bash
gokanon run -storage=./benchmark-results
gokanon list -storage=./benchmark-results
```

## Best Practices

1. **Consistent Environment** - Run benchmarks on consistent hardware and system load
2. **Multiple Runs** - Run benchmarks multiple times for reliable results
3. **Baseline** - Maintain a baseline for comparing optimizations
4. **CI Integration** - Catch regressions early with automated checks
5. **Appropriate Thresholds** - Set thresholds based on application requirements

## Development

```bash
# Build
make build

# Run tests
make test

# Generate coverage
make coverage

# Build for all platforms
make build-all
```

## License

See LICENSE file for details.

## Contributing

Contributions welcome! Please submit issues or pull requests at [github.com/alenon/gokanon](https://github.com/alenon/gokanon).

### For Maintainers: Releasing a New Version

See [RELEASE.md](RELEASE.md) for the complete release process. Quick summary:

```bash
# Bump version
make version-bump-minor  # or patch/major

# Commit version
git add VERSION
git commit -m "Bump version to $(cat VERSION)"
git push origin main

# Create and push tag
make tag-release
git push origin v$(cat VERSION)
```

The release workflow will automatically build binaries, generate release notes with changes, and create a GitHub release.
