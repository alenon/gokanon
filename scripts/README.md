# GoKanon Scripts

This directory contains utility scripts for testing, reporting, and development workflows.

## 📊 Test Report Script

### `test-report.sh`

Generates a beautiful, easy-to-understand test report with coverage metrics.

**Features:**
- 🎨 **Color-coded output** - Visual indicators for pass/fail and coverage levels
- 📊 **Progress bars** - Visual representation of coverage percentages
- 🎯 **Coverage grading** - A-F grading system for overall coverage
- 📦 **Package breakdown** - Detailed coverage for each package
- 🔍 **Command analysis** - Specific coverage for all GoKanon commands
- 💡 **Recommendations** - Actionable suggestions for improvement
- 📄 **File export** - Save reports for documentation or CI/CD

**Usage:**

```bash
# Run test report (console output)
./scripts/test-report.sh

# Save report to file
./scripts/test-report.sh test-report.txt

# View in CI/CD
./scripts/test-report.sh | tee test-report.log
```

**Sample Output:**

```
╔══════════════════════════════════════════════════════════════════╗
║                   GoKanon Test Report                            ║
╔══════════════════════════════════════════════════════════════════╗

Generated: 2025-11-07 08:46:54

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  TEST SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Status:        ✓ PASS
  Total Tests:   16
  Passed:        16

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  COVERAGE SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Overall Coverage: [███████████████████████████░░░░░░░░░░░░░]  69%

  Coverage Grade:   C - Fair

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PACKAGE COVERAGE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✓   internal/compare                           100.0%
  ✓   internal/stats                             100.0%
  ✓   internal/threshold                         100.0%
  ✓   internal/export                             91.8%
  ✓   internal/doctor                             91.5%
  ...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  COMMAND COVERAGE (Key Functions)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✓   Delete                                   100.0%
  ✓   Completion                                97.4%
  ✓   Stats                                     96.4%
  ✓   List                                      94.1%
  ...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RECOMMENDATIONS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ◐ Increase test coverage to at least 70%
  ✓ All tests passing!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Tests passed with good coverage!
```

**Coverage Grading System:**

| Coverage | Grade | Color |
|----------|-------|-------|
| ≥ 80%    | A - Excellent | Green |
| ≥ 70%    | B - Good | Cyan |
| ≥ 60%    | C - Fair | Yellow |
| ≥ 50%    | D - Needs Improvement | Yellow |
| < 50%    | F - Poor | Red |

**Package/Function Indicators:**

- ✓ **Green** - Good coverage (≥ 80%)
- ◐ **Yellow** - Fair coverage (60-79%)
- ✗ **Red** - Low coverage (< 60%)
- ○ **Cyan** - No tests

**Exit Codes:**

- `0` - All tests passed with ≥60% coverage
- `1` - Tests failed or coverage <60%

**Integration with CI/CD:**

See `.github/workflows/test-report.yml` for GitHub Actions integration.

The workflow automatically:
- Runs tests on push and PRs
- Generates beautiful test reports
- Posts reports as PR comments
- Uploads reports as artifacts
- Fails if tests fail

**Requirements:**

- Go 1.21 or higher
- bash shell
- `go test` command
- `bc` utility (for calculations)

**Tips:**

1. **Run before committing** to ensure tests pass
2. **Check recommendations** for areas needing improvement
3. **Save reports** for historical tracking
4. **Use in pre-commit hooks** to enforce quality gates
5. **Include in CI/CD pipelines** for automated testing

## Development

When adding new scripts:
1. Make them executable: `chmod +x scripts/your-script.sh`
2. Add documentation to this README
3. Include usage examples
4. Add error handling
5. Test in clean environment
