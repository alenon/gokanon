package webserver

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alenon/gokanon/internal/models"
	"github.com/alenon/gokanon/internal/storage"
	"github.com/google/pprof/profile"
)

// Server handles web serving of profile visualizations
type Server struct {
	storage *storage.Storage
	port    string
}

// NewServer creates a new web server
func NewServer(store *storage.Storage, port string) *Server {
	return &Server{
		storage: store,
		port:    port,
	}
}

// Start starts the web server
func (s *Server) Start(runID string) error {
	// Load the benchmark run
	run, err := s.storage.Load(runID)
	if err != nil {
		return fmt.Errorf("failed to load run: %w", err)
	}

	if run.CPUProfile == "" && run.MemoryProfile == "" {
		return fmt.Errorf("no profiles found for run %s", runID)
	}

	// Setup HTTP handlers
	mux := http.NewServeMux()

	// Main page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.handleIndex(w, r, run)
	})

	// CPU profile visualization
	if run.CPUProfile != "" {
		mux.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
			s.handleProfile(w, r, run.CPUProfile, "CPU Profile")
		})
		mux.HandleFunc("/cpu/flamegraph", func(w http.ResponseWriter, r *http.Request) {
			s.handleFlameGraph(w, r, run.CPUProfile, "CPU")
		})
	}

	// Memory profile visualization
	if run.MemoryProfile != "" {
		mux.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
			s.handleProfile(w, r, run.MemoryProfile, "Memory Profile")
		})
		mux.HandleFunc("/mem/flamegraph", func(w http.ResponseWriter, r *http.Request) {
			s.handleFlameGraph(w, r, run.MemoryProfile, "Memory")
		})
	}

	// Profile comparison
	if run.CPUProfile != "" && run.MemoryProfile != "" {
		mux.HandleFunc("/compare", func(w http.ResponseWriter, r *http.Request) {
			s.handleCompare(w, r, run)
		})
	}

	// Static assets (if needed)
	mux.HandleFunc("/static/", s.handleStatic)

	addr := ":" + s.port
	fmt.Printf("Starting profile visualization server at http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	return http.ListenAndServe(addr, mux)
}

// handleIndex shows the main page with links to different views
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request, run *models.BenchmarkRun) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"float64": func(i int64) float64 {
			return float64(i)
		},
	}

	tmpl := template.Must(template.New("index").Funcs(funcMap).Parse(indexTemplate))
	data := struct {
		Run        *models.BenchmarkRun
		HasCPU     bool
		HasMemory  bool
		HasSummary bool
	}{
		Run:        run,
		HasCPU:     run.CPUProfile != "",
		HasMemory:  run.MemoryProfile != "",
		HasSummary: run.ProfileSummary != nil,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// handleProfile serves pprof profile visualization
func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request, profilePath, title string) {
	// Read profile file
	data, err := os.ReadFile(profilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read profile: %v", err), http.StatusInternalServerError)
		return
	}

	// Serve using pprof handler
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.prof", strings.ToLower(title)))
	w.Write(data)
}

// handleFlameGraph generates and serves a flame graph
func (s *Server) handleFlameGraph(w http.ResponseWriter, r *http.Request, profilePath, profileType string) {
	// Check if profile exists
	if _, err := os.Stat(profilePath); err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %v", err), http.StatusNotFound)
		return
	}

	// Try to generate SVG using go tool pprof
	cmd := exec.Command("go", "tool", "pprof", "-http=:", "-no_browser", profilePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Generate a simple visualization using go tool pprof
	// For a basic flame graph, we'll generate the top output
	cmd = exec.Command("go", "tool", "pprof", "-top", "-cum", profilePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Fallback to simple visualization
		s.handleSimpleVisualization(w, profilePath, profileType)
		return
	}

	// Display as formatted text
	tmpl := template.Must(template.New("flamegraph").Parse(flamegraphTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{
		"Type":    profileType,
		"Content": string(output),
		"Path":    profilePath,
	})
}

// handleSimpleVisualization provides a fallback text-based visualization
func (s *Server) handleSimpleVisualization(w http.ResponseWriter, profilePath, profileType string) {
	data, err := os.ReadFile(profilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read profile: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the profile using google pprof
	prof, err := profile.Parse(bytes.NewReader(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse profile: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate a simple text summary
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Profile: %s\n", profileType))
	summary.WriteString(fmt.Sprintf("Sample Type: %v\n", prof.SampleType))
	summary.WriteString(fmt.Sprintf("Samples: %d\n\n", len(prof.Sample)))

	tmpl := template.Must(template.New("profile").Parse(profileTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{
		"Type":    profileType,
		"Profile": summary.String(),
	})
}

// handleCompare shows side-by-side profile comparison
func (s *Server) handleCompare(w http.ResponseWriter, r *http.Request, run *models.BenchmarkRun) {
	tmpl := template.Must(template.New("compare").Parse(compareTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, run)
}

// handleStatic serves static assets
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/static/"):]
	http.ServeFile(w, r, filepath.Join("static", path))
}

// HTML templates
const flamegraphTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Type}} Profile</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, monospace;
            margin: 0;
            padding: 20px;
            background: #1e1e1e;
            color: #d4d4d4;
        }
        .header {
            background: #2d2d30;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
        }
        h1 {
            margin: 0;
            color: #fff;
        }
        .actions {
            margin-top: 10px;
        }
        .btn {
            display: inline-block;
            padding: 8px 16px;
            background: #007acc;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            margin-right: 10px;
        }
        .btn:hover {
            background: #005a9e;
        }
        pre {
            background: #2d2d30;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            white-space: pre;
            font-family: 'Courier New', Courier, monospace;
            font-size: 13px;
            line-height: 1.5;
        }
        .hint {
            background: #3e3e42;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            border-left: 4px solid #007acc;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Type}} Profile</h1>
        <div class="actions">
            <a href="/" class="btn">← Back to Overview</a>
            <a href="{{.Path}}" class="btn">Download Profile</a>
        </div>
    </div>

    <div class="hint">
        <strong>💡 Tip:</strong> For interactive flame graphs, download the profile and use:<br>
        <code>go tool pprof -http=:8080 {{.Path}}</code>
    </div>

    <pre>{{.Content}}</pre>
</body>
</html>`

const indexTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Profile Viewer - {{.Run.ID}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .header {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        h1 {
            margin: 0;
            color: #333;
        }
        .meta {
            color: #666;
            margin-top: 10px;
        }
        .profiles {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        .card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .card h2 {
            margin-top: 0;
            color: #333;
        }
        .btn {
            display: inline-block;
            padding: 10px 20px;
            background: #007bff;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            margin-right: 10px;
            margin-top: 10px;
        }
        .btn:hover {
            background: #0056b3;
        }
        .summary {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .summary h2 {
            margin-top: 0;
        }
        .summary-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 15px;
            margin-top: 15px;
        }
        .stat {
            padding: 15px;
            background: #f8f9fa;
            border-radius: 4px;
            border-left: 4px solid #007bff;
        }
        .stat-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
        }
        .stat-value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔥 Profile Viewer</h1>
        <div class="meta">
            <strong>Run ID:</strong> {{.Run.ID}}<br>
            <strong>Timestamp:</strong> {{.Run.Timestamp.Format "2006-01-02 15:04:05"}}<br>
            <strong>Package:</strong> {{.Run.Package}}<br>
            <strong>Duration:</strong> {{.Run.Duration}}
        </div>
    </div>

    <div class="profiles">
        {{if .HasCPU}}
        <div class="card">
            <h2>🔥 CPU Profile</h2>
            <p>Analyze where your code spends time during execution.</p>
            <a href="/cpu/flamegraph" class="btn">View Flame Graph</a>
            <a href="/cpu" class="btn">Download Profile</a>
        </div>
        {{end}}

        {{if .HasMemory}}
        <div class="card">
            <h2>💾 Memory Profile</h2>
            <p>Identify memory allocations and potential leaks.</p>
            <a href="/mem/flamegraph" class="btn">View Flame Graph</a>
            <a href="/mem" class="btn">Download Profile</a>
        </div>
        {{end}}

        {{if and .HasCPU .HasMemory}}
        <div class="card">
            <h2>📊 Comparison</h2>
            <p>View CPU and memory profiles side-by-side.</p>
            <a href="/compare" class="btn">Compare Profiles</a>
        </div>
        {{end}}
    </div>

    {{if .HasSummary}}
    <div class="summary">
        <h2>📈 Profile Summary</h2>
        <div class="summary-grid">
            {{if .Run.ProfileSummary.TotalCPUSamples}}
            <div class="stat">
                <div class="stat-label">CPU Samples</div>
                <div class="stat-value">{{.Run.ProfileSummary.TotalCPUSamples}}</div>
            </div>
            {{end}}
            {{if .Run.ProfileSummary.TotalMemoryBytes}}
            <div class="stat">
                <div class="stat-label">Memory Allocated</div>
                <div class="stat-value">{{printf "%.1f MB" (div (float64 .Run.ProfileSummary.TotalMemoryBytes) 1048576)}}</div>
            </div>
            {{end}}
            <div class="stat">
                <div class="stat-label">Hot Functions</div>
                <div class="stat-value">{{len .Run.ProfileSummary.CPUTopFunctions}}</div>
            </div>
            <div class="stat">
                <div class="stat-label">Suggestions</div>
                <div class="stat-value">{{len .Run.ProfileSummary.Suggestions}}</div>
            </div>
        </div>
    </div>
    {{end}}
</body>
</html>`

const profileTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Type}} Profile</title>
    <style>
        body {
            font-family: monospace;
            padding: 20px;
            background: #1e1e1e;
            color: #d4d4d4;
        }
        pre {
            white-space: pre-wrap;
            word-wrap: break-word;
        }
    </style>
</head>
<body>
    <h1>{{.Type}} Profile</h1>
    <pre>{{.Profile}}</pre>
</body>
</html>`

const compareTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Profile Comparison</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 0;
        }
        .container {
            display: grid;
            grid-template-columns: 1fr 1fr;
            height: 100vh;
        }
        .pane {
            padding: 20px;
            overflow: auto;
        }
        .pane h2 {
            margin-top: 0;
        }
        .left {
            border-right: 2px solid #ccc;
        }
        iframe {
            width: 100%;
            height: calc(100% - 60px);
            border: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="pane left">
            <h2>🔥 CPU Profile</h2>
            <iframe src="/cpu/flamegraph"></iframe>
        </div>
        <div class="pane">
            <h2>💾 Memory Profile</h2>
            <iframe src="/mem/flamegraph"></iframe>
        </div>
    </div>
</body>
</html>`
