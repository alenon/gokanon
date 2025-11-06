# Example Workflows

This directory contains example GitHub Actions workflows demonstrating various ways to use the gokanon action in your CI/CD pipeline.

## Available Examples

### 1. Basic Workflow (`basic.yml`)

**Use case**: Simple benchmark testing on every push and PR.

**Features**:
- Runs benchmarks on all packages
- Compares against previous run
- Fails if performance degrades more than 10%
- Generates HTML report

**When to use**: Getting started with gokanon, simple projects.

---

### 2. Advanced with Caching (`advanced-with-caching.yml`)

**Use case**: Production-ready workflow with baseline comparison and PR integration.

**Features**:
- Compares PR benchmarks against main branch baseline
- Caches benchmark history
- Posts results as PR comments
- Uploads HTML and Markdown reports
- Detailed result interpretation

**When to use**: Production projects, teams that want PR visibility.

---

### 3. Scheduled Performance Testing (`scheduled-performance.yml`)

**Use case**: Nightly performance monitoring and trend analysis.

**Features**:
- Runs on a schedule (daily at 2 AM UTC)
- Analyzes last 30 runs for trends
- Creates issues on significant degradation
- Stores long-term performance data
- Doesn't fail on regressions (monitoring only)
- Exports multiple formats (HTML, CSV, Markdown)

**When to use**: Long-term performance monitoring, detecting gradual degradation.

---

### 4. Matrix Testing (`matrix-go-versions.yml`)

**Use case**: Testing performance across multiple Go versions.

**Features**:
- Tests on Go 1.20, 1.21, and 1.22
- Separate baseline for each version
- Individual reports per version
- Summary job that aggregates results

**When to use**: Libraries that support multiple Go versions, ensuring consistent performance.

---

## How to Use These Examples

1. **Choose an example** that matches your use case
2. **Copy the workflow file** to your repository's `.github/workflows/` directory
3. **Customize the settings**:
   - Update branch names if needed
   - Adjust `packages` to match your project structure
   - Set appropriate `threshold-percent` for your requirements
   - Modify `go-version` if needed
4. **Commit and push** to trigger the workflow

## Quick Setup

```bash
# Create workflows directory in your repo
mkdir -p .github/workflows

# Copy the example you want
cp examples/workflows/basic.yml .github/workflows/benchmark.yml

# Customize as needed
vim .github/workflows/benchmark.yml

# Commit and push
git add .github/workflows/benchmark.yml
git commit -m "Add benchmark workflow"
git push
```

## Customization Tips

### Adjusting Thresholds

Choose based on your performance requirements:
- **5%**: Strict, for critical performance paths
- **10%**: Moderate, for general application code (default)
- **15-20%**: Lenient, for less critical code or noisy benchmarks

### Targeting Specific Benchmarks

Use filters to run only specific benchmarks:

```yaml
with:
  benchmark-filter: 'BenchmarkCriticalPath'
  packages: './internal/core'
```

### Multiple Export Formats

Generate reports in multiple formats:

```yaml
with:
  export-format: 'html,markdown,csv'
```

### Adjusting Schedule

For scheduled workflows, modify the cron expression:

```yaml
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM
    - cron: '0 14 * * 1' # Weekly on Monday at 2 PM
```

## Common Patterns

### Pattern 1: PR Comparison

Compare PR performance against main branch baseline:

```yaml
- name: Restore baseline
  uses: actions/cache@v4
  with:
    path: .gokanon
    key: benchmark-baseline-main
```

### Pattern 2: Historical Tracking

Keep long-term benchmark history:

```yaml
- uses: actions/cache@v4
  with:
    path: .gokanon
    key: benchmark-history-${{ github.sha }}
    restore-keys: benchmark-history-
```

### Pattern 3: Conditional Failure

Only fail on severe regressions:

```yaml
with:
  threshold-percent: 20
  fail-on-regression: true
```

## Troubleshooting

### No Baseline Available

**Problem**: First run or cache miss shows "No baseline available for comparison"

**Solution**: This is normal for the first run. Subsequent runs will have a baseline.

### Permission Denied for PR Comments

**Problem**: Cannot post PR comments

**Solution**: Add permissions to your workflow:

```yaml
permissions:
  contents: read
  pull-requests: write
```

### Benchmarks Take Too Long

**Problem**: Workflow times out

**Solution**: Target specific packages or benchmarks:

```yaml
with:
  packages: './internal/critical'
  benchmark-filter: 'BenchmarkImportant'
```

## Next Steps

- Read the full [Action Documentation](../../ACTION.md)
- Check out the [gokanon README](../../README.md)
- Explore [benchmark examples](../benchmarks/)

## Contributing

Have a useful workflow pattern? Please contribute it by opening a PR!
