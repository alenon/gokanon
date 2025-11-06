package dashboard

// indexHTML is the main dashboard HTML page
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoKanon Dashboard</title>
    <link rel="stylesheet" href="/static/styles.css">
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
</head>
<body>
    <div class="dashboard-container">
        <!-- Header -->
        <header class="header">
            <div class="header-content">
                <h1>📊 GoKanon Dashboard</h1>
                <div class="header-controls">
                    <button id="darkModeToggle" class="btn btn-icon" title="Toggle dark mode">
                        <span class="icon-sun">☀️</span>
                        <span class="icon-moon">🌙</span>
                    </button>
                    <button id="refreshBtn" class="btn btn-primary" title="Refresh data">
                        🔄 Refresh
                    </button>
                </div>
            </div>
        </header>

        <!-- Embed Mode Notice -->
        <div id="embedNotice" class="embed-notice" style="display: none;">
            <span>Embedded View</span>
            <a href="/" target="_blank">Open Full Dashboard</a>
        </div>

        <!-- Main Content -->
        <main class="main-content">
            <!-- Stats Overview -->
            <section class="stats-section">
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-icon">🏃</div>
                        <div class="stat-content">
                            <div class="stat-value" id="totalRuns">-</div>
                            <div class="stat-label">Total Runs</div>
                        </div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon">📝</div>
                        <div class="stat-content">
                            <div class="stat-value" id="totalTests">-</div>
                            <div class="stat-label">Total Tests</div>
                        </div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon">📦</div>
                        <div class="stat-content">
                            <div class="stat-value" id="totalBenchmarks">-</div>
                            <div class="stat-label">Unique Benchmarks</div>
                        </div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon">📅</div>
                        <div class="stat-content">
                            <div class="stat-value" id="dateRange">-</div>
                            <div class="stat-label">Date Range</div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Search and Filter -->
            <section class="search-section">
                <div class="search-bar">
                    <input type="text" id="searchInput" placeholder="Search benchmarks, packages, or run IDs..." />
                    <button id="searchBtn" class="btn btn-primary">🔍 Search</button>
                </div>
                <div id="searchResults" class="search-results"></div>
            </section>

            <!-- Tabs -->
            <section class="tabs-section">
                <div class="tabs">
                    <button class="tab-btn active" data-tab="overview">Overview</button>
                    <button class="tab-btn" data-tab="trends">Trends</button>
                    <button class="tab-btn" data-tab="history">History</button>
                    <button class="tab-btn" data-tab="compare">Compare</button>
                </div>

                <!-- Tab Content -->
                <div class="tab-content">
                    <!-- Overview Tab -->
                    <div id="overview" class="tab-pane active">
                        <div class="chart-container">
                            <h2>Recent Benchmark Performance</h2>
                            <canvas id="overviewChart"></canvas>
                        </div>
                        <div class="recent-runs">
                            <h2>Recent Runs</h2>
                            <div id="recentRunsList"></div>
                        </div>
                    </div>

                    <!-- Trends Tab -->
                    <div id="trends" class="tab-pane">
                        <div class="trends-controls">
                            <label for="benchmarkSelect">Select Benchmark:</label>
                            <select id="benchmarkSelect" class="form-select">
                                <option value="">All Benchmarks</option>
                            </select>
                            <label for="limitSelect">Show Last:</label>
                            <select id="limitSelect" class="form-select">
                                <option value="10">10 runs</option>
                                <option value="25">25 runs</option>
                                <option value="50" selected>50 runs</option>
                                <option value="100">100 runs</option>
                            </select>
                            <button id="loadTrendsBtn" class="btn btn-primary">Load Trends</button>
                        </div>
                        <div class="chart-container">
                            <h2>Performance Trends</h2>
                            <canvas id="trendsChart"></canvas>
                        </div>
                        <div class="trends-stats" id="trendsStats"></div>
                    </div>

                    <!-- History Tab -->
                    <div id="history" class="tab-pane">
                        <div class="history-controls">
                            <input type="text" id="historyFilter" placeholder="Filter by package or ID..." />
                        </div>
                        <div id="historyTable" class="table-container"></div>
                    </div>

                    <!-- Compare Tab -->
                    <div id="compare" class="tab-pane">
                        <div class="compare-controls">
                            <div class="compare-select-group">
                                <label for="compareRun1">Baseline Run:</label>
                                <select id="compareRun1" class="form-select"></select>
                            </div>
                            <div class="compare-select-group">
                                <label for="compareRun2">Compare With:</label>
                                <select id="compareRun2" class="form-select"></select>
                            </div>
                            <button id="compareBtn" class="btn btn-primary">Compare</button>
                        </div>
                        <div id="compareResults" class="compare-results"></div>
                    </div>
                </div>
            </section>

            <!-- Share Modal -->
            <div id="shareModal" class="modal">
                <div class="modal-content">
                    <div class="modal-header">
                        <h2>Share This View</h2>
                        <button class="modal-close">&times;</button>
                    </div>
                    <div class="modal-body">
                        <div class="share-options">
                            <div class="share-option">
                                <label>Direct Link:</label>
                                <input type="text" id="shareUrl" readonly />
                                <button id="copyUrlBtn" class="btn btn-secondary">Copy</button>
                            </div>
                            <div class="share-option">
                                <label>Embed Code:</label>
                                <textarea id="embedCode" readonly rows="3"></textarea>
                                <button id="copyEmbedBtn" class="btn btn-secondary">Copy</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </main>

        <!-- Footer -->
        <footer class="footer">
            <p>GoKanon Dashboard v1.0 | <a href="https://github.com/alenon/gokanon" target="_blank">GitHub</a></p>
        </footer>
    </div>

    <script src="/static/app.js"></script>
</body>
</html>`

// stylesCSS is the dashboard CSS with dark mode support
const stylesCSS = `
:root {
    --bg-primary: #ffffff;
    --bg-secondary: #f8f9fa;
    --bg-card: #ffffff;
    --text-primary: #212529;
    --text-secondary: #6c757d;
    --border-color: #dee2e6;
    --accent-color: #0d6efd;
    --accent-hover: #0b5ed7;
    --success-color: #198754;
    --danger-color: #dc3545;
    --warning-color: #ffc107;
    --shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    --shadow-lg: 0 4px 12px rgba(0, 0, 0, 0.15);
}

[data-theme="dark"] {
    --bg-primary: #1a1a1a;
    --bg-secondary: #2d2d2d;
    --bg-card: #252525;
    --text-primary: #e9ecef;
    --text-secondary: #adb5bd;
    --border-color: #404040;
    --accent-color: #4dabf7;
    --accent-hover: #339af0;
    --success-color: #51cf66;
    --danger-color: #ff6b6b;
    --warning-color: #ffd43b;
    --shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    --shadow-lg: 0 4px 12px rgba(0, 0, 0, 0.5);
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    background-color: var(--bg-primary);
    color: var(--text-primary);
    line-height: 1.6;
    transition: background-color 0.3s, color 0.3s;
}

.dashboard-container {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}

/* Header */
.header {
    background-color: var(--bg-card);
    border-bottom: 1px solid var(--border-color);
    padding: 1rem 2rem;
    box-shadow: var(--shadow);
}

.header-content {
    max-width: 1400px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.header h1 {
    font-size: 1.8rem;
    font-weight: 600;
}

.header-controls {
    display: flex;
    gap: 0.5rem;
}

/* Buttons */
.btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    background-color: var(--bg-secondary);
    color: var(--text-primary);
}

.btn:hover {
    opacity: 0.9;
    transform: translateY(-1px);
}

.btn-primary {
    background-color: var(--accent-color);
    color: white;
}

.btn-primary:hover {
    background-color: var(--accent-hover);
}

.btn-secondary {
    background-color: var(--text-secondary);
    color: white;
}

.btn-icon {
    padding: 0.5rem;
    font-size: 1.2rem;
}

[data-theme="light"] .icon-moon,
[data-theme="dark"] .icon-sun {
    display: none;
}

/* Embed Notice */
.embed-notice {
    background-color: var(--warning-color);
    color: #000;
    padding: 0.5rem 2rem;
    text-align: center;
    font-size: 0.9rem;
}

.embed-notice a {
    color: #000;
    font-weight: 600;
    margin-left: 1rem;
}

/* Main Content */
.main-content {
    flex: 1;
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
    width: 100%;
}

/* Stats Section */
.stats-section {
    margin-bottom: 2rem;
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1rem;
}

.stat-card {
    background-color: var(--bg-card);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 1.5rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    box-shadow: var(--shadow);
    transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-lg);
}

.stat-icon {
    font-size: 2.5rem;
}

.stat-content {
    flex: 1;
}

.stat-value {
    font-size: 2rem;
    font-weight: 700;
    color: var(--accent-color);
}

.stat-label {
    font-size: 0.9rem;
    color: var(--text-secondary);
}

/* Search Section */
.search-section {
    margin-bottom: 2rem;
}

.search-bar {
    display: flex;
    gap: 0.5rem;
}

.search-bar input {
    flex: 1;
    padding: 0.75rem;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    font-size: 1rem;
    background-color: var(--bg-card);
    color: var(--text-primary);
}

.search-results {
    margin-top: 1rem;
    background-color: var(--bg-card);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    max-height: 400px;
    overflow-y: auto;
}

.search-result-item {
    padding: 1rem;
    border-bottom: 1px solid var(--border-color);
    cursor: pointer;
    transition: background-color 0.2s;
}

.search-result-item:hover {
    background-color: var(--bg-secondary);
}

.search-result-item:last-child {
    border-bottom: none;
}

/* Tabs */
.tabs-section {
    background-color: var(--bg-card);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    box-shadow: var(--shadow);
}

.tabs {
    display: flex;
    border-bottom: 1px solid var(--border-color);
    padding: 0 1rem;
    gap: 0.5rem;
}

.tab-btn {
    padding: 1rem 1.5rem;
    border: none;
    background: none;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    color: var(--text-secondary);
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
}

.tab-btn:hover {
    color: var(--text-primary);
}

.tab-btn.active {
    color: var(--accent-color);
    border-bottom-color: var(--accent-color);
}

.tab-content {
    padding: 2rem;
}

.tab-pane {
    display: none;
}

.tab-pane.active {
    display: block;
}

/* Charts */
.chart-container {
    margin-bottom: 2rem;
}

.chart-container h2 {
    margin-bottom: 1rem;
    font-size: 1.5rem;
}

.chart-container canvas {
    max-height: 400px;
}

/* Recent Runs */
.recent-runs h2 {
    margin-bottom: 1rem;
    font-size: 1.5rem;
}

.run-item {
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 6px;
    margin-bottom: 0.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    transition: all 0.2s;
}

.run-item:hover {
    background-color: var(--border-color);
    transform: translateX(4px);
}

/* Trends Controls */
.trends-controls {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
    align-items: center;
}

.trends-controls label {
    font-weight: 500;
}

.form-select {
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    background-color: var(--bg-card);
    color: var(--text-primary);
    font-size: 0.9rem;
}

.trends-stats {
    margin-top: 2rem;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
}

.trend-stat-card {
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 6px;
    border-left: 4px solid var(--accent-color);
}

.trend-stat-card h3 {
    font-size: 1rem;
    margin-bottom: 0.5rem;
}

.trend-stat-card.improving {
    border-left-color: var(--success-color);
}

.trend-stat-card.degrading {
    border-left-color: var(--danger-color);
}

/* History Controls */
.history-controls {
    margin-bottom: 1rem;
}

.history-controls input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    background-color: var(--bg-card);
    color: var(--text-primary);
}

/* Table */
.table-container {
    overflow-x: auto;
}

table {
    width: 100%;
    border-collapse: collapse;
}

th, td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
}

th {
    font-weight: 600;
    background-color: var(--bg-secondary);
}

tr:hover {
    background-color: var(--bg-secondary);
}

/* Compare Controls */
.compare-controls {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
}

.compare-select-group {
    flex: 1;
    min-width: 250px;
}

.compare-select-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
}

.compare-results {
    margin-top: 2rem;
}

.comparison-item {
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 6px;
    margin-bottom: 0.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.delta-improved {
    color: var(--success-color);
    font-weight: 600;
}

.delta-degraded {
    color: var(--danger-color);
    font-weight: 600;
}

.delta-same {
    color: var(--text-secondary);
}

/* Modal */
.modal {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 1000;
    align-items: center;
    justify-content: center;
}

.modal.active {
    display: flex;
}

.modal-content {
    background-color: var(--bg-card);
    border-radius: 8px;
    width: 90%;
    max-width: 600px;
    box-shadow: var(--shadow-lg);
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--border-color);
}

.modal-close {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text-primary);
}

.modal-body {
    padding: 1.5rem;
}

.share-option {
    margin-bottom: 1.5rem;
}

.share-option label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
}

.share-option input,
.share-option textarea {
    width: calc(100% - 90px);
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    background-color: var(--bg-secondary);
    color: var(--text-primary);
    font-family: monospace;
    margin-right: 0.5rem;
}

/* Footer */
.footer {
    background-color: var(--bg-card);
    border-top: 1px solid var(--border-color);
    padding: 1.5rem 2rem;
    text-align: center;
    color: var(--text-secondary);
}

.footer a {
    color: var(--accent-color);
    text-decoration: none;
}

.footer a:hover {
    text-decoration: underline;
}

/* Responsive */
@media (max-width: 768px) {
    .header-content {
        flex-direction: column;
        gap: 1rem;
    }

    .stats-grid {
        grid-template-columns: 1fr;
    }

    .tabs {
        overflow-x: auto;
    }

    .main-content {
        padding: 1rem;
    }
}

/* Loading Animation */
@keyframes spin {
    to { transform: rotate(360deg); }
}

.loading {
    animation: spin 1s linear infinite;
}
`
