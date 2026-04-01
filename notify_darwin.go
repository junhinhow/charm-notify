//go:build darwin

package notify

import (
	"fmt"
	"os/exec"
	"strings"
)

// sendPlatform envia notificacao no macOS via osascript (AppleScript).
func sendPlatform(n *Notification) error {
	// Escapa aspas duplas e barras invertidas para AppleScript
	title := strings.ReplaceAll(n.Title, `\`, `\\`)
	title = strings.ReplaceAll(title, `"`, `\"`)

	message := strings.ReplaceAll(n.Message, `\`, `\\`)
	message = strings.ReplaceAll(message, `"`, `\"`)

	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)

	if n.Sound {
		script += ` sound name "default"`
	}

	return exec.Command("osascript", "-e", script).Run()
}
