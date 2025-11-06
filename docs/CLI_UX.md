# CLI UX Enhancements

This document describes the user experience improvements added to gokanon.

## Features

### 1. Enhanced Colored Output

GoKanon now features rich colored terminal output for better readability:

- **Success messages** in green with ✓ checkmarks
- **Error messages** in red with ✗ marks
- **Warning messages** in yellow with ⚠ symbols
- **Info messages** in cyan with ℹ icons
- **Dimmed text** for less important information
- **Bold text** for emphasis

The color output automatically disables when:
- Running in a non-TTY environment (e.g., CI/CD pipelines)
- The `NO_COLOR` environment variable is set
- Output is redirected to a file

**Example:**
```bash
gokanon run -bench=.
```

Output will include colored status indicators, formatted performance metrics, and clear section headers.

### 2. Progress Indicators

Long-running operations now show visual feedback:

#### Spinners
Used for indeterminate operations like running benchmarks:
- Animated spinner characters (⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏)
- Clear operation descriptions
- Auto-clears when operation completes

#### Progress Bars
For operations with known progress:
- Visual bar showing completion percentage
- Customizable descriptions
- Color-coded progress indicators

### 3. Interactive Mode

Start an interactive gokanon session with tab completion and command history:

```bash
gokanon interactive
# or the short form:
gokanon i
```

**Features:**
- **Tab completion** for commands and flags
- **Command history** (stored in `/tmp/gokanon_history.tmp`)
- **Built-in commands:**
  - `help` - Show available commands
  - `clear` - Clear the screen
  - `exit` or `quit` - Exit interactive mode
  - All regular gokanon commands

**Example session:**
```
╔════════════════════════════════════════════════════════════╗
║          Welcome to gokanon Interactive Mode!          ║
╚════════════════════════════════════════════════════════════╝

  Type 'help' for available commands
  Use TAB for auto-completion
  Press Ctrl+C or type 'exit' to quit

gokanon> run -bench=.
gokanon> list
gokanon> compare --latest
gokanon> exit
```

### 4. Shell Auto-completion

Install auto-completion for your shell to get command and flag suggestions:

#### Bash
```bash
# System-wide installation
gokanon completion bash | sudo tee /etc/bash_completion.d/gokanon

# User-only installation
mkdir -p ~/.local/share/bash-completion/completions
gokanon completion bash > ~/.local/share/bash-completion/completions/gokanon
```

#### Zsh
```bash
# System-wide installation
gokanon completion zsh | sudo tee /usr/local/share/zsh/site-functions/_gokanon

# Add to .zshrc for auto-loading
echo 'source <(gokanon completion zsh)' >> ~/.zshrc
```

#### Fish
```bash
# User installation
gokanon completion fish > ~/.config/fish/completions/gokanon.fish
```

**Supported completions:**
- Main commands (run, list, compare, etc.)
- Command-specific flags
- Export format options (html, csv, markdown, json)
- Profile types (cpu, mem)

### 5. Improved Error Messages

Errors now include helpful suggestions for resolution:

**Example:**
```bash
$ gokanon run -profile=invalid

✗ Unknown profile type: invalid
  Cause: profile type not recognized

💡 Suggestions:
  → Valid profile types: cpu, mem
  → Example: -profile=cpu,mem
```

**Common error scenarios with suggestions:**
- No benchmark results found
- Invalid run IDs
- Benchmark execution failures
- Invalid threshold values
- Storage corruption
- Profile data not found
- Unsupported export formats
- Port already in use

### 6. Doctor Command

Diagnose issues with your gokanon setup:

```bash
gokanon doctor
```

**Checks performed:**
1. **Go Installation** - Verifies Go is installed and accessible
2. **Go Test Command** - Ensures `go test` can be executed
3. **Storage Directory** - Checks .gokanon directory status
4. **Storage Integrity** - Validates saved benchmark data
5. **Benchmark Files** - Looks for test files with benchmarks
6. **Git Repository** - Detects if you're in a git repo (optional)
7. **System Resources** - Checks available memory

**Example output:**
```
Running gokanon diagnostics...

📊 Diagnostic Results

✓ Go Installation: go version go1.24.7 linux/amd64
✓ Go Test Command: 'go test' command is available
✓ Storage Directory: Storage directory exists at: .gokanon
✓ Storage Integrity: Storage is healthy with 5 run(s)
✓ Benchmark Files: Found 3 test file(s) with benchmarks
✓ Git Repository: Git repository detected
✓ System Resources: System memory: 16384 MB, Go runtime: go1.24.7

Summary
✅ 7 checks passed

All checks passed! Your gokanon setup is healthy.
```

## UI Components

### Color Functions

The `internal/ui/colors.go` package provides:

- `Success()` - Green bold text
- `Error()` - Red bold text
- `Warning()` - Yellow bold text
- `Info()` - Cyan text
- `Dim()` - Faint text
- `Bold()` - Bold text

### Print Functions

- `PrintSuccess(format, args...)` - Print success message with ✓
- `PrintError(format, args...)` - Print error message with ✗
- `PrintWarning(format, args...)` - Print warning message with ⚠
- `PrintInfo(format, args...)` - Print info message with ℹ
- `PrintHeader(text)` - Print bold header with underline
- `PrintSection(emoji, title)` - Print section header with emoji

### Progress Components

- `NewProgressBar(max, description)` - Create a progress bar
- `NewIndeterminateSpinner(description)` - Create a spinner
- `NewSpinner(message)` - Create a simple spinner

### Error Handling

- `NewError(message, err, suggestions...)` - Create error with suggestions
- `ErrNoResults()` - No benchmark results found error
- `ErrInvalidRunID(id, available)` - Invalid run ID error
- `ErrBenchmarkFailed(err)` - Benchmark execution error
- `ErrStorageCorrupted(err)` - Storage corruption error
- `ErrProfileNotFound(runID)` - Profile data missing error

## Environment Variables

- `NO_COLOR` - Set to any value to disable colored output

## Best Practices

1. **Use `gokanon doctor`** when setting up a new project or troubleshooting
2. **Install shell completion** for faster command entry
3. **Use interactive mode** when experimenting with different commands
4. **Check error suggestions** when commands fail - they often contain the solution
5. **Disable colors in CI/CD** by setting `NO_COLOR=1` if needed (though gokanon auto-detects)

## Migration Notes

The enhanced UX features are fully backward compatible. Existing scripts and CI/CD pipelines will continue to work without modification. The colored output automatically disables in non-interactive environments.

## Technical Details

### Dependencies

- `github.com/fatih/color` - Terminal color support
- `github.com/schollz/progressbar/v3` - Progress bars and spinners
- `github.com/chzyer/readline` - Interactive mode with completion

### Architecture

```
internal/
├── ui/
│   ├── colors.go      # Color and formatting utilities
│   ├── progress.go    # Progress bars and spinners
│   └── errors.go      # Enhanced error messages
├── doctor/
│   └── doctor.go      # Diagnostic checks
└── interactive/
    └── interactive.go # Interactive mode implementation
```

### Performance

- Colors are rendered using ANSI escape codes (zero allocation)
- Progress indicators run in separate goroutines
- Spinners update at 100ms intervals
- Interactive mode stores history efficiently

## Examples

### Running benchmarks with progress
```bash
gokanon run -bench=. -profile=cpu,mem
```

### Using interactive mode
```bash
gokanon interactive
gokanon> run -bench=BenchmarkSort
gokanon> list
gokanon> compare --latest
```

### Diagnosing issues
```bash
gokanon doctor
```

### Installing completions
```bash
# Bash
gokanon completion bash > ~/.local/share/bash-completion/completions/gokanon

# Zsh
echo 'source <(gokanon completion zsh)' >> ~/.zshrc

# Fish
gokanon completion fish > ~/.config/fish/completions/gokanon.fish
```

## Future Enhancements

Potential future improvements:
- Real-time benchmark progress updates
- More interactive visualizations
- Configurable color schemes
- Extended diagnostic checks
- Integration with popular CI/CD platforms
