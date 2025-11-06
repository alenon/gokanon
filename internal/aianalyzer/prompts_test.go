package aianalyzer

import (
	"strings"
	"testing"
)

func TestBuildProfileAnalysisPrompt(t *testing.T) {
	context := "Test profiling data"
	prompt := buildProfileAnalysisPrompt(context)

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if !strings.Contains(prompt, context) {
		t.Error("Expected prompt to contain context")
	}

	if !strings.Contains(prompt, "JSON") {
		t.Error("Expected prompt to mention JSON format")
	}

	if !strings.Contains(prompt, "optimization") {
		t.Error("Expected prompt to mention optimization")
	}
}

func TestBuildComparisonAnalysisPrompt(t *testing.T) {
	context := "Test comparison data"
	prompt := buildComparisonAnalysisPrompt(context)

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if !strings.Contains(prompt, context) {
		t.Error("Expected prompt to contain context")
	}

	if !strings.Contains(prompt, "comparison") {
		t.Error("Expected prompt to mention comparison")
	}
}

func TestParseTextSuggestions(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedCount int
	}{
		{
			name: "JSON in markdown code block",
			text: "```json\n" +
				`[{"type":"cpu","severity":"high","function":"TestFunc","issue":"High CPU","suggestion":"Optimize","impact":"50%"}]` +
				"\n```",
			expectedCount: 1,
		},
		{
			name: "plain JSON array",
			text: `[{"type":"memory","severity":"medium","function":"MemFunc","issue":"Memory leak","suggestion":"Fix leak","impact":"20%"}]`,
			expectedCount: 1,
		},
		{
			name: "bullet points",
			text: "- Optimize TestFunc for better performance\n" +
				"- Reduce memory allocations in ProcessData\n" +
				"- Use sync.Pool for temporary buffers",
			expectedCount: 3,
		},
		{
			name:          "empty text",
			text:          "",
			expectedCount: 0,
		},
		{
			name:          "no suggestions",
			text:          "This is just some text without any suggestions.",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := parseTextSuggestions(tt.text)

			if len(suggestions) != tt.expectedCount {
				t.Errorf("Expected %d suggestions, got %d", tt.expectedCount, len(suggestions))
			}
		})
	}
}

func TestExtractBulletSuggestions(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedCount int
	}{
		{
			name: "dash bullets",
			text: "- First suggestion\n" +
				"- Second suggestion\n" +
				"- Third suggestion",
			expectedCount: 3,
		},
		{
			name: "asterisk bullets",
			text: "* First suggestion\n" +
				"* Second suggestion",
			expectedCount: 2,
		},
		{
			name: "numbered list",
			text: "1. First suggestion\n" +
				"2. Second suggestion\n" +
				"3. Third suggestion",
			expectedCount: 3,
		},
		{
			name: "mixed with empty lines",
			text: "- First suggestion\n" +
				"\n" +
				"- Second suggestion\n" +
				"\n",
			expectedCount: 2,
		},
		{
			name: "multiline suggestions",
			text: "- First suggestion\n" +
				"  continuation of first\n" +
				"- Second suggestion",
			expectedCount: 2,
		},
		{
			name:          "no bullets",
			text:          "Just plain text",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := extractBulletSuggestions(tt.text)

			if len(suggestions) != tt.expectedCount {
				t.Errorf("Expected %d suggestions, got %d", tt.expectedCount, len(suggestions))
			}

			// Verify suggestions have required fields
			for i, sug := range suggestions {
				if sug.Type == "" {
					t.Errorf("Suggestion %d missing type", i)
				}
				if sug.Severity == "" {
					t.Errorf("Suggestion %d missing severity", i)
				}
				if sug.Suggestion == "" {
					t.Errorf("Suggestion %d missing suggestion text", i)
				}
			}
		})
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		text         string
		expectedType string
	}{
		{"High CPU usage in function", "cpu"},
		{"Slow performance detected", "cpu"},
		{"Memory leak in handler", "memory"},
		{"Excessive allocations found", "memory"},
		{"Algorithm complexity issue", "algorithm"},
		{"General optimization needed", "general"},
		{"Time consuming operation", "cpu"},
		{"GC pressure detected", "memory"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := inferType(tt.text)
			if result != tt.expectedType {
				t.Errorf("For text %q, expected type %q, got %q", tt.text, tt.expectedType, result)
			}
		})
	}
}

func TestInferSeverity(t *testing.T) {
	tests := []struct {
		text             string
		expectedSeverity string
	}{
		{"Critical performance issue", "high"},
		{"Severe memory leak", "high"},
		{"Major bottleneck", "high"},
		{"Moderate optimization needed", "medium"},
		{"Significant improvement possible", "medium"},
		{"Minor tweak suggested", "low"},
		{"Small optimization", "low"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := inferSeverity(tt.text)
			if result != tt.expectedSeverity {
				t.Errorf("For text %q, expected severity %q, got %q", tt.text, tt.expectedSeverity, result)
			}
		})
	}
}

func TestExtractFunction(t *testing.T) {
	tests := []struct {
		text             string
		expectedFunction string
	}{
		{"Optimize ProcessData() for better performance", "ProcessData"},
		{"The `HandleRequest` function is slow", "HandleRequest"},
		{"Issue in database.Query() method", "database.Query"},
		{"General performance issue", "General"},
		{"No function mentioned here", "General"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := extractFunction(tt.text)
			if result != tt.expectedFunction {
				t.Errorf("For text %q, expected function %q, got %q", tt.text, tt.expectedFunction, result)
			}
		})
	}
}

func TestExtractIssue(t *testing.T) {
	tests := []struct {
		text          string
		expectedIssue string
	}{
		{
			text:          "High CPU usage: ProcessData is consuming 80% CPU",
			expectedIssue: "High CPU usage",
		},
		{
			text:          "Memory leak - allocations are growing",
			expectedIssue: "Memory leak",
		},
		{
			text:          "Performance degradation",
			expectedIssue: "Performance degradation",
		},
		{
			text:          strings.Repeat("a", 150),
			expectedIssue: strings.Repeat("a", 97) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedIssue, func(t *testing.T) {
			result := extractIssue(tt.text)
			if result != tt.expectedIssue {
				t.Errorf("Expected issue %q, got %q", tt.expectedIssue, result)
			}
		})
	}
}

func TestParseTextSuggestionsWithComplexJSON(t *testing.T) {
	text := `Here are my suggestions:

` + "```json" + `
[
  {
    "type": "cpu",
    "severity": "high",
    "function": "ProcessData",
    "issue": "High CPU usage",
    "suggestion": "Use more efficient algorithm",
    "impact": "50% improvement"
  },
  {
    "type": "memory",
    "severity": "medium",
    "function": "HandleRequest",
    "issue": "Memory leak",
    "suggestion": "Fix resource cleanup",
    "impact": "20% reduction"
  }
]
` + "```"

	suggestions := parseTextSuggestions(text)

	if len(suggestions) != 2 {
		t.Fatalf("Expected 2 suggestions, got %d", len(suggestions))
	}

	// Verify first suggestion
	if suggestions[0].Type != "cpu" {
		t.Errorf("Expected type 'cpu', got %q", suggestions[0].Type)
	}
	if suggestions[0].Severity != "high" {
		t.Errorf("Expected severity 'high', got %q", suggestions[0].Severity)
	}
	if suggestions[0].Function != "ProcessData" {
		t.Errorf("Expected function 'ProcessData', got %q", suggestions[0].Function)
	}

	// Verify second suggestion
	if suggestions[1].Type != "memory" {
		t.Errorf("Expected type 'memory', got %q", suggestions[1].Type)
	}
	if suggestions[1].Severity != "medium" {
		t.Errorf("Expected severity 'medium', got %q", suggestions[1].Severity)
	}
}

func TestExtractBulletSuggestionsWithContinuation(t *testing.T) {
	text := `- Optimize the ProcessData function
  which is currently consuming too much CPU
  and causing performance issues
- Reduce memory allocations in the handler`

	suggestions := extractBulletSuggestions(text)

	if len(suggestions) != 2 {
		t.Fatalf("Expected 2 suggestions, got %d", len(suggestions))
	}

	// Verify first suggestion includes continuation
	if !strings.Contains(suggestions[0].Suggestion, "consuming too much CPU") {
		t.Error("Expected first suggestion to include continuation text")
	}
	if !strings.Contains(suggestions[0].Suggestion, "performance issues") {
		t.Error("Expected first suggestion to include all continuation text")
	}
}
