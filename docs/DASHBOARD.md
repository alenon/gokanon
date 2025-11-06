# Interactive Web Dashboard

The GoKanon interactive web dashboard provides a powerful, real-time visualization interface for analyzing benchmark results, tracking performance trends, and sharing insights with your team.

## Features

### 🎯 Overview Dashboard
- **Real-time Statistics**: View total runs, tests, and unique benchmarks at a glance
- **Recent Performance Chart**: Visualize average performance across recent benchmark runs
- **Quick Access**: Jump directly to any benchmark run from the recent runs list

### 📊 Trend Analysis
- **Historical Performance Graphs**: Track how your benchmarks perform over time with interactive line charts
- **Statistical Insights**: View mean, median, standard deviation, and coefficient of variation for each benchmark
- **Trend Detection**: Automatically identify improving, degrading, or stable performance trends
- **Customizable Time Range**: Filter by number of runs (10, 25, 50, or 100 runs)
- **Benchmark Filtering**: Focus on specific benchmarks or view all at once

### 📜 History View
- **Complete Run History**: Browse all benchmark runs in a sortable, filterable table
- **Quick Search**: Filter runs by package name or run ID
- **Detailed Metrics**: View timestamp, Go version, test count, and average performance for each run
- **Click-to-View**: Click any run to see detailed information

### 🔍 Search & Filter
- **Global Search**: Search across all runs, packages, and benchmark names
- **Smart Results**: Search results categorized by type (runs vs. benchmarks)
- **Real-time Filtering**: Instant results as you type

### 🔄 Comparison Mode
- **Side-by-Side Comparison**: Compare any two benchmark runs
- **Performance Delta**: See exact performance improvements or degradations
- **Visual Indicators**: Color-coded results show improvements (green) and degradations (red)
- **Percentage Changes**: Understand the magnitude of performance changes

### 🌙 Dark Mode
- **Eye-Friendly**: Toggle between light and dark themes
- **Persistent Preference**: Your theme choice is saved automatically
- **Chart Adaptation**: All charts automatically adjust colors for optimal visibility

### 🔗 Shareable URLs
- **Direct Links**: Share specific benchmark runs via URL
- **Embed Mode**: Embed dashboard views in documentation sites
- **State Preservation**: URLs preserve your current view (tab, run selection, etc.)

### 📱 Responsive Design
- **Mobile-Friendly**: Works seamlessly on desktop, tablet, and mobile devices
- **Adaptive Layout**: UI elements reorganize for optimal viewing on any screen size

## Getting Started

### Start the Dashboard

```bash
# Start on default port (8080)
gokanon serve

# Start on custom port
gokanon serve -port=9000

# Bind to all interfaces (for remote access)
gokanon serve -addr=0.0.0.0 -port=8080

# Use custom storage directory
gokanon serve -storage=/path/to/.gokanon
```

### Access the Dashboard

Once started, open your browser and navigate to:
- Local access: `http://localhost:8080`
- Remote access: `http://<your-ip>:8080`

## Usage Examples

### Analyzing Performance Trends

1. Navigate to the **Trends** tab
2. Select a specific benchmark from the dropdown (or view all)
3. Choose the number of recent runs to analyze
4. Click **Load Trends** to see:
   - Line charts showing performance over time
   - Statistical summaries with trend direction
   - Coefficient of variation for stability assessment

### Comparing Benchmark Runs

1. Go to the **Compare** tab
2. Select a baseline run from the first dropdown
3. Select a comparison run from the second dropdown
4. Click **Compare** to see:
   - Performance deltas for each benchmark
   - Percentage improvements or degradations
   - Color-coded indicators

### Searching for Specific Results

1. Use the search bar at the top of the dashboard
2. Enter any of:
   - Package name (e.g., "mypackage")
   - Run ID (e.g., "run-123")
   - Benchmark name (e.g., "BenchmarkFoo")
3. Results appear instantly, categorized by type

### Sharing Results

To share a specific benchmark run:

1. Click on the run in the History or Overview tab
2. Copy the URL from your browser
3. Share the URL with your team

The URL includes the run ID and will load that specific run for anyone who opens it.

### Embedding in Documentation

To embed the dashboard in documentation:

1. Add `?embed=true` to any dashboard URL
2. Use an iframe in your documentation:

```html
<iframe
  src="http://localhost:8080?embed=true&tab=trends"
  width="100%"
  height="600px"
  frameborder="0">
</iframe>
```

Embed mode provides:
- A notice banner with a link to the full dashboard
- Hides the header and footer for cleaner embedding
- Maintains all interactive functionality

## API Endpoints

The dashboard exposes the following REST API endpoints:

### Get All Runs
```
GET /api/runs
```
Returns a summary of all benchmark runs.

### Get Run Details
```
GET /api/runs/{id}
```
Returns detailed information for a specific run.

### Get Trends
```
GET /api/trends?benchmark={name}&limit={number}
```
Returns trend data across multiple runs.

Parameters:
- `benchmark` (optional): Filter by benchmark name
- `limit` (optional): Number of runs to include (default: 50)

### Get Statistics
```
GET /api/stats
```
Returns aggregate statistics across all runs.

### Search
```
GET /api/search?q={query}
```
Searches across runs and benchmarks.

Parameters:
- `q`: Search query string

## Chart Types

### Overview Chart
- **Type**: Line chart
- **Data**: Average ns/op across recent runs
- **Purpose**: Quick visualization of recent performance

### Trends Chart
- **Type**: Multi-line chart with time-based X-axis
- **Data**: ns/op for selected benchmarks over time
- **Features**:
  - Multiple benchmarks on same chart
  - Hover tooltips for exact values
  - Automatic color coding

## Performance Considerations

The dashboard is optimized for:
- **Fast Loading**: Minimal initial data transfer
- **Lazy Loading**: Trends data loaded on-demand
- **Efficient Rendering**: Chart.js with hardware acceleration
- **Responsive Updates**: Instant UI feedback

## Browser Compatibility

The dashboard works with all modern browsers:
- Chrome/Edge (v90+)
- Firefox (v88+)
- Safari (v14+)

Requires JavaScript enabled and HTML5 support.

## Keyboard Shortcuts

- `?` - Show help (planned)
- `d` - Toggle dark mode (planned)
- `r` - Refresh data (planned)

## Tips & Best Practices

1. **Regular Benchmarking**: Run benchmarks frequently to build a comprehensive trend history
2. **Consistent Environment**: Run benchmarks in similar conditions for accurate comparisons
3. **Annotation**: Use commit messages or tags to correlate performance changes with code changes
4. **Threshold Monitoring**: Use the Compare tab to validate performance improvements
5. **Share Insights**: Use shareable URLs to collaborate with your team

## Troubleshooting

### Dashboard Won't Start

**Issue**: Port already in use
```
Error: failed to start dashboard server: listen tcp :8080: bind: address already in use
```

**Solution**: Use a different port
```bash
gokanon serve -port=9000
```

### No Data Showing

**Issue**: Empty dashboard with no runs

**Solution**: Ensure you have benchmark results saved
```bash
gokanon run
```

### Can't Access Remotely

**Issue**: Dashboard only accessible on localhost

**Solution**: Bind to all interfaces
```bash
gokanon serve -addr=0.0.0.0
```

### Charts Not Rendering

**Issue**: Blank chart areas

**Solution**:
- Check browser console for JavaScript errors
- Ensure JavaScript is enabled
- Try clearing browser cache

## Future Enhancements

Planned features for future releases:
- Real-time benchmark streaming (WebSocket support)
- Heat maps for visualizing performance across multiple dimensions
- Performance annotations linked to git commits
- Export dashboard as PDF/PNG
- Custom dashboard layouts
- Benchmark comparison across branches
- Alert notifications for performance regressions

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines on contributing to the dashboard.

## License

The dashboard is part of GoKanon and is distributed under the same license. See [LICENSE](../LICENSE) for details.
