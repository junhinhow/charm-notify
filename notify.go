// Pacote notify fornece notificacoes desktop multiplataforma.
//
// Suporta Windows (toast/msg), macOS (notification center) e Linux (notify-send).
package notify

import (
	"fmt"
	"time"
)

// Notification representa uma notificacao desktop.
type Notification struct {
	Title   string
	Message string
	Icon    string
	Timeout time.Duration
	Sound   bool
}

// Option configura uma Notification.
type Option func(*Notification)

// WithIcon define o icone da notificacao.
func WithIcon(icon string) Option {
	return func(n *Notification) {
		n.Icon = icon
	}
}

// WithTimeout define o tempo de exibicao da notificacao.
func WithTimeout(d time.Duration) Option {
	return func(n *Notification) {
		n.Timeout = d
	}
}

// WithSound ativa o som da notificacao.
func WithSound(enabled bool) Option {
	return func(n *Notification) {
		n.Sound = enabled
	}
}

// Send envia uma notificacao desktop.
// Detecta a plataforma automaticamente e usa o mecanismo nativo.
func Send(title, message string, opts ...Option) error {
	if title == "" {
		return fmt.Errorf("notify: titulo nao pode ser vazio")
	}
	if message == "" {
		return fmt.Errorf("notify: mensagem nao pode ser vazia")
	}

	n := &Notification{
		Title:   title,
		Message: message,
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(n)
	}

	return sendPlatform(n)
}
