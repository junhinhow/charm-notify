package notify

import (
	"strconv"
)

// notifySendArgs monta os argumentos usados pelo notify-send (Linux).
//
// O ícone e o prazo de expiração só são enviados quando informados: sem eles
// o notify-send usa o comportamento padrão do servidor de notificações.
func notifySendArgs(n *Notification) ([]string, error) {
	if err := validateSingleLine("título", n.Title); err != nil {
		return nil, err
	}
	if err := validateSingleLine("mensagem", n.Message); err != nil {
		return nil, err
	}

	args := []string{n.Title, n.Message}

	if n.Icon != "" {
		args = append(args, "--icon="+n.Icon)
	}

	if n.Timeout > 0 {
		args = append(args, "--expire-time="+strconv.FormatInt(n.Timeout.Milliseconds(), 10))
	}

	return args, nil
}
