//go:build !windows && !linux && !darwin

package notify

import "fmt"

// sendPlatform retorna erro em plataformas nao suportadas.
func sendPlatform(n *Notification) error {
	return fmt.Errorf("notify: notificacoes nao suportadas nesta plataforma")
}
