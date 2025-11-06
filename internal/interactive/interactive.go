package interactive

import (
	"fmt"
	"io"
	"strings"

	"github.com/alenon/gokanon/internal/ui"
	"github.com/chzyer/readline"
)

// Session represents an interactive gokanon session
type Session struct {
	rl       *readline.Instance
	commands map[string]CommandHandler
}

// CommandHandler is a function that handles a command
type CommandHandler func(args []string) error

// New creates a new interactive session
func New() (*Session, error) {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("run",
			readline.PcItem("-bench="),
			readline.PcItem("-pkg="),
			readline.PcItem("-profile="),
			readline.PcItem("-benchtime="),
			readline.PcItem("-count="),
		),
		readline.PcItem("list"),
		readline.PcItem("compare",
			readline.PcItem("--latest"),
		),
		readline.PcItem("export",
			readline.PcItem("--latest"),
			readline.PcItem("-format=html"),
			readline.PcItem("-format=csv"),
			readline.PcItem("-format=markdown"),
			readline.PcItem("-format=json"),
		),
		readline.PcItem("stats",
			readline.PcItem("-last="),
		),
		readline.PcItem("trend",
			readline.PcItem("-last="),
		),
		readline.PcItem("check",
			readline.PcItem("--latest"),
			readline.PcItem("-threshold="),
		),
		readline.PcItem("flamegraph"),
		readline.PcItem("serve",
			readline.PcItem("-port="),
		),
		readline.PcItem("delete"),
		readline.PcItem("doctor"),
		readline.PcItem("help"),
		readline.PcItem("clear"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ui.Info("gokanon> "),
		HistoryFile:     "/tmp/gokanon_history.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		return nil, err
	}

	return &Session{
		rl:       rl,
		commands: make(map[string]CommandHandler),
	}, nil
}

// RegisterCommand registers a command handler
func (s *Session) RegisterCommand(name string, handler CommandHandler) {
	s.commands[name] = handler
}

// Run starts the interactive session
func (s *Session) Run() error {
	defer s.rl.Close()

	s.printWelcome()

	for {
		line, err := s.rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle built-in commands
		if s.handleBuiltIn(line) {
			continue
		}

		// Parse command and arguments
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		// Execute command
		handler, exists := s.commands[command]
		if !exists {
			ui.PrintError("Unknown command: %s", command)
			fmt.Println("Type 'help' for available commands")
			continue
		}

		if err := handler(args); err != nil {
			ui.PrintError("Command failed: %v", err)
		}
	}

	s.printGoodbye()
	return nil
}

func (s *Session) printWelcome() {
	fmt.Println()
	fmt.Println(ui.Bold("╔════════════════════════════════════════════════════════════╗"))
	fmt.Println(ui.Bold("║") + "          " + ui.Info("Welcome to gokanon Interactive Mode!") + "          " + ui.Bold("║"))
	fmt.Println(ui.Bold("╚════════════════════════════════════════════════════════════╝"))
	fmt.Println()
	fmt.Println(ui.Dim("  Type 'help' for available commands"))
	fmt.Println(ui.Dim("  Use TAB for auto-completion"))
	fmt.Println(ui.Dim("  Press Ctrl+C or type 'exit' to quit"))
	fmt.Println()
}

func (s *Session) printGoodbye() {
	fmt.Println()
	fmt.Printf("%s Thanks for using gokanon! Happy benchmarking! %s\n", ui.Success(ui.SuccessIcon), ui.RocketEmoji)
	fmt.Println()
}

func (s *Session) handleBuiltIn(line string) bool {
	switch line {
	case "exit", "quit":
		return true
	case "clear", "cls":
		// Clear screen by printing ANSI escape code
		fmt.Print("\033[H\033[2J")
		return true
	case "help", "?":
		s.printHelp()
		return true
	default:
		return false
	}
}

func (s *Session) printHelp() {
	fmt.Println()
	ui.PrintSection(ui.InfoIcon, "Available Commands")
	fmt.Println()

	commands := []struct {
		name        string
		description string
	}{
		{"run", "Run benchmarks and save results"},
		{"list", "List all saved benchmark results"},
		{"compare", "Compare two benchmark results"},
		{"export", "Export comparison results to various formats"},
		{"stats", "Show statistical analysis of multiple runs"},
		{"trend", "Analyze performance trends over time"},
		{"check", "Check performance against thresholds"},
		{"flamegraph", "View CPU/memory flame graphs"},
		{"serve", "Start interactive web dashboard"},
		{"delete", "Delete a benchmark result"},
		{"doctor", "Run diagnostics"},
		{"help", "Show this help message"},
		{"clear", "Clear the screen"},
		{"exit", "Exit interactive mode"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %s %-12s %s\n",
			ui.Info(ui.ArrowIcon),
			ui.Bold(cmd.name),
			ui.Dim(cmd.description))
	}

	fmt.Println()
	fmt.Println(ui.Dim("For command-specific options, use: <command> -h"))
	fmt.Println()
}

// Close closes the interactive session
func (s *Session) Close() error {
	return s.rl.Close()
}
