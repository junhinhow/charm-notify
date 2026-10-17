// Package notify fornece notificações desktop multiplataforma.
//
// Suporta Windows (toast/msg), macOS (notification center) e Linux (notify-send).
//
// Nem toda opção funciona em toda plataforma. A tabela de suporte por
// sistema operacional está no README (README.md e README.pt-br.md) e pode ser
// consultada em tempo de execução com [PlatformCapabilities].
package notify

import (
	"runtime"
	"time"
)

// DefaultTimeout é o prazo de exibição padrão das notificações.
const DefaultTimeout = 5 * time.Second

// Notification representa uma notificação desktop.
type Notification struct {
	// Title título da notificação.
	Title string
	// Message corpo da notificação.
	Message string
	// Icon caminho de um ícone. Respeitado no Linux (--icon) e no Windows
	// (BurntToast -Icon). Ignorado no macOS.
	Icon string
	// Timeout prazo de exibição da notificação. Valores maiores que zero são
	// usados como expiração; valores menores ou iguais a zero fazem a
	// notificação usar o prazo padrão da plataforma.
	Timeout time.Duration
	// Sound habilita o som da notificação. Respeitado no macOS e no Windows
	// (BurntToast -Sound). Ignorado no Linux.
	Sound bool
}

// Option configura uma Notification.
type Option func(*Notification)

// WithIcon define o ícone da notificação.
//
// O ícone só é respeitado no Linux e no Windows (via BurntToast). No macOS a
// opção é ignorada porque o Notification Center não permite ícones
// personalizados.
func WithIcon(icon string) Option {
	return func(n *Notification) {
		n.Icon = icon
	}
}

// WithTimeout define o prazo de exibição da notificação.
//
// Valores maiores que zero são usados como expiração. Valores menores ou
// iguais a zero são normalizados para zero e fazem a notificação usar o prazo
// padrão da plataforma.
//
// Atenção: o msg.exe do Windows só aceita segundos inteiros, então prazos
// menores que 1 segundo são arredondados para 1 segundo.
func WithTimeout(d time.Duration) Option {
	return func(n *Notification) {
		if d < 0 {
			d = 0
		}
		n.Timeout = d
	}
}

// WithSound ativa o som da notificação.
//
// O som só é respeitado no macOS e no Windows (via BurntToast). No Linux a
// opção é ignorada porque o notify-send não expõe essa opção.
func WithSound(enabled bool) Option {
	return func(n *Notification) {
		n.Sound = enabled
	}
}

// Send envia uma notificação desktop.
// Detecta a plataforma automaticamente e usa o mecanismo nativo.
//
// O título e a mensagem são obrigatórios e não podem conter quebras de linha.
func Send(title, message string, opts ...Option) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if message == "" {
		return ErrEmptyMessage
	}
	if err := validateSingleLine("título", title); err != nil {
		return err
	}
	if err := validateSingleLine("mensagem", message); err != nil {
		return err
	}

	n := &Notification{
		Title:   title,
		Message: message,
		Timeout: DefaultTimeout,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(n)
		}
	}

	return sendPlatform(n)
}

// Capabilities descreve quais opções são respeitadas na plataforma atual.
type Capabilities struct {
	// Icon informa se o ícone da notificação é respeitado.
	Icon bool
	// Sound informa se o som da notificação é respeitado.
	Sound bool
	// Timeout informa se o prazo de exibição é respeitado.
	Timeout bool
}

// PlatformCapabilities retorna as opções respeitadas na plataforma atual.
//
// No Windows o prazo de exibição só é respeitado no fallback msg.exe (o
// msg.exe arredonda para segundos inteiros) e o ícone/som só são respeitados
// pelo módulo BurntToast.
func PlatformCapabilities() Capabilities {
	switch runtime.GOOS {
	case "linux":
		return Capabilities{Icon: true, Timeout: true}
	case "darwin":
		return Capabilities{Sound: true}
	case "windows":
		return Capabilities{Icon: true, Sound: true, Timeout: true}
	default:
		return Capabilities{}
	}
}

// Platform retorna o identificador da plataforma usada para enviar
// notificações: "linux", "darwin" ou "windows".
//
// Retorna uma string vazia em plataformas sem suporte.
func Platform() string {
	switch runtime.GOOS {
	case "linux":
		return "linux"
	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	default:
		return ""
	}
}

// Supported informa se a plataforma atual suporta notificações desktop.
func Supported() bool {
	return Platform() != ""
}
