package embedded

import (
	"embed"
)

//go:embed gokanon.bash
//go:embed gokanon.zsh
//go:embed gokanon.fish
var CompletionScripts embed.FS
