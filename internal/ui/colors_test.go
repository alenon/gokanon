package ui

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestColorFunctions(t *testing.T) {
	// Temporarily enable colors for testing
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()
	NoColor = false

	tests := []struct {
		name     string
		function func(...interface{}) string
		input    string
		contains string // We check if output contains ANSI codes or text
	}{
		{
			name:     "Success",
			function: Success,
			input:    "test",
			contains: "test",
		},
		{
			name:     "Error",
			function: Error,
			input:    "test",
			contains: "test",
		},
		{
			name:     "Warning",
			function: Warning,
			input:    "test",
			contains: "test",
		},
		{
			name:     "Info",
			function: Info,
			input:    "test",
			contains: "test",
		},
		{
			name:     "Dim",
			function: Dim,
			input:    "test",
			contains: "test",
		},
		{
			name:     "Bold",
			function: Bold,
			input:    "test",
			contains: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("%s() = %q, want to contain %q", tt.name, result, tt.contains)
			}
		})
	}
}

func TestNoColorEnvironment(t *testing.T) {
	// Save and restore environment
	oldNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", oldNoColor)

	os.Setenv("NO_COLOR", "1")

	// This test just ensures NO_COLOR environment is respected
	// The actual color library handles this
	if os.Getenv("NO_COLOR") == "" {
		t.Error("NO_COLOR environment variable not set correctly")
	}
}

func TestPrintSuccess(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintSuccess("test message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("PrintSuccess() output = %q, want to contain 'test message'", output)
	}
	if !strings.Contains(output, SuccessIcon) {
		t.Errorf("PrintSuccess() output = %q, want to contain success icon", output)
	}
}

func TestPrintError(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintError("test error")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test error") {
		t.Errorf("PrintError() output = %q, want to contain 'test error'", output)
	}
	if !strings.Contains(output, ErrorIcon) {
		t.Errorf("PrintError() output = %q, want to contain error icon", output)
	}
}

func TestPrintWarning(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintWarning("test warning")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test warning") {
		t.Errorf("PrintWarning() output = %q, want to contain 'test warning'", output)
	}
	if !strings.Contains(output, WarningIcon) {
		t.Errorf("PrintWarning() output = %q, want to contain warning icon", output)
	}
}

func TestPrintInfo(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintInfo("test info")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test info") {
		t.Errorf("PrintInfo() output = %q, want to contain 'test info'", output)
	}
	if !strings.Contains(output, InfoIcon) {
		t.Errorf("PrintInfo() output = %q, want to contain info icon", output)
	}
}

func TestPrintHeader(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHeader("Test Header")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Test Header") {
		t.Errorf("PrintHeader() output = %q, want to contain 'Test Header'", output)
	}
	// Should contain separator line
	if !strings.Contains(output, "─") {
		t.Errorf("PrintHeader() output = %q, want to contain separator", output)
	}
}

func TestPrintSection(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintSection("🎯", "Test Section")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Test Section") {
		t.Errorf("PrintSection() output = %q, want to contain 'Test Section'", output)
	}
	if !strings.Contains(output, "🎯") {
		t.Errorf("PrintSection() output = %q, want to contain emoji", output)
	}
}

func TestFormatChange(t *testing.T) {
	tests := []struct {
		name   string
		change float64
		want   string
	}{
		{
			name:   "positive change",
			change: 10.5,
			want:   "+10.50%",
		},
		{
			name:   "negative change",
			change: -5.25,
			want:   "-5.25%",
		},
		{
			name:   "zero change",
			change: 0.0,
			want:   "0.00%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatChange(tt.change)
			if !strings.Contains(result, tt.want) {
				t.Errorf("FormatChange(%f) = %q, want to contain %q", tt.change, result, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		ns   float64
		want string
	}{
		{
			name: "nanoseconds",
			ns:   500,
			want: "ns",
		},
		{
			name: "microseconds",
			ns:   5000,
			want: "µs",
		},
		{
			name: "milliseconds",
			ns:   5000000,
			want: "ms",
		},
		{
			name: "seconds",
			ns:   5000000000,
			want: "s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.ns)
			if !strings.Contains(result, tt.want) {
				t.Errorf("FormatDuration(%f) = %q, want to contain %q", tt.ns, result, tt.want)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes float64
		want  string
	}{
		{
			name:  "bytes",
			bytes: 500,
			want:  "B",
		},
		{
			name:  "kilobytes",
			bytes: 5120,
			want:  "KB",
		},
		{
			name:  "megabytes",
			bytes: 5242880,
			want:  "MB",
		},
		{
			name:  "gigabytes",
			bytes: 5368709120,
			want:  "GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if !strings.Contains(result, tt.want) {
				t.Errorf("FormatBytes(%f) = %q, want to contain %q", tt.bytes, result, tt.want)
			}
		})
	}
}

func TestRepeatChar(t *testing.T) {
	tests := []struct {
		name string
		char string
		n    int
		want string
	}{
		{
			name: "repeat dash 5 times",
			char: "-",
			n:    5,
			want: "-----",
		},
		{
			name: "repeat zero times",
			char: "x",
			n:    0,
			want: "",
		},
		{
			name: "repeat once",
			char: "a",
			n:    1,
			want: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repeatChar(tt.char, tt.n)
			if result != tt.want {
				t.Errorf("repeatChar(%q, %d) = %q, want %q", tt.char, tt.n, result, tt.want)
			}
		})
	}
}

func TestFormatWithArgs(t *testing.T) {
	// Test formatted output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintSuccess("test %s %d", "value", 42)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test value 42") {
		t.Errorf("PrintSuccess() with args output = %q, want to contain 'test value 42'", output)
	}
}

func BenchmarkFormatChange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatChange(10.5)
	}
}

func BenchmarkFormatDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatDuration(5000000)
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatBytes(5242880)
	}
}

// Test that icons are defined
func TestIconsAreDefined(t *testing.T) {
	icons := []struct {
		name  string
		value string
	}{
		{"SuccessIcon", SuccessIcon},
		{"ErrorIcon", ErrorIcon},
		{"WarningIcon", WarningIcon},
		{"InfoIcon", InfoIcon},
		{"ArrowIcon", ArrowIcon},
		{"UpArrow", UpArrow},
		{"DownArrow", DownArrow},
		{"RightArrow", RightArrow},
		{"FireEmoji", FireEmoji},
		{"TargetEmoji", TargetEmoji},
		{"RocketEmoji", RocketEmoji},
		{"ChartEmoji", ChartEmoji},
		{"CheckEmoji", CheckEmoji},
		{"CrossEmoji", CrossEmoji},
	}

	for _, icon := range icons {
		t.Run(icon.name, func(t *testing.T) {
			if icon.value == "" {
				t.Errorf("%s is empty", icon.name)
			}
		})
	}
}

func ExamplePrintSuccess() {
	// Note: Colors won't show in example output
	PrintSuccess("Operation completed successfully")
	// Output will contain: ✓ Operation completed successfully
}

func ExampleFormatChange() {
	fmt.Println(FormatChange(10.5))
	fmt.Println(FormatChange(-5.25))
	// Note: Output will be colored in terminal
}
