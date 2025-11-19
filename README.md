<div align="center">

# 🚀 gokanon

### A powerful CLI tool for running and comparing Go benchmark tests

*Profiling • Analysis • CI/CD Integration*

[![License](https://img.shields.io/github/license/alenon/gokanon)](LICENSE)
[![Release](https://img.shields.io/github/v/release/alenon/gokanon)](https://github.com/alenon/gokanon/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/alenon/gokanon)](https://goreportcard.com/report/github.com/alenon/gokanon)
[![Go Version](https://img.shields.io/github/go-mod/go-version/alenon/gokanon)](go.mod)

[Installation](#-installation) •
[Quick Start](#-quick-start) •
[Features](#-features) •
[Documentation](#-usage) •
[CI/CD](#-cicd-integration)

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎯 Core Capabilities
- ⚡ **Run & Save Benchmarks** - Execute and automatically save results
- 🔥 **CPU/Memory Profiling** - Flame graph visualization
- 📊 **Compare Results** - Detailed performance analysis
- 📈 **Statistical Analysis** - Mean, median, stability metrics

</td>
<td width="50%">

### 🚀 Advanced Features
- 📉 **Trend Analysis** - Regression detection over time
- 🎨 **Interactive Dashboard** - Web-based visualization
- 📝 **Export Reports** - HTML, CSV, Markdown formats
- 🤖 **AI Analysis** - Intelligent optimization suggestions

</td>
</tr>
<tr>
<td colspan="2">

### 🔧 Integration & Automation
✅ **CI/CD Ready** - GitHub Actions support with automated threshold checking
✅ **Multiple AI Providers** - Ollama, OpenAI, Claude, Gemini, Groq
✅ **Shell Completion** - Bash, Zsh, Fish support
✅ **Baseline Management** - Track and compare against reference points

</td>
</tr>
</table>

## 📦 Installation

### ⚡ Quick Install

<table>
<tr>
<td>

**🍎 macOS & Linux**
```bash
curl -sSL https://raw.githubusercontent.com/alenon/gokanon/main/install.sh | bash
```

</td>
<td>

**🍺 Homebrew**
```bash
brew install alenon/tap/gokanon
```

</td>
</tr>
</table>

### 📥 Pre-built Binaries

> 📌 Download from [GitHub Releases](https://github.com/alenon/gokanon/releases/latest)

<details>
<summary><b>🍎 macOS (Apple Silicon)</b></summary>

```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-darwin-arm64.tar.gz | tar xz
sudo mv gokanon-darwin-arm64 /usr/local/bin/gokanon
```
</details>

<details>
<summary><b>🐧 Linux (x86_64)</b></summary>

```bash
curl -L https://github.com/alenon/gokanon/releases/latest/download/gokanon-linux-amd64.tar.gz | tar xz
sudo mv gokanon-linux-amd64 /usr/local/bin/gokanon
```
</details>

### 🔨 Install with Go

```bash
go install github.com/alenon/gokanon@latest
```

---

## 🚀 Quick Start

```bash
# 1️⃣ Run benchmarks
gokanon run -pkg=./...

# 2️⃣ View interactive dashboard
gokanon serve

# 3️⃣ Run with profiling
gokanon run --profile=cpu,mem

# 4️⃣ View flame graphs
gokanon flamegraph --latest

# 5️⃣ Compare results
gokanon compare --latest

# 6️⃣ Export to HTML
gokanon export --latest -format=html

# 7️⃣ Check performance threshold (CI/CD)
gokanon check --latest -threshold=10
```

> 💡 **Tip:** Run `gokanon help` to see all available commands and options

---

## 📖 Usage

### ⚡ Running Benchmarks

```bash
# All benchmarks in current package
gokanon run

# Specific benchmark pattern
gokanon run -bench=BenchmarkStringBuilder

# All packages
gokanon run -pkg=./...

# With profiling
gokanon run --profile=cpu,mem

# Control CPU parallelism and benchmark duration
gokanon run -cpu=1,2,4 -benchtime=1s
```

### 🔥 Profiling & Analysis

Generate CPU and memory profiles to identify bottlenecks:

```bash
# Enable profiling
gokanon run --profile=cpu,mem

# View flame graphs
gokanon flamegraph --latest
```

**The profiler automatically:**
- 🎯 Identifies hot functions and memory allocation patterns
- 🔍 Detects potential memory leaks
- 💡 Provides optimization suggestions with impact analysis

### 🤖 AI-Powered Analysis

Enable AI analysis for intelligent insights:

<details>
<summary><b>🟢 Using Ollama (Free, Local)</b></summary>

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=ollama
```
</details>

<details>
<summary><b>🔵 Using OpenAI</b></summary>

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=openai
export GOKANON_AI_API_KEY=sk-your-key
export GOKANON_AI_MODEL=gpt-4o
```
</details>

<details>
<summary><b>🟣 Using Anthropic Claude</b></summary>

```bash
export GOKANON_AI_ENABLED=true
export GOKANON_AI_PROVIDER=anthropic
export GOKANON_AI_API_KEY=sk-ant-your-key
```
</details>

```bash
# Run with AI analysis
gokanon run --profile=cpu,mem
gokanon compare --latest
```

> 🔌 **Supported Providers:** Ollama, OpenAI, Anthropic, Gemini, Groq, OpenAI-compatible APIs

### 🎨 Interactive Dashboard

```bash
# Start dashboard
gokanon serve

# Custom port
gokanon serve -port=9000
```

Access at `http://localhost:8080` for:
- 📈 Real-time performance trends
- 📊 Historical data with charts
- ⚖️ Side-by-side comparisons
- 🌙 Dark mode support

### 📊 Comparing Results

```bash
# Last two runs
gokanon compare --latest

# Specific runs
gokanon compare run-123 run-456

# Against baseline
gokanon compare --baseline=v1.0
```

### 📈 Statistical & Trend Analysis

```bash
# Analyze last 5 runs
gokanon stats -last=5

# Track performance trends
gokanon trend -last=10
```

### 📝 Exporting Reports

```bash
# HTML report
gokanon export --latest -format=html -output=report.html

# CSV for analysis
gokanon export --latest -format=csv -output=results.csv

# Markdown for docs
gokanon export --latest -format=markdown -output=comparison.md
```

### 🔄 CI/CD Integration

```bash
# Fail if degradation > 10%
gokanon check --latest -threshold=10
```

**GitHub Action Example:**
```yaml
- name: Run benchmarks
  uses: alenon/gokanon@v1
  with:
    packages: './...'
    threshold-percent: 10
    enable-profiling: 'cpu,mem'
    cpu: '1,2,4'
    benchtime: '1s'
    export-format: 'html'
```

> 📋 See `action.yml` for complete GitHub Action configuration

### 🗂️ Managing Results

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

## 🔧 Commands Reference

<table>
<tr>
<td width="50%">

**Core Commands**
```bash
gokanon run         # Run & save benchmarks
gokanon list        # List saved results
gokanon compare     # Compare results
gokanon export      # Export to HTML/CSV/MD
gokanon stats       # Statistical analysis
gokanon trend       # Trend analysis
gokanon check       # Threshold checking
gokanon flamegraph  # View flame graphs
```

</td>
<td width="50%">

**Utility Commands**
```bash
gokanon serve        # Interactive dashboard
gokanon delete       # Delete results
gokanon baseline     # Manage baselines
gokanon doctor       # Run diagnostics
gokanon interactive  # Interactive mode
gokanon completion   # Shell completion
gokanon version      # Version info
gokanon help         # Show help
```

</td>
</tr>
</table>

## 💾 Storage

Results are stored in `.gokanon` directory by default. Use `-storage` flag to customize:

```bash
gokanon run -storage=./benchmark-results
gokanon list -storage=./benchmark-results
```

---

## 💡 Best Practices

| Practice | Description |
|----------|-------------|
| 🖥️ **Consistent Environment** | Run benchmarks on consistent hardware and system load |
| 🔄 **Multiple Runs** | Run benchmarks multiple times for reliable results |
| 📏 **Baseline** | Maintain a baseline for comparing optimizations |
| 🔗 **CI Integration** | Catch regressions early with automated checks |
| 🎯 **Appropriate Thresholds** | Set thresholds based on application requirements |

---

## 🛠️ Development

```bash
# Build the project
make build

# Run tests
make test

# Generate coverage report
make coverage

# Build for all platforms
make build-all
```

---

## 📄 License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## 🤝 Contributing

Contributions are welcome! We appreciate your help in making gokanon better.

- 🐛 **Found a bug?** [Open an issue](https://github.com/alenon/gokanon/issues)
- 💡 **Have an idea?** [Start a discussion](https://github.com/alenon/gokanon/discussions)
- 🔧 **Want to contribute?** [Submit a pull request](https://github.com/alenon/gokanon/pulls)

---

<div align="center">

**Made with ❤️ for the Go community**

[![GitHub](https://img.shields.io/badge/GitHub-alenon%2Fgokanon-blue?logo=github)](https://github.com/alenon/gokanon)
[![Issues](https://img.shields.io/github/issues/alenon/gokanon)](https://github.com/alenon/gokanon/issues)
[![Stars](https://img.shields.io/github/stars/alenon/gokanon)](https://github.com/alenon/gokanon/stargazers)

</div>
