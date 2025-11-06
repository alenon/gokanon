package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	// Disable colors if NO_COLOR environment variable is set or not a TTY
	NoColor = os.Getenv("NO_COLOR") != "" || !isTerminal()

	// Color functions
	Success = color.New(color.FgGreen, color.Bold).SprintFunc()
	Error   = color.New(color.FgRed, color.Bold).SprintFunc()
	Warning = color.New(color.FgYellow, color.Bold).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
	Dim     = color.New(color.Faint).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()

	// Status indicators
	SuccessIcon = "✓"
	ErrorIcon   = "✗"
	WarningIcon = "⚠"
	InfoIcon    = "ℹ"
	ArrowIcon   = "→"

	// Trend indicators
	UpArrow    = "↑"
	DownArrow  = "↓"
	RightArrow = "→"

	// Emoji indicators (for enhanced output)
	FireEmoji   = "🔥"
	TargetEmoji = "🎯"
	RocketEmoji = "🚀"
	ChartEmoji  = "📊"
	CheckEmoji  = "✅"
	CrossEmoji  = "❌"
)

func init() {
	if NoColor {
		color.NoColor = true
	}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// PrintSuccess prints a success message with a green checkmark
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Success(SuccessIcon), fmt.Sprintf(format, args...))
}

// PrintError prints an error message with a red X
func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s\n", Error(ErrorIcon), fmt.Sprintf(format, args...))
}

// PrintWarning prints a warning message with a yellow warning sign
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Warning(WarningIcon), fmt.Sprintf(format, args...))
}

// PrintInfo prints an info message with a cyan info icon
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Info(InfoIcon), fmt.Sprintf(format, args...))
}

// PrintHeader prints a bold header
func PrintHeader(text string) {
	fmt.Println()
	fmt.Println(Bold(text))
	fmt.Println(Dim(repeatChar("─", len(text))))
}

// PrintSection prints a section header
func PrintSection(emoji, title string) {
	fmt.Printf("\n%s %s\n", emoji, Bold(title))
}

// FormatChange formats a performance change with appropriate color
func FormatChange(percentChange float64) string {
	if percentChange > 0 {
		return Error(fmt.Sprintf("+%.2f%%", percentChange))
	} else if percentChange < 0 {
		return Success(fmt.Sprintf("%.2f%%", percentChange))
	}
	return Dim(fmt.Sprintf("%.2f%%", percentChange))
}

// FormatDuration formats a duration with color based on magnitude
func FormatDuration(ns float64) string {
	if ns < 1000 {
		return Info(fmt.Sprintf("%.2f ns", ns))
	} else if ns < 1000000 {
		return Info(fmt.Sprintf("%.2f µs", ns/1000))
	} else if ns < 1000000000 {
		return Warning(fmt.Sprintf("%.2f ms", ns/1000000))
	}
	return Error(fmt.Sprintf("%.2f s", ns/1000000000))
}

// FormatBytes formats bytes with appropriate units and color
func FormatBytes(bytes float64) string {
	if bytes < 1024 {
		return Info(fmt.Sprintf("%.0f B", bytes))
	} else if bytes < 1024*1024 {
		return Info(fmt.Sprintf("%.2f KB", bytes/1024))
	} else if bytes < 1024*1024*1024 {
		return Warning(fmt.Sprintf("%.2f MB", bytes/(1024*1024)))
	}
	return Error(fmt.Sprintf("%.2f GB", bytes/(1024*1024*1024)))
}

// repeatChar repeats a character n times
func repeatChar(char string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += char
	}
	return result
}
