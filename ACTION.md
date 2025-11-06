# GitHub Action for gokanon

This GitHub Action allows you to easily integrate gokanon benchmark testing into your CI/CD workflows. It automatically runs benchmarks, compares them against baselines, and can fail builds if performance degrades beyond acceptable thresholds.

## Features

- **Automated benchmark execution** - Run Go benchmarks on any package
- **Baseline comparison** - Compare against previous runs automatically
- **Regression detection** - Fail builds if performance degrades beyond threshold
- **Beautiful HTML reports** - Interactive reports with charts and graphs powered by Chart.js
- **Multiple export formats** - Generate HTML, CSV, and Markdown reports
- **CPU/Memory profiling** - Enable profiling to identify performance bottlenecks
- **Trend analysis** - Track performance trends over time with statistical insights
- **Statistical analysis** - Analyze stability across multiple runs
- **Artifact uploads** - Automatically upload reports and profiles for easy access
- **PR summaries** - Add benchmark results to workflow summaries

## Quick Start

### Basic Usage

Add this to your workflow file (e.g., `.github/workflows/benchmark.yml`):

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
```

### Advanced Usage with Caching

For comparing against a baseline from the main branch:

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

      - name: Restore baseline benchmarks
        uses: actions/cache@v4
        with:
          path: .gokanon
          key: benchmark-baseline-${{ github.base_ref || github.ref_name }}
          restore-keys: |
            benchmark-baseline-

      - name: Run benchmarks
        id: benchmark
        uses: alenon/gokanon@v1
        with:
          go-version: '1.21'
          packages: './...'
          threshold-percent: 10
          export-format: 'html,markdown'
          fail-on-regression: true

      - name: Save baseline
        if: github.ref == 'refs/heads/main'
        uses: actions/cache/save@v4
        with:
          path: .gokanon
          key: benchmark-baseline-${{ github.ref_name }}-${{ github.sha }}

      - name: Comment PR
        if: github.event_name == 'pull_request' && steps.benchmark.outputs.comparison-summary != ''
        uses: actions/github-script@v8
        with:
          script: |
            const summary = `## 📊 Benchmark Results

            ${{ steps.benchmark.outputs.comparison-summary }}

            - Max degradation: ${{ steps.benchmark.outputs.max-degradation-percent }}%
            - Threshold: ${{ inputs.threshold-percent }}%
            - Status: ${{ steps.benchmark.outputs.passed == 'true' ? '✅ Passed' : '❌ Failed' }}

            [View detailed report](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }})`;

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: summary
            });
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `go-version` | Go version to use | No | `1.21` |
| `packages` | Package paths to benchmark (e.g., `./...`, `./app`) | No | `./...` |
| `benchmark-filter` | Benchmark filter pattern (e.g., `BenchmarkFoo`) | No | `.` (all) |
| `storage-dir` | Directory to store benchmark results | No | `.gokanon` |
| `threshold-percent` | Max degradation % allowed (0-100) | No | `10` |
| `compare-baseline` | Compare against baseline | No | `true` |
| `export-format` | Export format: `html`, `csv`, `markdown`, or comma-separated list | No | `html` |
| `export-output` | Output filename for reports (without extension) | No | `benchmark-report` |
| `stats-last-n` | Number of runs for statistical analysis (0 to skip) | No | `0` |
| `upload-artifact` | Upload reports as workflow artifacts | No | `true` |
| `fail-on-regression` | Fail workflow if threshold exceeded | No | `true` |
| `working-directory` | Working directory for benchmarks | No | `.` |
| `install-gokanon` | Whether to install gokanon | No | `true` |
| `gokanon-version` | Version to install (`latest` or specific tag) | No | `latest` |
| `enable-profiling` | Enable CPU/memory profiling: `cpu`, `mem`, or `cpu,mem` | No | `` (disabled) |
| `trend-analysis-runs` | Number of runs for trend analysis (0 to skip) | No | `0` |

## Outputs

| Output | Description |
|--------|-------------|
| `run-id` | ID of the benchmark run |
| `comparison-summary` | Plain text comparison summary |
| `passed` | Whether threshold check passed (`true`/`false`) |
| `degraded-count` | Number of degraded benchmarks |
| `improved-count` | Number of improved benchmarks |
| `unchanged-count` | Number of unchanged benchmarks |
| `report-path` | Path to generated report file(s) |
| `max-degradation-percent` | Maximum degradation percentage found |

## Usage Examples

### Example 1: Simple Benchmark Testing

Run benchmarks on every push without baseline comparison:

```yaml
name: Simple Benchmarks

on: [push]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5

      - name: Run benchmarks
        uses: alenon/gokanon@v1
        with:
          packages: './...'
          compare-baseline: false
```

### Example 2: Specific Package and Filter

Run benchmarks only for specific packages and functions:

```yaml
- name: Run string benchmarks
  uses: alenon/gokanon@v1
  with:
    packages: './internal/benchmarks'
    benchmark-filter: 'BenchmarkString'
    threshold-percent: 5
```

### Example 3: Multiple Export Formats

Generate reports in multiple formats:

```yaml
- name: Run benchmarks with multiple exports
  uses: alenon/gokanon@v1
  with:
    export-format: 'html,markdown,csv'
    export-output: 'performance-report'
```

### Example 4: Statistical Analysis

Analyze performance trends over last 10 runs:

```yaml
- name: Run benchmarks with statistics
  uses: alenon/gokanon@v1
  with:
    stats-last-n: 10
    threshold-percent: 15
```

### Example 5: Matrix Testing

Test across multiple Go versions:

```yaml
jobs:
  benchmark:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.20', '1.21', '1.22']
    steps:
      - uses: actions/checkout@v5

      - name: Run benchmarks
        uses: alenon/gokanon@v1
        with:
          go-version: ${{ matrix.go-version }}
          storage-dir: '.gokanon-${{ matrix.go-version }}'
```

### Example 6: Conditional Failure

Only fail on severe regressions:

```yaml
- name: Run benchmarks
  uses: alenon/gokanon@v1
  with:
    threshold-percent: 20
    fail-on-regression: true
```

### Example 7: Custom Working Directory

Run benchmarks in a subdirectory:

```yaml
- name: Run benchmarks
  uses: alenon/gokanon@v1
  with:
    working-directory: './services/api'
    packages: './...'
```

### Example 8: Scheduled Performance Testing

Run nightly performance tests:

```yaml
name: Nightly Performance

on:
  schedule:
    - cron: '0 2 * * *' # 2 AM daily

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5

      - name: Restore historical data
        uses: actions/cache@v4
        with:
          path: .gokanon
          key: benchmark-history-${{ github.sha }}
          restore-keys: benchmark-history-

      - name: Run benchmarks
        uses: alenon/gokanon@v1
        with:
          stats-last-n: 30
          export-format: 'html,markdown'
          fail-on-regression: false

      - name: Save historical data
        uses: actions/cache/save@v4
        with:
          path: .gokanon
          key: benchmark-history-${{ github.sha }}
```

### Example 9: Performance Profiling

Run benchmarks with CPU and memory profiling enabled:

```yaml
- name: Run benchmarks with profiling
  uses: alenon/gokanon@v1
  with:
    packages: './...'
    enable-profiling: 'cpu,mem'
    export-format: 'html'
    upload-artifact: true
```

This will:
- Generate CPU and memory profiles during benchmark execution
- Analyze profiles for hot functions and memory issues
- Upload profile files (.prof) as artifacts for detailed analysis
- Include profile insights in the terminal output

### Example 10: Trend Analysis

Analyze performance trends over multiple runs:

```yaml
- name: Run benchmarks with trend analysis
  uses: alenon/gokanon@v1
  with:
    packages: './...'
    trend-analysis-runs: 10
    stats-last-n: 10
    export-format: 'html,markdown'
```

This will:
- Run benchmarks and compare against baseline
- Analyze performance trends over the last 10 runs
- Show improving/degrading/stable trends
- Generate statistical analysis with mean, median, and standard deviation

### Example 11: Complete CI Pipeline with All Features

A comprehensive example using all advanced features:

```yaml
name: Complete Benchmark Pipeline

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

      - name: Restore benchmark history
        uses: actions/cache@v4
        with:
          path: .gokanon
          key: benchmark-${{ github.ref_name }}-${{ github.sha }}
          restore-keys: |
            benchmark-${{ github.ref_name }}-
            benchmark-

      - name: Run comprehensive benchmarks
        id: benchmark
        uses: alenon/gokanon@v1
        with:
          go-version: '1.21'
          packages: './...'
          enable-profiling: 'cpu,mem'
          threshold-percent: 10
          export-format: 'html,markdown,csv'
          stats-last-n: 5
          trend-analysis-runs: 10
          upload-artifact: true
          fail-on-regression: true

      - name: Save baseline
        if: github.ref == 'refs/heads/main'
        uses: actions/cache/save@v4
        with:
          path: .gokanon
          key: benchmark-${{ github.ref_name }}-${{ github.sha }}

      - name: Comment PR with results
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v8
        with:
          script: |
            const summary = ` ## 📊 Benchmark Results

            **Performance Summary:**
            - 🟢 Improved: ${{ steps.benchmark.outputs.improved-count }}
            - 🔴 Degraded: ${{ steps.benchmark.outputs.degraded-count }}
            - ⚪ Unchanged: ${{ steps.benchmark.outputs.unchanged-count }}
            - 📈 Max degradation: ${{ steps.benchmark.outputs.max-degradation-percent }}%

            **Status:** ${{ steps.benchmark.outputs.passed == 'true' ? '✅ Passed' : '❌ Failed' }}

            **Profiling:** CPU and Memory profiling enabled
            **Trend Analysis:** Last 10 runs analyzed

            📄 [View Detailed Report](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }})
            `;

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: summary
            });
```

## Best Practices

### 1. Use Caching for Baseline Comparison

Always cache the `.gokanon` directory to enable proper baseline comparison:

```yaml
- uses: actions/cache@v4
  with:
    path: .gokanon
    key: benchmark-${{ github.ref_name }}-${{ github.sha }}
    restore-keys: benchmark-${{ github.ref_name }}-
```

### 2. Set Appropriate Thresholds

Choose thresholds based on your performance requirements:
- **Strict (5%)**: Critical performance paths
- **Moderate (10%)**: General application code
- **Lenient (20%)**: Non-critical utilities

### 3. Export Multiple Formats

Generate both HTML (for humans) and CSV/Markdown (for automation):

```yaml
export-format: 'html,csv'
```

### 4. Use PR Comments

Add benchmark results to PR comments for better visibility:

```yaml
- uses: actions/github-script@v8
  with:
    script: |
      github.rest.issues.createComment({
        issue_number: context.issue.number,
        owner: context.repo.owner,
        repo: context.repo.repo,
        body: 'Benchmark results: ...'
      });
```

### 5. Separate Baseline and PR Workflows

Run different strategies for main branch (save baseline) vs PRs (compare):

```yaml
- name: Save baseline
  if: github.ref == 'refs/heads/main'
  uses: actions/cache/save@v4
  # ...

- name: Compare against baseline
  if: github.event_name == 'pull_request'
  # ...
```

## Beautiful HTML Reports

The HTML reports generated by gokanon are designed to be extremely user-friendly and visually appealing:

### Features
- **Modern, responsive design** - Works beautifully on desktop, tablet, and mobile
- **Interactive charts** - Powered by Chart.js for visualizing performance data
- **Gradient backgrounds** - Eye-catching purple gradient design
- **Animated cards** - Summary cards with hover effects and animations
- **Color-coded results** - Green for improvements, red for degradations, gray for unchanged
- **Performance comparison charts** - Side-by-side bar charts comparing old vs new results
- **Delta distribution graphs** - Visualize the percentage change for each benchmark
- **Professional typography** - Clean, modern font stack for excellent readability
- **Data badges** - Color-coded badges showing performance deltas
- **Mobile-optimized** - Fully responsive layout that works on all screen sizes

### Report Sections

1. **Header** - Beautiful gradient title with subtitle
2. **Metadata** - Run IDs and timestamps in an elegant card
3. **Summary Cards** - Three prominent cards showing:
   - Number of improved benchmarks (green)
   - Number of degraded benchmarks (red)
   - Number of unchanged benchmarks (gray)
4. **Performance Comparison Chart** - Interactive bar chart showing old vs new performance
5. **Delta Distribution Chart** - Bar chart showing percentage changes for each benchmark
6. **Detailed Table** - Complete data table with all metrics and color-coded indicators
7. **Footer** - Credits with link to gokanon project

### Example Report

When you run benchmarks with `export-format: 'html'`, you'll get a stunning HTML report that includes:

- **Visual performance trends** at a glance
- **Interactive tooltips** on hover for detailed information
- **Smooth animations** for a polished user experience
- **Professional presentation** suitable for sharing with stakeholders

The reports are self-contained HTML files that can be:
- Viewed directly in any browser
- Downloaded from GitHub Actions artifacts
- Shared with team members
- Embedded in documentation
- Archived for historical comparison

## Troubleshooting

### "Not enough runs for comparison"

This happens when there's no baseline to compare against. Solutions:
1. Run the workflow on the main branch first to establish a baseline
2. Set `compare-baseline: false` to skip comparison
3. Use caching to preserve benchmark history

### Benchmarks timing out

If benchmarks take too long:
1. Use `benchmark-filter` to run specific benchmarks
2. Use `packages` to target specific packages
3. Increase workflow timeout: `timeout-minutes: 30`

### Action fails to install gokanon

If installation fails:
1. Check network connectivity
2. Verify `gokanon-version` is valid
3. Set `install-gokanon: false` if pre-installed

## Migration from Manual Scripts

If you're currently using manual gokanon commands in your workflow:

**Before:**
```yaml
- run: go install github.com/alenon/gokanon@latest
- run: gokanon run -pkg=./...
- run: gokanon compare --latest
- run: gokanon check --latest -threshold=10
```

**After:**
```yaml
- uses: alenon/gokanon@v1
  with:
    packages: './...'
    threshold-percent: 10
```

## Contributing

Found a bug or have a feature request? Please open an issue at [github.com/alenon/gokanon](https://github.com/alenon/gokanon).

## License

This action follows the same license as the gokanon project.
