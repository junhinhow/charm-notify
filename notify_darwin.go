//go:build darwin

package notify

import (
	"fmt"
)

// sendPlatform envia a notificação no macOS via osascript.
//
// O Notification Center do macOS aceita som, mas não aceita ícones
// personalizados, portanto Icon é ignorado aqui. O osascript também não
// oferece prazo de expiração, portanto Timeout é ignorado.
func sendPlatform(n *Notification) error {
	script, err := appleScript(n)
	if err != nil {
		return err
	}

	if err := runCommand("osascript", "-e", script); err != nil {
		return fmt.Errorf("notify: falha ao enviar a notificação via osascript: %w", err)
	}

	return nil
}
