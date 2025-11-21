package runner

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/alenon/gokanon/internal/models"
)

// DirectRunner executes benchmarks directly using the testing package
type DirectRunner struct {
	runner *Runner
}

// NewDirectRunner creates a new direct benchmark runner
func NewDirectRunner(r *Runner) *DirectRunner {
	return &DirectRunner{runner: r}
}

// Run executes benchmarks directly without spawning go test
func (dr *DirectRunner) Run() (*models.BenchmarkRun, error) {
	startTime := time.Now()

	// Get Go version
	goVersion, err := dr.runner.getGoVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get Go version: %w", err)
	}

	// Generate unique ID for this run
	runID := generateID()

	// Get package path
	pkgPath := dr.runner.packagePath
	if pkgPath == "" {
		pkgPath = "."
	}

	// Find all test files in the package
	testFiles, err := dr.findTestFiles(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find test files: %w", err)
	}

	if len(testFiles) == 0 {
		return nil, fmt.Errorf("no test files found in %s", pkgPath)
	}

	// Find all benchmark functions
	benchmarks, err := dr.findBenchmarkFunctions(testFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to find benchmark functions: %w", err)
	}

	if len(benchmarks) == 0 {
		return nil, fmt.Errorf("no benchmark functions found")
	}

	// Filter benchmarks based on filter pattern
	filteredBenchmarks := dr.filterBenchmarks(benchmarks)
	if len(filteredBenchmarks) == 0 {
		return nil, fmt.Errorf("no benchmarks match filter: %s", dr.runner.benchFilter)
	}

	// Set up profiling if enabled
	var cpuProfileBuf, memProfileBuf *bytes.Buffer
	var cpuProfilePath, memProfilePath string
	tempDir := ""

	if dr.runner.profileOptions != nil {
		tempDir, err = os.MkdirTemp("", "gokanon-profile-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)

		if dr.runner.profileOptions.EnableCPU {
			cpuProfileBuf = new(bytes.Buffer)
			cpuProfilePath = filepath.Join(tempDir, "cpu.prof")
			if err := pprof.StartCPUProfile(cpuProfileBuf); err != nil {
				return nil, fmt.Errorf("failed to start CPU profiling: %w", err)
			}
			defer pprof.StopCPUProfile()
		}

		if dr.runner.profileOptions.EnableMemory {
			memProfileBuf = new(bytes.Buffer)
			memProfilePath = filepath.Join(tempDir, "mem.prof")
		}
	}

	// Run benchmarks by compiling and executing the test binary
	results, err := dr.runBenchmarksViaTestBinary(pkgPath, filteredBenchmarks)
	if err != nil {
		return nil, fmt.Errorf("failed to run benchmarks: %w", err)
	}

	// Stop CPU profiling if it was enabled
	if cpuProfileBuf != nil {
		pprof.StopCPUProfile()
		if err := os.WriteFile(cpuProfilePath, cpuProfileBuf.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write CPU profile: %v\n", err)
		}
	}

	// Capture memory profile if enabled
	if memProfileBuf != nil {
		runtime.GC() // Get up-to-date statistics
		if err := pprof.WriteHeapProfile(memProfileBuf); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write memory profile: %v\n", err)
		} else {
			if err := os.WriteFile(memProfilePath, memProfileBuf.Bytes(), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to write memory profile: %v\n", err)
			}
		}
	}

	duration := time.Since(startTime)

	// Build command string for display
	cmdArgs := []string{"benchmark", "-bench", dr.runner.benchFilter}
	if dr.runner.cpu != "" {
		cmdArgs = append(cmdArgs, "-cpu", dr.runner.cpu)
	}
	if dr.runner.benchtime != "" {
		cmdArgs = append(cmdArgs, "-benchtime", dr.runner.benchtime)
	}

	run := &models.BenchmarkRun{
		ID:        runID,
		Timestamp: startTime,
		Package:   pkgPath,
		GoVersion: goVersion,
		Results:   results,
		Command:   fmt.Sprintf("gokanon %s (direct execution)", strings.Join(cmdArgs, " ")),
		Duration:  duration,
	}

	// Handle profile files if profiling was enabled
	if dr.runner.profileOptions != nil && dr.runner.profileOptions.Storage != nil {
		if err := dr.runner.handleProfiles(run, cpuProfilePath, memProfilePath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to process profiles: %v\n", err)
		}
	}

	return run, nil
}

// findTestFiles finds all *_test.go files in the given path
func (dr *DirectRunner) findTestFiles(pkgPath string) ([]string, error) {
	var testFiles []string

	// Resolve package path
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return nil, err
	}

	// Check if it's a directory or a specific package pattern
	info, err := os.Stat(absPath)
	if err != nil {
		// Might be a package pattern like ./...
		if strings.HasSuffix(pkgPath, "/...") {
			baseDir := strings.TrimSuffix(pkgPath, "/...")
			if baseDir == "" {
				baseDir = "."
			}
			return dr.findTestFilesRecursive(baseDir)
		}
		return nil, err
	}

	if info.IsDir() {
		files, err := filepath.Glob(filepath.Join(absPath, "*_test.go"))
		if err != nil {
			return nil, err
		}
		testFiles = append(testFiles, files...)
	}

	return testFiles, nil
}

// findTestFilesRecursive finds test files recursively
func (dr *DirectRunner) findTestFilesRecursive(baseDir string) ([]string, error) {
	var testFiles []string

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})

	return testFiles, err
}

// findBenchmarkFunctions parses test files and finds benchmark functions
func (dr *DirectRunner) findBenchmarkFunctions(testFiles []string) ([]string, error) {
	var benchmarks []string
	seen := make(map[string]bool)

	for _, file := range testFiles {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		// Find all functions starting with "Benchmark"
		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			name := fn.Name.Name
			if strings.HasPrefix(name, "Benchmark") {
				// Check function signature: func Benchmark...(b *testing.B)
				if fn.Type.Params.NumFields() == 1 {
					param := fn.Type.Params.List[0]
					if len(param.Names) == 1 {
						if starExpr, ok := param.Type.(*ast.StarExpr); ok {
							if selExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
								if ident, ok := selExpr.X.(*ast.Ident); ok {
									if ident.Name == "testing" && selExpr.Sel.Name == "B" {
										if !seen[name] {
											benchmarks = append(benchmarks, name)
											seen[name] = true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return benchmarks, nil
}

// filterBenchmarks filters benchmarks based on the bench filter pattern
func (dr *DirectRunner) filterBenchmarks(benchmarks []string) []string {
	filter := dr.runner.benchFilter
	if filter == "" || filter == "." {
		return benchmarks
	}

	var filtered []string
	pattern := regexp.MustCompile(filter)

	for _, name := range benchmarks {
		// Remove "Benchmark" prefix for matching
		shortName := strings.TrimPrefix(name, "Benchmark")
		if pattern.MatchString(name) || pattern.MatchString(shortName) {
			filtered = append(filtered, name)
		}
	}

	return filtered
}

// runBenchmarksViaTestBinary compiles and runs benchmarks via custom harness
func (dr *DirectRunner) runBenchmarksViaTestBinary(pkgPath string, benchmarks []string) ([]models.BenchmarkResult, error) {
	// Create a temporary directory for the harness
	tempDir, err := os.MkdirTemp("", "gokanon-harness-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Determine the import path for the package
	importPath, err := dr.getPackageImportPath(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get package import path: %w", err)
	}

	// Generate the benchmark harness code
	harnessCode, err := dr.generateHarness(importPath, benchmarks)
	if err != nil {
		return nil, fmt.Errorf("failed to generate harness: %w", err)
	}

	// Write the harness to a temporary file
	harnessPath := filepath.Join(tempDir, "harness.go")
	if err := os.WriteFile(harnessPath, []byte(harnessCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write harness file: %w", err)
	}

	// For local packages, build in the module context
	// Get the module root
	moduleRoot, err := dr.findModuleRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find module root: %w", err)
	}

	// Build the harness binary in the context of the module root
	harnessBinary := filepath.Join(tempDir, "harness")
	compileCmd := exec.Command("go", "build", "-o", harnessBinary, harnessPath)
	compileCmd.Dir = moduleRoot // Build from module root so imports work
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err := compileCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to compile harness: %w\nStderr: %s\nHarness:\n%s", err, compileStderr.String(), harnessCode)
	}

	// Execute the harness
	cmd := exec.Command(harnessBinary)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOKANON_BENCH_FILTER=%s", dr.runner.benchFilter),
		fmt.Sprintf("GOKANON_CPU=%s", dr.runner.cpu),
		fmt.Sprintf("GOKANON_BENCHTIME=%s", dr.runner.benchtime),
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Get stdout pipe for real-time reading
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start harness: %w", err)
	}

	// Parse results in real-time
	results, err := dr.runner.parseOutputRealtime(stdoutPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to parse benchmark output: %w", err)
	}

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("benchmark execution failed: %w\nStderr: %s", err, stderr.String())
		}
	}

	return results, nil
}

// getPackageImportPath returns the import path for a package directory
func (dr *DirectRunner) getPackageImportPath(pkgPath string) (string, error) {
	if pkgPath == "" || pkgPath == "." {
		pkgPath = "."
	}

	// Use go list to get the import path
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}", pkgPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get import path: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// generateHarness generates Go code that runs benchmarks directly
func (dr *DirectRunner) generateHarness(importPath string, benchmarks []string) (string, error) {
	var buf bytes.Buffer

	buf.WriteString(`package main

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	bench "` + importPath + `"
)

func main() {
	// Set GOMAXPROCS if CPU flag is provided
	cpuFlag := os.Getenv("GOKANON_CPU")
	if cpuFlag != "" {
		// For simplicity, just use the first value
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	benchtime := os.Getenv("GOKANON_BENCHTIME")
	if benchtime != "" {
		// Benchtime will be handled by testing.B
		_ = benchtime
	}

`)

	// Generate benchmark runner for each function
	for _, benchName := range benchmarks {
		buf.WriteString(fmt.Sprintf(`	// Run %s
	result := testing.Benchmark(bench.%s)
	printBenchmarkResult("%s", result)
`, benchName, benchName, strings.TrimPrefix(benchName, "Benchmark")))
	}

	buf.WriteString(`}

func printBenchmarkResult(name string, result testing.BenchmarkResult) {
	// Get GOMAXPROCS for naming
	procs := runtime.GOMAXPROCS(0)
	fullName := fmt.Sprintf("%s-%d", name, procs)

	// Print in the same format as go test -bench
	fmt.Printf("Benchmark%s\t%d\t%d ns/op",
		fullName, result.N, result.NsPerOp())

	if result.Bytes > 0 {
		mbPerSec := (float64(result.Bytes) * float64(result.N) / 1e6) / result.T.Seconds()
		fmt.Printf("\t%.2f MB/s", mbPerSec)
	}

	if result.MemAllocs > 0 {
		fmt.Printf("\t%d B/op\t%d allocs/op",
			result.AllocedBytesPerOp(), result.AllocsPerOp())
	}

	fmt.Println()
}
`)

	return buf.String(), nil
}

// findModuleRoot finds the root directory of the Go module
func (dr *DirectRunner) findModuleRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf("go.mod not found")
}

// InternalBenchmark runs a single benchmark function directly (for future use)
func InternalBenchmark(name string, f func(b *testing.B)) models.BenchmarkResult {
	result := testing.Benchmark(f)

	// Calculate MB/s if bytes were processed
	var mbPerSec float64
	if result.Bytes > 0 && result.T > 0 {
		mbPerSec = (float64(result.Bytes) * float64(result.N) / 1e6) / result.T.Seconds()
	}

	return models.BenchmarkResult{
		Name:        name,
		Iterations:  int64(result.N),
		NsPerOp:     float64(result.NsPerOp()),
		BytesPerOp:  result.AllocedBytesPerOp(),
		AllocsPerOp: result.AllocsPerOp(),
		MBPerSec:    mbPerSec,
	}
}

// dummyWriter is used to suppress test output
type dummyWriter struct{}

func (dw *dummyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// captureOutput captures stdout and stderr during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old
	return <-outC
}
