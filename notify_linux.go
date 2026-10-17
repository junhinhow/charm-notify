//go:build linux

package notify

import (
	"errors"
	"fmt"
)

// sendPlatform envia a notificação no Linux via notify-send.
//
// São respeitados o ícone (--icon) e o prazo de expiração (--expire-time).
// O notify-send não expõe opção de som, portanto Sound é ignorado aqui.
func sendPlatform(n *Notification) error {
	args, err := notifySendArgs(n)
	if err != nil {
		return err
	}

	if err := runCommand("notify-send", args...); err != nil {
		if errors.Is(err, ErrCommandNotFound) {
			return fmt.Errorf(
				"notify: notify-send não encontrado — instale o libnotify (ex: apt install libnotify-bin): %w",
				err,
			)
		}
		return fmt.Errorf("notify: falha ao enviar a notificação via notify-send: %w", err)
	}

	return nil
}
