package aianalyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AIProvider is the interface for AI service providers
type AIProvider interface {
	Analyze(prompt string) (string, error)
}

// OllamaProvider implements AIProvider for Ollama
type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(config Config) (*OllamaProvider, error) {
	return &OllamaProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze sends a prompt to Ollama and returns the response
func (p *OllamaProvider) Analyze(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model":  p.model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.7,
			"top_p":       0.9,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", p.baseURL)
	resp, err := p.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w (is Ollama running? try: ollama serve)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return response.Response, nil
}

// GroqProvider implements AIProvider for Groq
type GroqProvider struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

// NewGroqProvider creates a new Groq provider
func NewGroqProvider(config Config) (*GroqProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Groq API key is required (set GOKANON_AI_API_KEY environment variable)")
	}

	return &GroqProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		apiKey:  config.APIKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze sends a prompt to Groq and returns the response
func (p *GroqProvider) Analyze(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert Go performance analyst. Provide concise, actionable insights about benchmark results.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", p.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Groq: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Groq API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode Groq response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from Groq")
	}

	return response.Choices[0].Message.Content, nil
}

// OpenAIProvider implements AIProvider for OpenAI (GPT-4o, GPT-4-turbo, etc.)
type OpenAIProvider struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config Config) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required (set GOKANON_AI_API_KEY environment variable)")
	}

	return &OpenAIProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		apiKey:  config.APIKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze sends a prompt to OpenAI and returns the response
func (p *OpenAIProvider) Analyze(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert Go performance analyst. Provide concise, actionable insights about benchmark results.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", p.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode OpenAI response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

// AnthropicProvider implements AIProvider for Anthropic Claude (Sonnet 4.5, Haiku 4.5, etc.)
type AnthropicProvider struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(config Config) (*AnthropicProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required (set GOKANON_AI_API_KEY environment variable)")
	}

	return &AnthropicProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		apiKey:  config.APIKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze sends a prompt to Anthropic and returns the response
func (p *AnthropicProvider) Analyze(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model":      p.model,
		"max_tokens": 2000,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"system": "You are an expert Go performance analyst. Provide concise, actionable insights about benchmark results.",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/messages", p.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Anthropic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Anthropic API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode Anthropic response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	// Concatenate all text blocks
	var result string
	for _, content := range response.Content {
		if content.Type == "text" {
			result += content.Text
		}
	}

	if result == "" {
		return "", fmt.Errorf("no text content in Anthropic response")
	}

	return result, nil
}

// OpenAICompatibleProvider implements AIProvider for OpenAI-compatible endpoints
// This works with any service that implements the OpenAI Chat Completions API format,
// including Cursor, LM Studio, LocalAI, and many other services
type OpenAICompatibleProvider struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

// NewOpenAICompatibleProvider creates a new OpenAI-compatible provider
func NewOpenAICompatibleProvider(config Config) (*OpenAICompatibleProvider, error) {
	// API key is optional for some local services
	return &OpenAICompatibleProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		apiKey:  config.APIKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze sends a prompt to an OpenAI-compatible endpoint and returns the response
func (p *OpenAICompatibleProvider) Analyze(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert Go performance analyst. Provide concise, actionable insights about benchmark results.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Try both with and without /v1 prefix
	url := fmt.Sprintf("%s/chat/completions", p.baseURL)
	if !bytes.Contains([]byte(p.baseURL), []byte("/v1")) {
		url = fmt.Sprintf("%s/v1/chat/completions", p.baseURL)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenAI-compatible endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI-compatible API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI-compatible endpoint")
	}

	return response.Choices[0].Message.Content, nil
}

// GeminiProvider implements AIProvider for Google Gemini (Gemini 2.5 Flash, 2.0 Flash, etc.)
type GeminiProvider struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(config Config) (*GeminiProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Gemini API key is required (set GOKANON_AI_API_KEY environment variable)")
	}

	return &GeminiProvider{
		baseURL: config.BaseURL,
		model:   config.Model,
		apiKey:  config.APIKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Analyze sends a prompt to Gemini and returns the response
func (p *GeminiProvider) Analyze(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": fmt.Sprintf("You are an expert Go performance analyst. Provide concise, actionable insights about benchmark results.\n\n%s", prompt),
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,
			"maxOutputTokens": 2000,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Gemini uses model in the URL path
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", p.baseURL, p.model, p.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode Gemini response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Concatenate all text parts
	var result string
	for _, part := range response.Candidates[0].Content.Parts {
		result += part.Text
	}

	if result == "" {
		return "", fmt.Errorf("no text content in Gemini response")
	}

	return result, nil
}
