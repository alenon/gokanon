# CI/CD Integration Guide

This guide shows you how to integrate gokanon into your CI/CD pipeline to automatically catch performance regressions.

## Overview

gokanon can be used in CI/CD pipelines to:
- Run benchmarks automatically on every PR or commit
- Compare performance against a baseline (e.g., main branch)
- Fail builds if performance degrades beyond a threshold
- Generate and archive benchmark reports

## GitHub Actions

A complete example is provided in `.github/workflows/benchmark.yml`. Here's how it works:

### Basic Setup

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
      run: gokanon run -pkg=./... -bench=.

    - name: Check threshold
      run: gokanon check --latest -threshold=10
```

### Advanced: Compare Against Baseline

To compare PR benchmarks against the main branch:

```yaml
- name: Checkout baseline (main branch)
  run: |
    git fetch origin main:main
    git checkout main
    gokanon run -pkg=./... -bench=.

- name: Checkout PR branch
  run: |
    git checkout ${{ github.event.pull_request.head.sha }}
    gokanon run -pkg=./... -bench=.

- name: Compare results
  run: |
    RUNS=($(gokanon list | tail -n +2 | awk '{print $1}'))
    gokanon compare ${RUNS[1]} ${RUNS[0]}
    gokanon check ${RUNS[1]} ${RUNS[0]} -threshold=10
```

### Generate Reports

```yaml
- name: Generate HTML report
  run: gokanon export --latest -format=html -output=report.html

- name: Upload report
  uses: actions/upload-artifact@v5
  with:
    name: benchmark-report
    path: report.html
```

### Comment on PR (Optional)

Using [peter-evans/create-or-update-comment](https://github.com/peter-evans/create-or-update-comment):

```yaml
- name: Generate comparison
  id: comparison
  run: |
    COMPARISON=$(gokanon compare --latest)
    echo "comparison<<EOF" >> $GITHUB_OUTPUT
    echo "$COMPARISON" >> $GITHUB_OUTPUT
    echo "EOF" >> $GITHUB_OUTPUT

- name: Comment PR
  uses: peter-evans/create-or-update-comment@v3
  with:
    issue-number: ${{ github.event.pull_request.number }}
    body: |
      ## Benchmark Results

      ```
      ${{ steps.comparison.outputs.comparison }}
      ```
```

## GitLab CI

Example `.gitlab-ci.yml`:

```yaml
stages:
  - test
  - benchmark

benchmark:
  stage: benchmark
  image: golang:1.21
  before_script:
    - go install github.com/alenon/gokanon@latest
  script:
    # Run benchmarks
    - gokanon run -pkg=./... -bench=.

    # Download baseline from artifacts
    - 'curl --location --output baseline.tar.gz "$CI_API_V4_URL/projects/$CI_PROJECT_ID/jobs/artifacts/main/download?job=benchmark" || true'
    - tar -xzf baseline.tar.gz || true

    # Compare if baseline exists
    - |
      if [ $(gokanon list | wc -l) -gt 1 ]; then
        gokanon compare --latest
        gokanon check --latest -threshold=10
      fi

    # Export report
    - gokanon export --latest -format=html -output=benchmark-report.html || true

  artifacts:
    paths:
      - .gokanon/
      - benchmark-report.html
    expire_in: 30 days

  only:
    - merge_requests
    - main
```

## Jenkins

Example `Jenkinsfile`:

```groovy
pipeline {
    agent any

    environment {
        GOPATH = "${WORKSPACE}/go"
        PATH = "${GOPATH}/bin:${PATH}"
    }

    stages {
        stage('Setup') {
            steps {
                sh 'go install github.com/alenon/gokanon@latest'
            }
        }

        stage('Benchmark') {
            steps {
                sh 'gokanon run -pkg=./... -bench=.'
            }
        }

        stage('Compare') {
            when {
                changeRequest()
            }
            steps {
                script {
                    // Copy baseline from main branch
                    sh '''
                        git checkout main
                        gokanon run -pkg=./... -bench=.
                        git checkout ${CHANGE_BRANCH}
                        gokanon run -pkg=./... -bench=.
                    '''

                    // Compare
                    def comparison = sh(
                        script: 'gokanon compare --latest',
                        returnStdout: true
                    ).trim()

                    echo "Benchmark Comparison:\n${comparison}"

                    // Check threshold
                    sh 'gokanon check --latest -threshold=10'
                }
            }
        }

        stage('Report') {
            steps {
                sh 'gokanon export --latest -format=html -output=benchmark-report.html'
                publishHTML([
                    reportDir: '.',
                    reportFiles: 'benchmark-report.html',
                    reportName: 'Benchmark Report'
                ])
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: '.gokanon/**', allowEmptyArchive: true
        }
    }
}
```

## CircleCI

Example `.circleci/config.yml`:

```yaml
version: 2.1

jobs:
  benchmark:
    docker:
      - image: cimg/go:1.21

    steps:
      - checkout

      - restore_cache:
          keys:
            - gokanon-{{ .Branch }}
            - gokanon-main

      - run:
          name: Install gokanon
          command: go install github.com/alenon/gokanon@latest

      - run:
          name: Run benchmarks
          command: gokanon run -pkg=./... -bench=.

      - run:
          name: Compare with baseline
          command: |
            if [ $(gokanon list | wc -l) -gt 1 ]; then
              gokanon compare --latest
              gokanon check --latest -threshold=10
            fi

      - run:
          name: Generate report
          command: gokanon export --latest -format=html -output=benchmark-report.html
          when: always

      - save_cache:
          key: gokanon-{{ .Branch }}-{{ .Revision }}
          paths:
            - .gokanon

      - store_artifacts:
          path: benchmark-report.html

workflows:
  version: 2
  benchmark-check:
    jobs:
      - benchmark
```

## Best Practices

### 1. Set Appropriate Thresholds

Choose thresholds based on your needs:
- **Strict (1-5%)**: Critical performance-sensitive code
- **Moderate (5-10%)**: General application code
- **Lenient (10-20%)**: Non-critical paths or variable benchmarks

```bash
# Strict
gokanon check --latest -threshold=5

# Moderate
gokanon check --latest -threshold=10

# Lenient
gokanon check --latest -threshold=20
```

### 2. Use Statistical Analysis

For more reliable results, run benchmarks multiple times:

```bash
# Run benchmarks 5 times
for i in {1..5}; do
  gokanon run -pkg=./... -bench=.
done

# Analyze stability
gokanon stats -last=5
```

### 3. Isolate CI Runners

For consistent results:
- Use dedicated CI runners for benchmarks
- Avoid running other workloads during benchmarks
- Use consistent hardware specifications

### 4. Store Historical Data

Keep benchmark history for trend analysis:

```bash
# Archive results
tar -czf benchmarks-$(date +%Y%m%d).tar.gz .gokanon/

# Upload to artifact storage (S3, etc.)
aws s3 cp benchmarks-*.tar.gz s3://my-bucket/benchmarks/
```

### 5. Fail Fast

Use the `check` command to fail builds early:

```bash
# This will exit with code 1 if threshold is exceeded
gokanon check --latest -threshold=10 || exit 1
```

## Environment Variables

You can use environment variables for configuration:

```bash
export GOKANON_STORAGE="/tmp/gokanon-results"
export GOKANON_THRESHOLD="10"

gokanon run -storage="$GOKANON_STORAGE"
gokanon check --latest -threshold="$GOKANON_THRESHOLD"
```

## Troubleshooting

### Benchmarks Too Variable

If benchmarks show high variation:

```bash
# Check coefficient of variation
gokanon stats -last=5

# Increase benchmark time
go test -bench=. -benchtime=10s
```

### No Baseline Found

Ensure your baseline is saved:

```bash
# On main branch
gokanon run -pkg=./...

# Save to artifact storage
tar -czf baseline.tar.gz .gokanon/
```

### CI Runner Performance

Different CI runners may show different results:
- Use the same runner type for baseline and comparison
- Consider using self-hosted runners for consistency
- Document the runner specs in your results

## Example: Complete PR Workflow

Here's a complete example that:
1. Runs benchmarks on both base and PR branches
2. Compares results
3. Generates HTML report
4. Comments on PR
5. Fails if degradation > 10%

```yaml
name: PR Benchmark Check

on:
  pull_request:
    branches: [ main ]

jobs:
  benchmark:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v5
      with:
        fetch-depth: 0

    - uses: actions/setup-go@v6
      with:
        go-version: '1.21'

    - name: Install gokanon
      run: go install github.com/alenon/gokanon@latest

    - name: Run baseline benchmarks
      run: |
        git checkout ${{ github.event.pull_request.base.sha }}
        gokanon run -pkg=./... -bench=.

    - name: Run PR benchmarks
      run: |
        git checkout ${{ github.event.pull_request.head.sha }}
        gokanon run -pkg=./... -bench=.

    - name: Compare and check
      id: compare
      run: |
        # Get comparison output
        COMPARISON=$(gokanon compare --latest)
        echo "$COMPARISON"

        # Save for PR comment
        echo "comparison<<EOF" >> $GITHUB_OUTPUT
        echo "$COMPARISON" >> $GITHUB_OUTPUT
        echo "EOF" >> $GITHUB_OUTPUT

        # Check threshold
        gokanon check --latest -threshold=10

    - name: Generate HTML report
      if: always()
      run: gokanon export --latest -format=html -output=report.html

    - name: Upload report
      if: always()
      uses: actions/upload-artifact@v5
      with:
        name: benchmark-report
        path: report.html

    - name: Comment PR
      if: always()
      uses: peter-evans/create-or-update-comment@v3
      with:
        issue-number: ${{ github.event.pull_request.number }}
        body: |
          ## 📊 Benchmark Results

          ```
          ${{ steps.compare.outputs.comparison }}
          ```

          [View detailed HTML report](https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }})
```

## Next Steps

- Review the example workflow in `.github/workflows/benchmark.yml`
- Customize thresholds for your project
- Set up artifact storage for historical data
- Configure notifications for performance regressions
