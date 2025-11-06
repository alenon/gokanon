package dashboard

// appJS is the dashboard JavaScript application
const appJS = `
// Dashboard App
const App = {
    charts: {},
    data: {
        runs: [],
        stats: null,
        trends: null,
        selectedRun: null
    },

    init() {
        this.setupEventListeners();
        this.checkEmbedMode();
        this.loadTheme();
        this.loadData();
        this.loadURLParams();
    },

    setupEventListeners() {
        // Dark mode toggle
        document.getElementById('darkModeToggle').addEventListener('click', () => {
            this.toggleTheme();
        });

        // Refresh button
        document.getElementById('refreshBtn').addEventListener('click', () => {
            this.loadData();
        });

        // Tab switching
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.switchTab(e.target.dataset.tab);
            });
        });

        // Search
        document.getElementById('searchBtn').addEventListener('click', () => {
            this.performSearch();
        });

        document.getElementById('searchInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.performSearch();
            }
        });

        // Trends controls
        document.getElementById('loadTrendsBtn').addEventListener('click', () => {
            this.loadTrends();
        });

        // History filter
        document.getElementById('historyFilter').addEventListener('input', (e) => {
            this.filterHistory(e.target.value);
        });

        // Compare
        document.getElementById('compareBtn').addEventListener('click', () => {
            this.compareRuns();
        });
    },

    checkEmbedMode() {
        const params = new URLSearchParams(window.location.search);
        if (params.get('embed') === 'true') {
            document.getElementById('embedNotice').style.display = 'block';
            document.querySelector('.header').style.display = 'none';
            document.querySelector('.footer').style.display = 'none';
        }
    },

    loadTheme() {
        const theme = localStorage.getItem('theme') || 'light';
        document.documentElement.setAttribute('data-theme', theme);
    },

    toggleTheme() {
        const current = document.documentElement.getAttribute('data-theme');
        const next = current === 'dark' ? 'light' : 'dark';
        document.documentElement.setAttribute('data-theme', next);
        localStorage.setItem('theme', next);

        // Update charts for new theme
        if (this.charts.overview) {
            this.updateChartTheme();
        }
    },

    updateChartTheme() {
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const textColor = isDark ? '#e9ecef' : '#212529';
        const gridColor = isDark ? '#404040' : '#dee2e6';

        Object.values(this.charts).forEach(chart => {
            if (chart && chart.options) {
                chart.options.scales.x.ticks.color = textColor;
                chart.options.scales.y.ticks.color = textColor;
                chart.options.scales.x.grid.color = gridColor;
                chart.options.scales.y.grid.color = gridColor;
                chart.options.plugins.legend.labels.color = textColor;
                chart.update();
            }
        });
    },

    async loadData() {
        try {
            // Load stats
            const statsRes = await fetch('/api/stats');
            this.data.stats = await statsRes.json();
            this.updateStats();

            // Load runs
            const runsRes = await fetch('/api/runs');
            this.data.runs = await runsRes.json();
            this.updateRecentRuns();
            this.createOverviewChart();
            this.populateCompareSelects();
            this.populateBenchmarkSelect();
            this.updateHistory();
        } catch (error) {
            console.error('Failed to load data:', error);
            alert('Failed to load dashboard data. Please check if the server is running.');
        }
    },

    updateStats() {
        const stats = this.data.stats;
        document.getElementById('totalRuns').textContent = stats.totalRuns || 0;
        document.getElementById('totalTests').textContent = stats.totalTests || 0;
        document.getElementById('totalBenchmarks').textContent = stats.benchmarks?.length || 0;

        if (stats.dateRange && stats.dateRange.oldest) {
            const oldest = new Date(stats.dateRange.oldest);
            const newest = new Date(stats.dateRange.newest);
            const days = Math.ceil((newest - oldest) / (1000 * 60 * 60 * 24));
            document.getElementById('dateRange').textContent = days + ' days';
        } else {
            document.getElementById('dateRange').textContent = 'N/A';
        }
    },

    updateRecentRuns() {
        const container = document.getElementById('recentRunsList');
        const runs = this.data.stats.recentRuns || [];

        if (runs.length === 0) {
            container.innerHTML = '<p>No benchmark runs found. Run some benchmarks to get started!</p>';
            return;
        }

        container.innerHTML = runs.map(run => {
            const date = new Date(run.timestamp);
            return '<div class="run-item" onclick="App.viewRun(\'' + run.id + '\')">' +
                '<div>' +
                '<strong>' + run.package + '</strong><br>' +
                '<small>' + run.numTests + ' tests</small>' +
                '</div>' +
                '<div>' +
                '<small>' + date.toLocaleString() + '</small>' +
                '</div>' +
                '</div>';
        }).join('');
    },

    createOverviewChart() {
        const runs = this.data.runs.slice(0, 10).reverse();

        if (runs.length === 0) {
            return;
        }

        const labels = runs.map(run => {
            const date = new Date(run.timestamp);
            return date.toLocaleDateString();
        });

        const avgData = runs.map(run => run.avgNsPerOp || 0);

        const ctx = document.getElementById('overviewChart');
        if (this.charts.overview) {
            this.charts.overview.destroy();
        }

        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const textColor = isDark ? '#e9ecef' : '#212529';
        const gridColor = isDark ? '#404040' : '#dee2e6';

        this.charts.overview = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Avg ns/op',
                    data: avgData,
                    borderColor: '#4dabf7',
                    backgroundColor: 'rgba(77, 171, 247, 0.1)',
                    tension: 0.4,
                    fill: true
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        labels: { color: textColor }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return 'Avg: ' + context.parsed.y.toFixed(2) + ' ns/op';
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: { color: textColor },
                        grid: { color: gridColor }
                    },
                    x: {
                        ticks: { color: textColor },
                        grid: { color: gridColor }
                    }
                }
            }
        });
    },

    async loadTrends() {
        const benchmark = document.getElementById('benchmarkSelect').value;
        const limit = document.getElementById('limitSelect').value;

        try {
            const url = '/api/trends?limit=' + limit + (benchmark ? '&benchmark=' + encodeURIComponent(benchmark) : '');
            const res = await fetch(url);
            this.data.trends = await res.json();
            this.createTrendsChart();
            this.updateTrendsStats();
        } catch (error) {
            console.error('Failed to load trends:', error);
        }
    },

    createTrendsChart() {
        const trends = this.data.trends.trends;
        if (!trends || Object.keys(trends).length === 0) {
            return;
        }

        const colors = ['#4dabf7', '#51cf66', '#ff6b6b', '#ffd43b', '#a78bfa', '#fb923c'];
        const datasets = [];
        let colorIndex = 0;

        for (const [name, points] of Object.entries(trends)) {
            if (points.length === 0) continue;

            datasets.push({
                label: name,
                data: points.map(p => ({
                    x: new Date(p.timestamp),
                    y: p.nsPerOp
                })),
                borderColor: colors[colorIndex % colors.length],
                backgroundColor: colors[colorIndex % colors.length] + '33',
                tension: 0.4,
                fill: false
            });
            colorIndex++;
        }

        const ctx = document.getElementById('trendsChart');
        if (this.charts.trends) {
            this.charts.trends.destroy();
        }

        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const textColor = isDark ? '#e9ecef' : '#212529';
        const gridColor = isDark ? '#404040' : '#dee2e6';

        this.charts.trends = new Chart(ctx, {
            type: 'line',
            data: { datasets: datasets },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        labels: { color: textColor }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return context.dataset.label + ': ' + context.parsed.y.toFixed(2) + ' ns/op';
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        type: 'time',
                        time: {
                            unit: 'day',
                            displayFormats: {
                                day: 'MMM d'
                            }
                        },
                        ticks: { color: textColor },
                        grid: { color: gridColor }
                    },
                    y: {
                        beginAtZero: true,
                        ticks: { color: textColor },
                        grid: { color: gridColor }
                    }
                }
            }
        });
    },

    updateTrendsStats() {
        const stats = this.data.trends.statistics;
        if (!stats) return;

        const container = document.getElementById('trendsStats');
        container.innerHTML = '';

        for (const [name, stat] of Object.entries(stats)) {
            const trendClass = stat.trend === 'improving' ? 'improving' :
                             stat.trend === 'degrading' ? 'degrading' : '';

            const card = document.createElement('div');
            card.className = 'trend-stat-card ' + trendClass;
            card.innerHTML = '<h3>' + name + '</h3>' +
                '<p><strong>Mean:</strong> ' + stat.mean.toFixed(2) + ' ns/op</p>' +
                '<p><strong>Median:</strong> ' + stat.median.toFixed(2) + ' ns/op</p>' +
                '<p><strong>Std Dev:</strong> ' + stat.stdDev.toFixed(2) + '</p>' +
                '<p><strong>CV:</strong> ' + (stat.cv * 100).toFixed(2) + '%</p>' +
                '<p><strong>Trend:</strong> ' + stat.trend + '</p>';

            container.appendChild(card);
        }
    },

    populateBenchmarkSelect() {
        const select = document.getElementById('benchmarkSelect');
        const benchmarks = this.data.stats.benchmarks || [];

        select.innerHTML = '<option value="">All Benchmarks</option>';
        benchmarks.forEach(name => {
            const option = document.createElement('option');
            option.value = name;
            option.textContent = name;
            select.appendChild(option);
        });
    },

    updateHistory() {
        const container = document.getElementById('historyTable');
        const runs = this.data.runs;

        if (runs.length === 0) {
            container.innerHTML = '<p>No benchmark runs found.</p>';
            return;
        }

        let html = '<table><thead><tr>' +
            '<th>ID</th>' +
            '<th>Timestamp</th>' +
            '<th>Package</th>' +
            '<th>Go Version</th>' +
            '<th>Tests</th>' +
            '<th>Avg ns/op</th>' +
            '</tr></thead><tbody>';

        runs.forEach(run => {
            const date = new Date(run.timestamp);
            html += '<tr onclick="App.viewRun(\'' + run.id + '\')">' +
                '<td>' + run.id.substring(0, 8) + '</td>' +
                '<td>' + date.toLocaleString() + '</td>' +
                '<td>' + run.package + '</td>' +
                '<td>' + run.goVersion + '</td>' +
                '<td>' + run.numTests + '</td>' +
                '<td>' + (run.avgNsPerOp ? run.avgNsPerOp.toFixed(2) : 'N/A') + '</td>' +
                '</tr>';
        });

        html += '</tbody></table>';
        container.innerHTML = html;
    },

    filterHistory(query) {
        const lowerQuery = query.toLowerCase();
        const rows = document.querySelectorAll('#historyTable tbody tr');

        rows.forEach(row => {
            const text = row.textContent.toLowerCase();
            row.style.display = text.includes(lowerQuery) ? '' : 'none';
        });
    },

    populateCompareSelects() {
        const runs = this.data.runs;
        const select1 = document.getElementById('compareRun1');
        const select2 = document.getElementById('compareRun2');

        select1.innerHTML = '';
        select2.innerHTML = '';

        runs.forEach(run => {
            const date = new Date(run.timestamp);
            const text = run.id.substring(0, 8) + ' - ' + run.package + ' (' + date.toLocaleDateString() + ')';

            const option1 = document.createElement('option');
            option1.value = run.id;
            option1.textContent = text;
            select1.appendChild(option1);

            const option2 = document.createElement('option');
            option2.value = run.id;
            option2.textContent = text;
            select2.appendChild(option2);
        });

        // Select first and second by default
        if (runs.length >= 2) {
            select1.selectedIndex = 0;
            select2.selectedIndex = 1;
        }
    },

    async compareRuns() {
        const id1 = document.getElementById('compareRun1').value;
        const id2 = document.getElementById('compareRun2').value;

        if (!id1 || !id2) {
            alert('Please select two runs to compare');
            return;
        }

        if (id1 === id2) {
            alert('Please select two different runs');
            return;
        }

        try {
            const [run1Res, run2Res] = await Promise.all([
                fetch('/api/runs/' + id1),
                fetch('/api/runs/' + id2)
            ]);

            const run1 = await run1Res.json();
            const run2 = await run2Res.json();

            this.displayComparison(run1, run2);
        } catch (error) {
            console.error('Failed to compare runs:', error);
            alert('Failed to load run data');
        }
    },

    displayComparison(run1, run2) {
        const container = document.getElementById('compareResults');

        // Create a map of benchmarks
        const benchMap = new Map();

        run1.Results.forEach(result => {
            benchMap.set(result.Name, { old: result, new: null });
        });

        run2.Results.forEach(result => {
            if (benchMap.has(result.Name)) {
                benchMap.get(result.Name).new = result;
            } else {
                benchMap.set(result.Name, { old: null, new: result });
            }
        });

        let html = '<h3>Comparison Results</h3>';
        html += '<p>Baseline: ' + run1.ID.substring(0, 8) + ' vs ' + run2.ID.substring(0, 8) + '</p>';

        benchMap.forEach((data, name) => {
            if (!data.old || !data.new) return;

            const delta = data.new.NsPerOp - data.old.NsPerOp;
            const deltaPercent = (delta / data.old.NsPerOp) * 100;

            let deltaClass = 'delta-same';
            let deltaText = 'No change';

            if (Math.abs(deltaPercent) > 5) {
                if (deltaPercent < 0) {
                    deltaClass = 'delta-improved';
                    deltaText = deltaPercent.toFixed(2) + '% faster';
                } else {
                    deltaClass = 'delta-degraded';
                    deltaText = '+' + deltaPercent.toFixed(2) + '% slower';
                }
            }

            html += '<div class="comparison-item">' +
                '<div><strong>' + name + '</strong></div>' +
                '<div class="' + deltaClass + '">' + deltaText + '</div>' +
                '</div>';
        });

        container.innerHTML = html;
    },

    async performSearch() {
        const query = document.getElementById('searchInput').value.trim();
        if (!query) return;

        try {
            const res = await fetch('/api/search?q=' + encodeURIComponent(query));
            const data = await res.json();
            this.displaySearchResults(data);
        } catch (error) {
            console.error('Search failed:', error);
        }
    },

    displaySearchResults(data) {
        const container = document.getElementById('searchResults');

        if (data.count === 0) {
            container.innerHTML = '<div class="search-result-item">No results found</div>';
            return;
        }

        container.innerHTML = data.results.map(result => {
            const date = new Date(result.timestamp);
            if (result.type === 'run') {
                return '<div class="search-result-item" onclick="App.viewRun(\'' + result.id + '\')">' +
                    '<strong>Run: ' + result.id.substring(0, 8) + '</strong><br>' +
                    '<small>' + result.package + ' - ' + date.toLocaleString() + '</small>' +
                    '</div>';
            } else {
                return '<div class="search-result-item" onclick="App.viewRun(\'' + result.runId + '\')">' +
                    '<strong>Benchmark: ' + result.name + '</strong><br>' +
                    '<small>' + result.nsPerOp.toFixed(2) + ' ns/op - ' + date.toLocaleString() + '</small>' +
                    '</div>';
            }
        }).join('');
    },

    async viewRun(id) {
        try {
            const res = await fetch('/api/runs/' + id);
            const run = await res.json();

            // Switch to history tab and highlight
            this.switchTab('history');

            // Update URL for sharing
            const url = new URL(window.location);
            url.searchParams.set('run', id);
            window.history.pushState({}, '', url);

            alert('Run Details:\\n' +
                'ID: ' + run.ID + '\\n' +
                'Package: ' + run.Package + '\\n' +
                'Tests: ' + run.Results.length + '\\n' +
                'Go Version: ' + run.GoVersion);
        } catch (error) {
            console.error('Failed to load run:', error);
        }
    },

    switchTab(tabName) {
        // Update buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.tab === tabName) {
                btn.classList.add('active');
            }
        });

        // Update content
        document.querySelectorAll('.tab-pane').forEach(pane => {
            pane.classList.remove('active');
            if (pane.id === tabName) {
                pane.classList.add('active');
            }
        });

        // Load data for trends tab if needed
        if (tabName === 'trends' && !this.data.trends) {
            this.loadTrends();
        }
    },

    loadURLParams() {
        const params = new URLSearchParams(window.location.search);
        const runId = params.get('run');
        const tab = params.get('tab');

        if (tab) {
            this.switchTab(tab);
        }

        if (runId) {
            setTimeout(() => this.viewRun(runId), 500);
        }
    }
};

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    App.init();
});
`
