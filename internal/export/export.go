package export

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/alenon/gokanon/internal/models"
)

// Exporter handles exporting benchmark comparisons to various formats
type Exporter struct{}

// NewExporter creates a new exporter
func NewExporter() *Exporter {
	return &Exporter{}
}

// ToCSV exports comparisons to CSV format
func (e *Exporter) ToCSV(comparisons []models.Comparison, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Benchmark", "Old (ns/op)", "New (ns/op)", "Delta (ns/op)", "Delta (%)", "Status"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, comp := range comparisons {
		record := []string{
			comp.Name,
			fmt.Sprintf("%.2f", comp.OldNsPerOp),
			fmt.Sprintf("%.2f", comp.NewNsPerOp),
			fmt.Sprintf("%.2f", comp.Delta),
			fmt.Sprintf("%.2f", comp.DeltaPercent),
			comp.Status,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// ToMarkdown exports comparisons to Markdown format
func (e *Exporter) ToMarkdown(comparisons []models.Comparison, oldID, newID string, filename string) error {
	var sb strings.Builder

	sb.WriteString("# Benchmark Comparison\n\n")
	sb.WriteString(fmt.Sprintf("Comparing: `%s` vs `%s`\n\n", oldID, newID))
	sb.WriteString("| Status | Benchmark | Old (ns/op) | New (ns/op) | Delta | Delta (%) |\n")
	sb.WriteString("|--------|-----------|-------------|-------------|-------|----------|\n")

	for _, comp := range comparisons {
		status := "⚪"
		switch comp.Status {
		case "improved":
			status = "🟢"
		case "degraded":
			status = "🔴"
		}

		sb.WriteString(fmt.Sprintf("| %s | %s | %.2f | %.2f | %.2f | %+.2f%% |\n",
			status,
			comp.Name,
			comp.OldNsPerOp,
			comp.NewNsPerOp,
			comp.Delta,
			comp.DeltaPercent,
		))
	}

	// Add summary
	improved, degraded, same := countStatus(comparisons)
	sb.WriteString(fmt.Sprintf("\n## Summary\n\n"))
	sb.WriteString(fmt.Sprintf("- 🟢 Improved: %d\n", improved))
	sb.WriteString(fmt.Sprintf("- 🔴 Degraded: %d\n", degraded))
	sb.WriteString(fmt.Sprintf("- ⚪ Unchanged: %d\n", same))

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

// ToHTML exports comparisons to HTML format
func (e *Exporter) ToHTML(comparisons []models.Comparison, oldID, newID, oldTimestamp, newTimestamp string, filename string) error {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Benchmark Comparison Report</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --primary-color: #4f46e5;
            --success-color: #10b981;
            --danger-color: #ef4444;
            --warning-color: #f59e0b;
            --neutral-color: #6b7280;
            --bg-color: #f9fafb;
            --card-bg: #ffffff;
            --text-primary: #111827;
            --text-secondary: #6b7280;
            --border-color: #e5e7eb;
            --shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            --shadow-lg: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
            color: var(--text-primary);
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        header {
            background: var(--card-bg);
            border-radius: 16px;
            padding: 40px;
            margin-bottom: 30px;
            box-shadow: var(--shadow-lg);
            animation: slideDown 0.5s ease-out;
        }

        @keyframes slideDown {
            from {
                opacity: 0;
                transform: translateY(-20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        h1 {
            font-size: 2.5rem;
            font-weight: 800;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            margin-bottom: 10px;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 1.1rem;
        }

        .metadata {
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
            padding: 20px;
            border-radius: 12px;
            margin: 20px 0;
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 15px;
        }

        .metadata-item {
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .metadata-item strong {
            color: var(--text-primary);
            font-weight: 600;
        }

        .metadata-item span {
            color: var(--text-secondary);
        }

        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin: 30px 0;
        }

        .summary-card {
            background: var(--card-bg);
            padding: 30px;
            border-radius: 16px;
            text-align: center;
            box-shadow: var(--shadow);
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }

        .summary-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: linear-gradient(90deg, var(--card-color) 0%, var(--card-color-light) 100%);
        }

        .summary-card:hover {
            transform: translateY(-5px);
            box-shadow: var(--shadow-lg);
        }

        .summary-card.improved-card {
            --card-color: var(--success-color);
            --card-color-light: #34d399;
        }

        .summary-card.degraded-card {
            --card-color: var(--danger-color);
            --card-color-light: #f87171;
        }

        .summary-card.same-card {
            --card-color: var(--neutral-color);
            --card-color-light: #9ca3af;
        }

        .summary-card h3 {
            font-size: 0.875rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 1px;
            color: var(--text-secondary);
            margin-bottom: 15px;
        }

        .summary-card .number {
            font-size: 3rem;
            font-weight: 800;
            color: var(--card-color);
            line-height: 1;
        }

        .summary-card .label {
            margin-top: 10px;
            font-size: 0.875rem;
            color: var(--text-secondary);
        }

        .chart-container {
            background: var(--card-bg);
            border-radius: 16px;
            padding: 30px;
            margin: 30px 0;
            box-shadow: var(--shadow);
        }

        .chart-container h2 {
            font-size: 1.5rem;
            font-weight: 700;
            margin-bottom: 20px;
            color: var(--text-primary);
        }

        .chart-wrapper {
            position: relative;
            height: 400px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            background: var(--card-bg);
            border-radius: 16px;
            overflow: hidden;
            box-shadow: var(--shadow);
            margin: 30px 0;
        }

        thead {
            background: linear-gradient(135deg, var(--primary-color) 0%, #6366f1 100%);
        }

        th {
            color: white;
            padding: 16px;
            text-align: left;
            font-weight: 600;
            font-size: 0.875rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        tbody tr {
            border-bottom: 1px solid var(--border-color);
            transition: background-color 0.2s ease;
        }

        tbody tr:hover {
            background-color: #f9fafb;
        }

        tbody tr:last-child {
            border-bottom: none;
        }

        td {
            padding: 16px;
            font-size: 0.95rem;
        }

        .status {
            font-size: 1.5rem;
        }

        .benchmark-name {
            font-weight: 600;
            color: var(--text-primary);
        }

        .metric {
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
        }

        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.875rem;
            font-weight: 600;
        }

        .badge.improved {
            background-color: #d1fae5;
            color: var(--success-color);
        }

        .badge.degraded {
            background-color: #fee2e2;
            color: var(--danger-color);
        }

        .badge.same {
            background-color: #f3f4f6;
            color: var(--neutral-color);
        }

        .footer {
            text-align: center;
            padding: 40px 20px;
            color: white;
            font-size: 0.875rem;
        }

        .footer a {
            color: white;
            text-decoration: underline;
        }

        @media (max-width: 768px) {
            h1 {
                font-size: 1.75rem;
            }

            .summary {
                grid-template-columns: 1fr;
            }

            table {
                font-size: 0.875rem;
            }

            th, td {
                padding: 12px 8px;
            }

            .chart-wrapper {
                height: 300px;
            }
        }

        .loading {
            text-align: center;
            padding: 60px 20px;
            color: var(--text-secondary);
        }

        .spinner {
            border: 3px solid #f3f4f6;
            border-top: 3px solid var(--primary-color);
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 20px;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 Benchmark Comparison Report</h1>
            <p class="subtitle">Performance Analysis & Regression Detection</p>
        </header>

        <div class="metadata">
            <div class="metadata-item">
                <strong>📦 Old Run:</strong>
                <span>{{.OldID}} ({{.OldTimestamp}})</span>
            </div>
            <div class="metadata-item">
                <strong>📦 New Run:</strong>
                <span>{{.NewID}} ({{.NewTimestamp}})</span>
            </div>
        </div>

        <div class="summary">
            <div class="summary-card improved-card">
                <h3>Improved</h3>
                <div class="number">{{.Improved}}</div>
                <div class="label">Faster benchmarks</div>
            </div>
            <div class="summary-card degraded-card">
                <h3>Degraded</h3>
                <div class="number">{{.Degraded}}</div>
                <div class="label">Slower benchmarks</div>
            </div>
            <div class="summary-card same-card">
                <h3>Unchanged</h3>
                <div class="number">{{.Same}}</div>
                <div class="label">Stable benchmarks</div>
            </div>
        </div>

        <div class="chart-container">
            <h2>Performance Comparison</h2>
            <div class="chart-wrapper">
                <canvas id="performanceChart"></canvas>
            </div>
        </div>

        <div class="chart-container">
            <h2>Delta Distribution</h2>
            <div class="chart-wrapper">
                <canvas id="deltaChart"></canvas>
            </div>
        </div>

        <table>
            <thead>
                <tr>
                    <th>Status</th>
                    <th>Benchmark</th>
                    <th>Old (ns/op)</th>
                    <th>New (ns/op)</th>
                    <th>Delta (ns/op)</th>
                    <th>Delta (%)</th>
                </tr>
            </thead>
            <tbody>
                {{range .Comparisons}}
                <tr>
                    <td class="status">
                        {{if eq .Status "improved"}}✅{{else if eq .Status "degraded"}}❌{{else}}⚪{{end}}
                    </td>
                    <td class="benchmark-name">{{.Name}}</td>
                    <td class="metric">{{printf "%.2f" .OldNsPerOp}}</td>
                    <td class="metric">{{printf "%.2f" .NewNsPerOp}}</td>
                    <td class="metric">{{printf "%+.2f" .Delta}}</td>
                    <td>
                        <span class="badge {{.Status}}">{{printf "%+.2f%%" .DeltaPercent}}</span>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>

        <div class="footer">
            <p>Generated by <a href="https://github.com/alenon/gokanon" target="_blank">gokanon</a></p>
            <p>A powerful CLI tool for Go benchmark testing and performance analysis</p>
        </div>
    </div>

    <script>
        // Prepare data for charts
        const comparisons = [
            {{range .Comparisons}}
            {
                name: "{{.Name}}",
                oldValue: {{.OldNsPerOp}},
                newValue: {{.NewNsPerOp}},
                delta: {{.Delta}},
                deltaPercent: {{.DeltaPercent}},
                status: "{{.Status}}"
            },
            {{end}}
        ];

        // Performance Comparison Chart
        const ctx1 = document.getElementById('performanceChart').getContext('2d');
        new Chart(ctx1, {
            type: 'bar',
            data: {
                labels: comparisons.map(c => c.name.length > 30 ? c.name.substring(0, 30) + '...' : c.name),
                datasets: [
                    {
                        label: 'Old (ns/op)',
                        data: comparisons.map(c => c.oldValue),
                        backgroundColor: 'rgba(107, 114, 128, 0.7)',
                        borderColor: 'rgba(107, 114, 128, 1)',
                        borderWidth: 2
                    },
                    {
                        label: 'New (ns/op)',
                        data: comparisons.map(c => c.newValue),
                        backgroundColor: comparisons.map(c =>
                            c.status === 'improved' ? 'rgba(16, 185, 129, 0.7)' :
                            c.status === 'degraded' ? 'rgba(239, 68, 68, 0.7)' :
                            'rgba(107, 114, 128, 0.7)'
                        ),
                        borderColor: comparisons.map(c =>
                            c.status === 'improved' ? 'rgba(16, 185, 129, 1)' :
                            c.status === 'degraded' ? 'rgba(239, 68, 68, 1)' :
                            'rgba(107, 114, 128, 1)'
                        ),
                        borderWidth: 2
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'top',
                        labels: {
                            font: {
                                size: 14,
                                weight: '600'
                            }
                        }
                    },
                    tooltip: {
                        callbacks: {
                            afterLabel: function(context) {
                                const index = context.dataIndex;
                                const comp = comparisons[index];
                                return 'Delta: ' + comp.deltaPercent.toFixed(2) + '%';
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Nanoseconds per operation',
                            font: {
                                size: 14,
                                weight: '600'
                            }
                        }
                    }
                }
            }
        });

        // Delta Distribution Chart
        const ctx2 = document.getElementById('deltaChart').getContext('2d');
        new Chart(ctx2, {
            type: 'bar',
            data: {
                labels: comparisons.map(c => c.name.length > 30 ? c.name.substring(0, 30) + '...' : c.name),
                datasets: [{
                    label: 'Performance Delta (%)',
                    data: comparisons.map(c => c.deltaPercent),
                    backgroundColor: comparisons.map(c =>
                        c.deltaPercent < 0 ? 'rgba(16, 185, 129, 0.7)' :
                        c.deltaPercent > 0 ? 'rgba(239, 68, 68, 0.7)' :
                        'rgba(107, 114, 128, 0.7)'
                    ),
                    borderColor: comparisons.map(c =>
                        c.deltaPercent < 0 ? 'rgba(16, 185, 129, 1)' :
                        c.deltaPercent > 0 ? 'rgba(239, 68, 68, 1)' :
                        'rgba(107, 114, 128, 1)'
                    ),
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return 'Delta: ' + context.parsed.y.toFixed(2) + '%';
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        title: {
                            display: true,
                            text: 'Performance Change (%)',
                            font: {
                                size: 14,
                                weight: '600'
                            }
                        },
                        ticks: {
                            callback: function(value) {
                                return value + '%';
                            }
                        }
                    }
                }
            }
        });
    </script>
</body>
</html>`

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	improved, degraded, same := countStatus(comparisons)

	data := struct {
		OldID        string
		NewID        string
		OldTimestamp string
		NewTimestamp string
		Comparisons  []models.Comparison
		Improved     int
		Degraded     int
		Same         int
	}{
		OldID:        oldID,
		NewID:        newID,
		OldTimestamp: oldTimestamp,
		NewTimestamp: newTimestamp,
		Comparisons:  comparisons,
		Improved:     improved,
		Degraded:     degraded,
		Same:         same,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	return t.Execute(file, data)
}

// countStatus counts the number of each status type
func countStatus(comparisons []models.Comparison) (improved, degraded, same int) {
	for _, comp := range comparisons {
		switch comp.Status {
		case "improved":
			improved++
		case "degraded":
			degraded++
		case "same":
			same++
		}
	}
	return
}
