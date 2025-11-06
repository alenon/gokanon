package aianalyzer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alenon/gokanon/internal/models"
)

// Config holds AI analyzer configuration
type Config struct {
	Enabled  bool
	Provider string // "ollama" or "groq"
	Model    string // Model name to use
	APIKey   string // API key for cloud providers (not needed for Ollama)
	BaseURL  string // Base URL for the provider
}

// Analyzer provides AI-powered analysis of benchmark results
type Analyzer struct {
	config   Config
	provider AIProvider
}

// NewAnalyzer creates a new AI analyzer
func NewAnalyzer(config Config) (*Analyzer, error) {
	if !config.Enabled {
		return &Analyzer{config: config}, nil
	}

	var provider AIProvider
	var err error

	switch config.Provider {
	case "ollama":
		provider, err = NewOllamaProvider(config)
	case "groq":
		provider, err = NewGroqProvider(config)
	case "openai":
		provider, err = NewOpenAIProvider(config)
	case "anthropic", "claude":
		provider, err = NewAnthropicProvider(config)
	case "gemini":
		provider, err = NewGeminiProvider(config)
	case "openai-compatible", "custom":
		provider, err = NewOpenAICompatibleProvider(config)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s (supported: ollama, groq, openai, anthropic, gemini, openai-compatible)", config.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	return &Analyzer{
		config:   config,
		provider: provider,
	}, nil
}

// NewFromEnv creates an analyzer from environment variables
func NewFromEnv() (*Analyzer, error) {
	config := Config{
		Enabled:  os.Getenv("GOKANON_AI_ENABLED") == "true",
		Provider: getEnvWithDefault("GOKANON_AI_PROVIDER", "ollama"),
		Model:    getEnvWithDefault("GOKANON_AI_MODEL", ""),
		APIKey:   os.Getenv("GOKANON_AI_API_KEY"),
		BaseURL:  getEnvWithDefault("GOKANON_AI_BASE_URL", ""),
	}

	// Set default models if not specified
	if config.Model == "" {
		switch config.Provider {
		case "ollama":
			config.Model = "llama3.2"
		case "groq":
			config.Model = "llama-3.3-70b-versatile"
		case "openai":
			config.Model = "gpt-4o"
		case "anthropic", "claude":
			config.Model = "claude-sonnet-4-5-20250929"
		case "gemini":
			config.Model = "gemini-2.5-flash"
		case "openai-compatible", "custom":
			config.Model = "default" // Let the service use its default model
		}
	}

	// Set default base URLs if not specified
	if config.BaseURL == "" {
		switch config.Provider {
		case "ollama":
			config.BaseURL = "http://localhost:11434"
		case "groq":
			config.BaseURL = "https://api.groq.com/openai/v1"
		case "openai":
			config.BaseURL = "https://api.openai.com"
		case "anthropic", "claude":
			config.BaseURL = "https://api.anthropic.com"
		case "gemini":
			config.BaseURL = "https://generativelanguage.googleapis.com"
		case "openai-compatible", "custom":
			config.BaseURL = "http://localhost:8080" // Placeholder, should be set by user
		}
	}

	return NewAnalyzer(config)
}

// EnhanceProfileSummary enhances a profile summary with AI insights
func (a *Analyzer) EnhanceProfileSummary(summary *models.ProfileSummary) (*models.ProfileSummary, error) {
	if !a.config.Enabled || a.provider == nil {
		return summary, nil
	}

	// Prepare context for AI
	context, err := a.prepareProfileContext(summary)
	if err != nil {
		return summary, fmt.Errorf("failed to prepare context: %w", err)
	}

	// Get AI analysis
	prompt := buildProfileAnalysisPrompt(context)
	response, err := a.provider.Analyze(prompt)
	if err != nil {
		return summary, fmt.Errorf("AI analysis failed: %w", err)
	}

	// Parse AI response and enhance suggestions
	enhancedSuggestions, err := a.parseAISuggestions(response, summary)
	if err != nil {
		// If parsing fails, keep original suggestions
		return summary, fmt.Errorf("failed to parse AI suggestions: %w", err)
	}

	// Create enhanced summary
	enhanced := *summary
	enhanced.Suggestions = enhancedSuggestions
	return &enhanced, nil
}

// AnalyzeComparison provides AI insights on benchmark comparison
func (a *Analyzer) AnalyzeComparison(oldRun, newRun *models.BenchmarkRun, comparisons []models.Comparison) (string, error) {
	if !a.config.Enabled || a.provider == nil {
		return "", nil
	}

	// Prepare comparison context
	context := a.prepareComparisonContext(oldRun, newRun, comparisons)

	// Get AI analysis
	prompt := buildComparisonAnalysisPrompt(context)
	response, err := a.provider.Analyze(prompt)
	if err != nil {
		return "", fmt.Errorf("AI comparison analysis failed: %w", err)
	}

	return response, nil
}

// prepareProfileContext converts profile summary to AI-friendly format
func (a *Analyzer) prepareProfileContext(summary *models.ProfileSummary) (string, error) {
	context := map[string]interface{}{
		"cpu_top_functions":    summary.CPUTopFunctions,
		"memory_top_functions": summary.MemoryTopFunctions,
		"memory_leaks":         summary.MemoryLeaks,
		"hot_paths":            summary.HotPaths,
		"total_cpu_samples":    summary.TotalCPUSamples,
		"total_memory_bytes":   summary.TotalMemoryBytes,
		"existing_suggestions": summary.Suggestions,
	}

	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// prepareComparisonContext converts comparison data to AI-friendly format
func (a *Analyzer) prepareComparisonContext(oldRun, newRun *models.BenchmarkRun, comparisons []models.Comparison) string {
	context := map[string]interface{}{
		"old_run": map[string]interface{}{
			"timestamp":  oldRun.Timestamp,
			"go_version": oldRun.GoVersion,
			"package":    oldRun.Package,
		},
		"new_run": map[string]interface{}{
			"timestamp":  newRun.Timestamp,
			"go_version": newRun.GoVersion,
			"package":    newRun.Package,
		},
		"comparisons": comparisons,
	}

	data, _ := json.MarshalIndent(context, "", "  ")
	return string(data)
}

// parseAISuggestions parses AI response and merges with existing suggestions
func (a *Analyzer) parseAISuggestions(aiResponse string, summary *models.ProfileSummary) ([]models.Suggestion, error) {
	// Try to parse as JSON array first
	var aiSuggestions []models.Suggestion
	if err := json.Unmarshal([]byte(aiResponse), &aiSuggestions); err == nil {
		// Merge AI suggestions with existing ones
		return a.mergeSuggestions(summary.Suggestions, aiSuggestions), nil
	}

	// If not JSON, try to extract structured suggestions from markdown/text
	parsed := parseTextSuggestions(aiResponse)
	if len(parsed) > 0 {
		return a.mergeSuggestions(summary.Suggestions, parsed), nil
	}

	// If we can't parse, add the raw AI response as a general suggestion
	if aiResponse != "" {
		general := models.Suggestion{
			Type:       "general",
			Severity:   "info",
			Function:   "Overall Analysis",
			Issue:      "AI Analysis Results",
			Suggestion: aiResponse,
			Impact:     "See detailed analysis above",
		}
		return append(summary.Suggestions, general), nil
	}

	return summary.Suggestions, nil
}

// mergeSuggestions combines original and AI suggestions, removing duplicates
func (a *Analyzer) mergeSuggestions(original, ai []models.Suggestion) []models.Suggestion {
	// Use a map to track suggestions by function+type to avoid duplicates
	seen := make(map[string]bool)
	var merged []models.Suggestion

	// Add original suggestions first
	for _, s := range original {
		key := fmt.Sprintf("%s:%s", s.Function, s.Type)
		if !seen[key] {
			merged = append(merged, s)
			seen[key] = true
		}
	}

	// Add AI suggestions, skipping duplicates
	for _, s := range ai {
		key := fmt.Sprintf("%s:%s", s.Function, s.Type)
		if !seen[key] {
			// Mark AI-enhanced suggestions
			if s.Impact == "" {
				s.Impact = "AI-suggested optimization"
			}
			merged = append(merged, s)
			seen[key] = true
		}
	}

	return merged
}

// getEnvWithDefault gets environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
