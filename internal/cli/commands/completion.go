package commands

import (
	"fmt"
	"os"

	"github.com/alenon/gokanon/embedded"
	"github.com/alenon/gokanon/internal/ui"
)

// Completion installs shell completion scripts
func Completion() error {
	if len(os.Args) < 3 {
		fmt.Println(ui.Bold("gokanon completion - Install shell completion"))
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  gokanon completion <shell>")
		fmt.Println()
		fmt.Println("Available shells:")
		fmt.Println("  bash    Bash completion")
		fmt.Println("  zsh     Zsh completion")
		fmt.Println("  fish    Fish completion")
		fmt.Println()
		fmt.Println("Installation instructions:")
		fmt.Println()
		fmt.Println(ui.Bold("Bash:"))
		fmt.Println("  gokanon completion bash > /etc/bash_completion.d/gokanon")
		fmt.Println("  # or for user-only installation:")
		fmt.Println("  gokanon completion bash > ~/.local/share/bash-completion/completions/gokanon")
		fmt.Println()
		fmt.Println(ui.Bold("Zsh:"))
		fmt.Println("  gokanon completion zsh > /usr/local/share/zsh/site-functions/_gokanon")
		fmt.Println("  # or add to your .zshrc:")
		fmt.Println("  source <(gokanon completion zsh)")
		fmt.Println()
		fmt.Println(ui.Bold("Fish:"))
		fmt.Println("  gokanon completion fish > ~/.config/fish/completions/gokanon.fish")
		fmt.Println()
		return nil
	}

	shell := os.Args[2]

	// Read the completion script from embedded file
	scriptPath := ""
	switch shell {
	case "bash":
		scriptPath = "gokanon.bash"
	case "zsh":
		scriptPath = "gokanon.zsh"
	case "fish":
		scriptPath = "gokanon.fish"
	default:
		return ui.NewError(
			fmt.Sprintf("Unsupported shell: %s", shell),
			nil,
			"Supported shells: bash, zsh, fish",
			"Example: gokanon completion bash",
		)
	}

	// Read the completion script from embedded files
	content, err := embedded.CompletionScripts.ReadFile(scriptPath)
	if err != nil {
		return ui.NewError(
			"Failed to read completion script",
			err,
			"Ensure gokanon is properly installed",
			"You may need to reinstall gokanon",
		)
	}

	fmt.Print(string(content))
	return nil
}
