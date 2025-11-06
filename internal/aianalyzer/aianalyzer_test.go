package aianalyzer

import (
	"os"
	"strings"
	"testing"

	"github.com/alenon/gokanon/internal/models"
)

func TestNewAnalyzer(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "disabled analyzer",
			config: Config{
				Enabled: false,
			},
			expectError: false,
		},
		{
			name: "unsupported provider",
			config: Config{
				Enabled:  true,
				Provider: "unsupported",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer, err := NewAnalyzer(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if analyzer == nil {
				t.Error("Expected non-nil analyzer")
			}

			if analyzer.config.Enabled != tt.config.Enabled {
				t.Errorf("Expected enabled %v, got %v", tt.config.Enabled, analyzer.config.Enabled)
			}
		})
	}
}

func TestNewFromEnv(t *testing.T) {
	// Save original env vars
	origEnabled := os.Getenv("GOKANON_AI_ENABLED")
	origProvider := os.Getenv("GOKANON_AI_PROVIDER")
	origModel := os.Getenv("GOKANON_AI_MODEL")
	origAPIKey := os.Getenv("GOKANON_AI_API_KEY")
	origBaseURL := os.Getenv("GOKANON_AI_BASE_URL")

	defer func() {
		os.Setenv("GOKANON_AI_ENABLED", origEnabled)
		os.Setenv("GOKANON_AI_PROVIDER", origProvider)
		os.Setenv("GOKANON_AI_MODEL", origModel)
		os.Setenv("GOKANON_AI_API_KEY", origAPIKey)
		os.Setenv("GOKANON_AI_BASE_URL", origBaseURL)
	}()

	tests := []struct {
		name            string
		enabled         string
		provider        string
		model           string
		apiKey          string
		baseURL         string
		expectedEnabled bool
	}{
		{
			name:            "disabled via env",
			enabled:         "false",
			provider:        "ollama",
			expectedEnabled: false,
		},
		{
			name:            "custom provider disabled",
			enabled:         "false",
			provider:        "groq",
			model:           "llama-3.3-70b-versatile",
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("GOKANON_AI_ENABLED", tt.enabled)
			os.Setenv("GOKANON_AI_PROVIDER", tt.provider)
			os.Setenv("GOKANON_AI_MODEL", tt.model)
			os.Setenv("GOKANON_AI_API_KEY", tt.apiKey)
			os.Setenv("GOKANON_AI_BASE_URL", tt.baseURL)

			analyzer, err := NewFromEnv()

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if analyzer != nil && analyzer.config.Enabled != tt.expectedEnabled {
				t.Errorf("Expected enabled %v, got %v", tt.expectedEnabled, analyzer.config.Enabled)
			}
		})
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		setValue     string
		expected     string
	}{
		{
			name:         "env var not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default",
			setValue:     "",
			expected:     "default",
		},
		{
			name:         "env var set",
			key:          "TEST_VAR_SET",
			defaultValue: "default",
			setValue:     "custom",
			expected:     "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvWithDefault(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestEnhanceProfileSummaryDisabled(t *testing.T) {
	analyzer := &Analyzer{
		config: Config{Enabled: false},
	}

	summary := &models.ProfileSummary{
		TotalCPUSamples: 1000,
	}

	result, err := analyzer.EnhanceProfileSummary(summary)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != summary {
		t.Error("Expected same summary when analyzer is disabled")
	}
}

func TestAnalyzeComparisonDisabled(t *testing.T) {
	analyzer := &Analyzer{
		config: Config{Enabled: false},
	}

	oldRun := &models.BenchmarkRun{
		ID: "test-1",
	}
	newRun := &models.BenchmarkRun{
		ID: "test-2",
	}
	comparisons := []models.Comparison{
		{Name: "BenchmarkTest", OldNsPerOp: 100, NewNsPerOp: 200},
	}

	result, err := analyzer.AnalyzeComparison(oldRun, newRun, comparisons)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "" {
		t.Error("Expected empty string when analyzer is disabled")
	}
}

// Mock provider for testing
type mockProvider struct {
	analyzeResult string
	analyzeError  error
}

func (m *mockProvider) Analyze(prompt string) (string, error) {
	return m.analyzeResult, m.analyzeError
}

func TestEnhanceProfileSummaryWithMockProvider(t *testing.T) {
	mockResult := `[
		{
			"type": "cpu",
			"severity": "high",
			"function": "TestFunc",
			"issue": "Test issue",
			"suggestion": "Test suggestion",
			"impact": "50% improvement"
		}
	]`

	mock := &mockProvider{
		analyzeResult: mockResult,
		analyzeError:  nil,
	}

	analyzer := &Analyzer{
		config:   Config{Enabled: true},
		provider: mock,
	}

	summary := &models.ProfileSummary{
		TotalCPUSamples: 1000,
		HotPaths: []models.HotPath{
			{Path: []string{"TestFunc"}, Percentage: 50.0},
		},
	}

	result, err := analyzer.EnhanceProfileSummary(summary)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Suggestions) == 0 {
		t.Error("Expected suggestions to be populated")
	}
}

func TestAnalyzeComparisonWithMockProvider(t *testing.T) {
	mockResult := "Analysis: Performance has regressed significantly."

	mock := &mockProvider{
		analyzeResult: mockResult,
		analyzeError:  nil,
	}

	analyzer := &Analyzer{
		config:   Config{Enabled: true},
		provider: mock,
	}

	oldRun := &models.BenchmarkRun{
		ID: "test-1",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkTest", NsPerOp: 100.0},
		},
	}
	newRun := &models.BenchmarkRun{
		ID: "test-2",
		Results: []models.BenchmarkResult{
			{Name: "BenchmarkTest", NsPerOp: 200.0},
		},
	}
	comparisons := []models.Comparison{
		{
			Name:         "BenchmarkTest",
			OldNsPerOp:   100.0,
			NewNsPerOp:   200.0,
			DeltaPercent: 100.0,
			Status:       "degraded",
		},
	}

	result, err := analyzer.AnalyzeComparison(oldRun, newRun, comparisons)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	if !strings.Contains(result, "regressed") {
		t.Error("Expected analysis to mention regression")
	}
}

func TestPrepareProfileContext(t *testing.T) {
	analyzer := &Analyzer{}

	summary := &models.ProfileSummary{
		TotalCPUSamples: 1000,
		HotPaths: []models.HotPath{
			{Path: []string{"TestFunc"}, Percentage: 50.0},
		},
		MemoryLeaks: []models.MemoryLeak{
			{Function: "LeakFunc", Bytes: 1024, Description: "Test leak"},
		},
	}

	context, err := analyzer.prepareProfileContext(summary)
	if err != nil {
		t.Fatalf("prepareProfileContext failed: %v", err)
	}

	if context == "" {
		t.Error("Expected non-empty context")
	}

	if len(context) < 50 {
		t.Error("Expected context to be substantial")
	}
}

func TestPrepareComparisonContext(t *testing.T) {
	analyzer := &Analyzer{}

	oldRun := &models.BenchmarkRun{
		ID:      "test-1",
		Package: "./test",
	}
	newRun := &models.BenchmarkRun{
		ID:      "test-2",
		Package: "./test",
	}
	comparisons := []models.Comparison{
		{
			Name:         "BenchmarkTest",
			OldNsPerOp:   100.0,
			NewNsPerOp:   200.0,
			DeltaPercent: 100.0,
			Status:       "degraded",
		},
	}

	context := analyzer.prepareComparisonContext(oldRun, newRun, comparisons)

	if context == "" {
		t.Error("Expected non-empty context")
	}

	if len(context) < 50 {
		t.Error("Expected context to be substantial")
	}
}
