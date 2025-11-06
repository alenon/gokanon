package ui

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ProgressBar wraps the progressbar library with sensible defaults
type ProgressBar struct {
	bar *progressbar.ProgressBar
}

// NewProgressBar creates a new progress bar
func NewProgressBar(max int, description string) *ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	return &ProgressBar{bar: bar}
}

// NewIndeterminateSpinner creates a spinner for indeterminate operations
func NewIndeterminateSpinner(description string) *ProgressBar {
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWidth(10),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	return &ProgressBar{bar: bar}
}

// Add increments the progress bar
func (p *ProgressBar) Add(num int) error {
	return p.bar.Add(num)
}

// Set sets the progress bar to a specific value
func (p *ProgressBar) Set(num int) error {
	return p.bar.Set(num)
}

// Finish completes the progress bar
func (p *ProgressBar) Finish() error {
	return p.bar.Finish()
}

// Clear clears the progress bar
func (p *ProgressBar) Clear() error {
	return p.bar.Clear()
}

// Describe updates the description
func (p *ProgressBar) Describe(description string) {
	p.bar.Describe(description)
}

// Spinner is a simple spinner for long operations
type Spinner struct {
	writer      io.Writer
	message     string
	stopChan    chan bool
	isRunning   bool
	spinChars   []string
	currentChar int
	mu          sync.RWMutex // Protects message field
}

// isCI checks if we're running in a CI environment
func isCI() bool {
	// Check common CI environment variables
	ciVars := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"GITHUB_ACTIONS",
		"GITLAB_CI",
		"CIRCLECI",
		"TRAVIS",
		"JENKINS_URL",
		"BUILDKITE",
		"DRONE",
	}

	for _, v := range ciVars {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return false
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		writer:    os.Stdout,
		message:   message,
		stopChan:  make(chan bool),
		spinChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if s.isRunning {
		return
	}

	s.isRunning = true

	// In CI environments, just print once without spinning
	if isCI() {
		s.mu.RLock()
		msg := s.message
		s.mu.RUnlock()
		fmt.Fprintf(s.writer, "%s %s %s\n",
			Info("⠋"),
			msg,
			Dim("..."))
		return
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopChan:
				return
			case <-ticker.C:
				s.mu.RLock()
				msg := s.message
				s.mu.RUnlock()

				fmt.Fprintf(s.writer, "\r%s %s %s",
					Info(s.spinChars[s.currentChar]),
					msg,
					Dim("..."))
				s.currentChar = (s.currentChar + 1) % len(s.spinChars)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	if !s.isRunning {
		return
	}

	// In CI, we just printed once, so nothing to clean up
	if isCI() {
		s.isRunning = false
		return
	}

	s.stopChan <- true
	s.isRunning = false

	s.mu.RLock()
	msgLen := len(s.message)
	s.mu.RUnlock()

	fmt.Fprintf(s.writer, "\r%s\r", repeatChar(" ", msgLen+20))
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}
