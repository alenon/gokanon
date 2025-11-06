package ui

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrorWithSuggestion_Error(t *testing.T) {
	tests := []struct {
		name        string
		err         *ErrorWithSuggestion
		wantMessage string
		wantCause   string
		wantIcon    string
	}{
		{
			name: "error with suggestions",
			err: &ErrorWithSuggestion{
				Message:     "Test error",
				Suggestions: []string{"Try this", "Or this"},
				Err:         nil,
			},
			wantMessage: "Test error",
			wantIcon:    "✗",
		},
		{
			name: "error with cause",
			err: &ErrorWithSuggestion{
				Message:     "Operation failed",
				Suggestions: []string{},
				Err:         errors.New("underlying cause"),
			},
			wantMessage: "Operation failed",
			wantCause:   "underlying cause",
		},
		{
			name: "error with everything",
			err: &ErrorWithSuggestion{
				Message:     "Complete error",
				Suggestions: []string{"Suggestion 1", "Suggestion 2"},
				Err:         errors.New("root cause"),
			},
			wantMessage: "Complete error",
			wantCause:   "root cause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()

			if !strings.Contains(result, tt.wantMessage) {
				t.Errorf("Error() = %q, want to contain %q", result, tt.wantMessage)
			}

			if tt.wantCause != "" && !strings.Contains(result, tt.wantCause) {
				t.Errorf("Error() = %q, want to contain cause %q", result, tt.wantCause)
			}

			if tt.wantIcon != "" && !strings.Contains(result, tt.wantIcon) {
				t.Errorf("Error() = %q, want to contain icon %q", result, tt.wantIcon)
			}

			for _, suggestion := range tt.err.Suggestions {
				if !strings.Contains(result, suggestion) {
					t.Errorf("Error() = %q, want to contain suggestion %q", result, suggestion)
				}
			}
		})
	}
}

func TestErrorWithSuggestion_Unwrap(t *testing.T) {
	rootErr := errors.New("root cause")
	err := &ErrorWithSuggestion{
		Message:     "Wrapped error",
		Suggestions: []string{},
		Err:         rootErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != rootErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, rootErr)
	}
}

func TestErrorWithSuggestion_UnwrapNil(t *testing.T) {
	err := &ErrorWithSuggestion{
		Message:     "Error without cause",
		Suggestions: []string{},
		Err:         nil,
	}

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Errorf("Unwrap() = %v, want nil", unwrapped)
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		err         error
		suggestions []string
	}{
		{
			name:        "basic error",
			message:     "Something went wrong",
			err:         nil,
			suggestions: []string{"Try again"},
		},
		{
			name:        "error with cause",
			message:     "Failed to process",
			err:         errors.New("file not found"),
			suggestions: []string{"Check the file path", "Ensure file exists"},
		},
		{
			name:        "error with multiple suggestions",
			message:     "Configuration error",
			err:         nil,
			suggestions: []string{"Check config file", "Verify syntax", "See documentation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewError(tt.message, tt.err, tt.suggestions...)

			if err == nil {
				t.Fatal("NewError() returned nil")
			}

			if err.Message != tt.message {
				t.Errorf("Message = %q, want %q", err.Message, tt.message)
			}

			if err.Err != tt.err {
				t.Errorf("Err = %v, want %v", err.Err, tt.err)
			}

			if len(err.Suggestions) != len(tt.suggestions) {
				t.Errorf("Got %d suggestions, want %d", len(err.Suggestions), len(tt.suggestions))
			}

			for i, suggestion := range tt.suggestions {
				if err.Suggestions[i] != suggestion {
					t.Errorf("Suggestions[%d] = %q, want %q", i, err.Suggestions[i], suggestion)
				}
			}
		})
	}
}

func TestErrNoResults(t *testing.T) {
	err := ErrNoResults()

	if err == nil {
		t.Fatal("ErrNoResults() returned nil")
	}

	errStr := err.Error()

	// Should contain key information
	expectedStrings := []string{
		"No benchmark results found",
		"gokanon run",
		".gokanon",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error string should contain %q, got %q", expected, errStr)
		}
	}
}

func TestErrInvalidRunID(t *testing.T) {
	availableIDs := []string{"run-1", "run-2", "run-3"}
	err := ErrInvalidRunID("invalid-run", availableIDs)

	if err == nil {
		t.Fatal("ErrInvalidRunID() returned nil")
	}

	errStr := err.Error()

	if !strings.Contains(errStr, "invalid-run") {
		t.Errorf("Error should contain the invalid ID, got %q", errStr)
	}

	if !strings.Contains(errStr, "gokanon list") {
		t.Errorf("Error should suggest gokanon list, got %q", errStr)
	}

	// Should show available IDs
	for _, id := range availableIDs {
		if !strings.Contains(errStr, id) {
			t.Errorf("Error should contain available ID %q, got %q", id, errStr)
		}
	}
}

func TestErrInvalidRunID_EmptyList(t *testing.T) {
	err := ErrInvalidRunID("invalid-run", []string{})

	if err == nil {
		t.Fatal("ErrInvalidRunID() returned nil")
	}

	errStr := err.Error()

	if !strings.Contains(errStr, "invalid-run") {
		t.Errorf("Error should contain the invalid ID, got %q", errStr)
	}
}

func TestErrBenchmarkFailed(t *testing.T) {
	rootErr := errors.New("compilation failed")
	err := ErrBenchmarkFailed(rootErr)

	if err == nil {
		t.Fatal("ErrBenchmarkFailed() returned nil")
	}

	errStr := err.Error()

	expectedStrings := []string{
		"Benchmark execution failed",
		"BenchmarkXxx",
		"go test",
		"go version",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error should contain %q, got %q", expected, errStr)
		}
	}

	// Should be able to unwrap to root cause
	wrapped := err.(*ErrorWithSuggestion)
	if wrapped.Unwrap() != rootErr {
		t.Errorf("Unwrap() should return root error")
	}
}

func TestErrInvalidThreshold(t *testing.T) {
	err := ErrInvalidThreshold("invalid")

	if err == nil {
		t.Fatal("ErrInvalidThreshold() returned nil")
	}

	errStr := err.Error()

	expectedStrings := []string{
		"Invalid threshold value",
		"invalid",
		"positive number",
		"10",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error should contain %q, got %q", expected, errStr)
		}
	}
}

func TestErrStorageCorrupted(t *testing.T) {
	rootErr := errors.New("json parse error")
	err := ErrStorageCorrupted(rootErr)

	if err == nil {
		t.Fatal("ErrStorageCorrupted() returned nil")
	}

	errStr := err.Error()

	expectedStrings := []string{
		"storage",
		"corrupted",
		"gokanon doctor",
		".gokanon",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error should contain %q, got %q", expected, errStr)
		}
	}
}

func TestErrProfileNotFound(t *testing.T) {
	err := ErrProfileNotFound("run-123")

	if err == nil {
		t.Fatal("ErrProfileNotFound() returned nil")
	}

	errStr := err.Error()

	expectedStrings := []string{
		"Profile data not found",
		"run-123",
		"-profile",
		"cpu,mem",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error should contain %q, got %q", expected, errStr)
		}
	}
}

func TestErrInvalidFormat(t *testing.T) {
	err := ErrInvalidFormat("xml")

	if err == nil {
		t.Fatal("ErrInvalidFormat() returned nil")
	}

	errStr := err.Error()

	expectedStrings := []string{
		"Unsupported export format",
		"xml",
		"html",
		"csv",
		"markdown",
		"json",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error should contain %q, got %q", expected, errStr)
		}
	}
}

func TestErrPortInUse(t *testing.T) {
	rootErr := errors.New("bind: address already in use")
	err := ErrPortInUse(8080, rootErr)

	if err == nil {
		t.Fatal("ErrPortInUse() returned nil")
	}

	errStr := err.Error()

	expectedStrings := []string{
		"Port 8080",
		"already in use",
		"-port=",
		"lsof",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(errStr, expected) {
			t.Errorf("Error should contain %q, got %q", expected, errStr)
		}
	}
}

func TestPrintErrorAndExit_NilError(t *testing.T) {
	// Should not panic with nil error
	// We can't actually test the exit, but we can ensure it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintErrorAndExit panicked with nil error: %v", r)
		}
	}()

	// This won't actually exit in test, but shouldn't panic
	go func() {
		PrintErrorAndExit(nil, 1)
	}()
}

func TestErrorWithSuggestion_Format(t *testing.T) {
	err := NewError(
		"Test error message",
		errors.New("underlying error"),
		"First suggestion",
		"Second suggestion",
		"Third suggestion",
	)

	errStr := err.Error()

	// Should contain all parts
	if !strings.Contains(errStr, "Test error message") {
		t.Error("Should contain main message")
	}
	if !strings.Contains(errStr, "underlying error") {
		t.Error("Should contain underlying error")
	}
	if !strings.Contains(errStr, "First suggestion") {
		t.Error("Should contain first suggestion")
	}
	if !strings.Contains(errStr, "Second suggestion") {
		t.Error("Should contain second suggestion")
	}
	if !strings.Contains(errStr, "Third suggestion") {
		t.Error("Should contain third suggestion")
	}
	if !strings.Contains(errStr, "💡") {
		t.Error("Should contain suggestion icon")
	}
}

func TestErrorWithSuggestion_NoSuggestions(t *testing.T) {
	err := NewError("Just an error", nil)

	errStr := err.Error()

	if !strings.Contains(errStr, "Just an error") {
		t.Error("Should contain error message")
	}

	// Should not contain suggestions section
	if strings.Contains(errStr, "Suggestions:") {
		t.Error("Should not contain suggestions section when there are no suggestions")
	}
}

func TestErrorImplementsError(t *testing.T) {
	var _ error = &ErrorWithSuggestion{}
}

func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewError("test error", nil, "suggestion 1", "suggestion 2")
	}
}

func BenchmarkErrorWithSuggestion_Error(b *testing.B) {
	err := NewError("test error", errors.New("cause"), "suggestion 1", "suggestion 2")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func ExampleNewError() {
	err := NewError(
		"Failed to connect to database",
		errors.New("connection timeout"),
		"Check your network connection",
		"Verify database server is running",
		"Check firewall settings",
	)

	fmt.Println(err)
}

func ExampleErrNoResults() {
	err := ErrNoResults()
	fmt.Println(err)
}

// Test that all error constructors return non-nil errors
func TestAllErrorConstructorsReturnNonNil(t *testing.T) {
	constructors := []struct {
		name string
		err  error
	}{
		{"ErrNoResults", ErrNoResults()},
		{"ErrInvalidRunID", ErrInvalidRunID("test", []string{})},
		{"ErrBenchmarkFailed", ErrBenchmarkFailed(errors.New("test"))},
		{"ErrInvalidThreshold", ErrInvalidThreshold("test")},
		{"ErrStorageCorrupted", ErrStorageCorrupted(errors.New("test"))},
		{"ErrProfileNotFound", ErrProfileNotFound("test")},
		{"ErrInvalidFormat", ErrInvalidFormat("test")},
		{"ErrPortInUse", ErrPortInUse(8080, errors.New("test"))},
	}

	for _, tc := range constructors {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Errorf("%s returned nil", tc.name)
			}
		})
	}
}

// Test error chaining with errors.Is and errors.As
func TestErrorChaining(t *testing.T) {
	rootErr := errors.New("root cause")
	err := NewError("wrapper", rootErr, "suggestion")

	// Should be able to unwrap
	if !errors.Is(err, rootErr) {
		t.Error("errors.Is should find root cause")
	}

	// Should be able to cast
	var ewsErr *ErrorWithSuggestion
	if !errors.As(err, &ewsErr) {
		t.Error("errors.As should work with ErrorWithSuggestion")
	}
}
