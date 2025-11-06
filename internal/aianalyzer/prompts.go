package aianalyzer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/alenon/gokanon/internal/models"
)

// buildProfileAnalysisPrompt creates a prompt for profile analysis
func buildProfileAnalysisPrompt(context string) string {
	return fmt.Sprintf(`You are analyzing Go benchmark profiling data. Based on the data below, provide optimization suggestions.

PROFILING DATA:
%s

Please analyze this data and provide actionable optimization suggestions. For each suggestion, provide:
1. Type: "cpu", "memory", or "algorithm"
2. Severity: "low", "medium", or "high"
3. Function: The affected function name
4. Issue: A brief description of the problem
5. Suggestion: Specific, actionable advice on how to fix it
6. Impact: Expected performance improvement

Respond with a JSON array of suggestions in this format:
[
  {
    "type": "cpu",
    "severity": "high",
    "function": "functionName",
    "issue": "Brief issue description",
    "suggestion": "Specific actionable advice",
    "impact": "Expected improvement description"
  }
]

Focus on the most impactful optimizations. Consider:
- Hot functions consuming significant CPU/memory
- Potential memory leaks (high allocation with low in-use memory)
- Hot paths that could be optimized
- Common Go performance patterns (e.g., unnecessary allocations, inefficient algorithms)
- Opportunities for sync.Pool, buffering, or pre-allocation

Be specific and actionable in your suggestions.`, context)
}

// buildComparisonAnalysisPrompt creates a prompt for comparison analysis
func buildComparisonAnalysisPrompt(context string) string {
	return fmt.Sprintf(`You are analyzing benchmark comparison results between two Go benchmark runs.

COMPARISON DATA:
%s

Please analyze the performance changes and provide insights about:
1. Significant improvements or regressions
2. Possible causes for the changes
3. Whether the changes are concerning or expected
4. Recommendations for next steps

Provide a concise analysis (2-3 paragraphs) focusing on the most important findings.`, context)
}

// parseTextSuggestions attempts to parse suggestions from markdown/text format
func parseTextSuggestions(text string) []models.Suggestion {
	var suggestions []models.Suggestion

	// Try to extract JSON blocks from markdown
	jsonBlockRegex := regexp.MustCompile("```(?:json)?\n(\\[\\s*\\{[\\s\\S]*?\\}\\s*\\])\\n```")
	matches := jsonBlockRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		var parsed []models.Suggestion
		if err := json.Unmarshal([]byte(matches[1]), &parsed); err == nil {
			return parsed
		}
	}

	// Try to find JSON array directly in text
	jsonArrayRegex := regexp.MustCompile(`\[\s*\{[^\]]*"type"\s*:[^\]]*"suggestion"\s*:[^\]]*\}\s*\]`)
	match := jsonArrayRegex.FindString(text)
	if match != "" {
		var parsed []models.Suggestion
		if err := json.Unmarshal([]byte(match), &parsed); err == nil {
			return parsed
		}
	}

	// If no structured format found, try to extract bullet points
	suggestions = extractBulletSuggestions(text)
	if len(suggestions) > 0 {
		return suggestions
	}

	return nil
}

// extractBulletSuggestions extracts suggestions from bullet-point format
func extractBulletSuggestions(text string) []models.Suggestion {
	var suggestions []models.Suggestion

	// Common bullet point patterns
	lines := strings.Split(text, "\n")
	var currentSuggestion *models.Suggestion

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for bullet points or numbered lists
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") ||
		   regexp.MustCompile(`^\d+\.`).MatchString(line) {
			// Save previous suggestion if exists
			if currentSuggestion != nil {
				suggestions = append(suggestions, *currentSuggestion)
			}

			// Start new suggestion
			cleaned := regexp.MustCompile(`^[-*\d.]\s*`).ReplaceAllString(line, "")
			currentSuggestion = &models.Suggestion{
				Type:       inferType(cleaned),
				Severity:   inferSeverity(cleaned),
				Function:   extractFunction(cleaned),
				Issue:      extractIssue(cleaned),
				Suggestion: cleaned,
				Impact:     "Potential performance improvement",
			}
		} else if currentSuggestion != nil {
			// Continuation of current suggestion
			currentSuggestion.Suggestion += " " + line
		}
	}

	// Add last suggestion
	if currentSuggestion != nil {
		suggestions = append(suggestions, *currentSuggestion)
	}

	return suggestions
}

// inferType tries to infer the suggestion type from text
func inferType(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "cpu") || strings.Contains(lower, "time") ||
	   strings.Contains(lower, "slow") || strings.Contains(lower, "performance") {
		return "cpu"
	}
	if strings.Contains(lower, "memory") || strings.Contains(lower, "allocation") ||
	   strings.Contains(lower, "leak") || strings.Contains(lower, "gc") {
		return "memory"
	}
	if strings.Contains(lower, "algorithm") || strings.Contains(lower, "complexity") {
		return "algorithm"
	}
	return "general"
}

// inferSeverity tries to infer severity from text
func inferSeverity(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "critical") || strings.Contains(lower, "severe") ||
	   strings.Contains(lower, "major") {
		return "high"
	}
	if strings.Contains(lower, "moderate") || strings.Contains(lower, "significant") {
		return "medium"
	}
	return "low"
}

// extractFunction tries to extract function name from text
func extractFunction(text string) string {
	// Look for common function patterns
	funcRegex := regexp.MustCompile(`(\w+(?:\.\w+)*)\s*\(`)
	if match := funcRegex.FindStringSubmatch(text); len(match) > 1 {
		return match[1]
	}

	// Look for backtick-quoted names
	backtickRegex := regexp.MustCompile("`([^`]+)`")
	if match := backtickRegex.FindStringSubmatch(text); len(match) > 1 {
		return match[1]
	}

	return "General"
}

// extractIssue tries to extract the issue description
func extractIssue(text string) string {
	// Split on common delimiters
	parts := regexp.MustCompile(`[:|-]`).Split(text, 2)
	if len(parts) > 0 {
		issue := strings.TrimSpace(parts[0])
		if len(issue) > 100 {
			return issue[:97] + "..."
		}
		return issue
	}
	return "Performance issue detected"
}
