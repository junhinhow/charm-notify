package notify

import (
	"fmt"
)

// appleScript monta o comando AppleScript usado no macOS.
//
// O Notification Center do macOS não aceita ícones personalizados, portanto a
// opção Icon é ignorada nesta plataforma.
func appleScript(n *Notification) (string, error) {
	if err := validateSingleLine("título", n.Title); err != nil {
		return "", err
	}
	if err := validateSingleLine("mensagem", n.Message); err != nil {
		return "", err
	}

	script := fmt.Sprintf(
		`display notification "%s" with title "%s"`,
		escapeAppleScript(n.Message),
		escapeAppleScript(n.Title),
	)

	if n.Sound {
		script += ` sound name "default"`
	}

	return script, nil
}
