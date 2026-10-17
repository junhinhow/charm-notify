//go:build windows

package notify

import (
	"errors"
	"fmt"
)

// sendPlatform envia a notificação no Windows.
//
// Tenta primeiro o módulo BurntToast (PowerShell) e, se ele não estiver
// disponível ou falhar, usa o msg.exe. Quando os dois caminhos falham, os
// dois erros são devolvidos juntos para que a causa real (por exemplo, o
// BurntToast não instalado) não seja perdida.
func sendPlatform(n *Notification) error {
	errToast := sendBurntToast(n)
	if errToast == nil {
		return nil
	}

	errMsg := sendMsg(n)
	if errMsg == nil {
		return nil
	}

	return errors.Join(
		fmt.Errorf("notify: falha ao enviar a notificação via BurntToast: %w", errToast),
		fmt.Errorf("notify: falha ao enviar a notificação via msg.exe: %w", errMsg),
	)
}

// sendBurntToast usa o módulo BurntToast do PowerShell.
//
// São respeitados o ícone (-Icon) e o som (-Sound). O módulo não oferece
// prazo de expiração, portanto Timeout só é respeitado no fallback msg.exe.
func sendBurntToast(n *Notification) error {
	cmd, err := burntToastCommand(n)
	if err != nil {
		return err
	}

	return runCommand("powershell", "-NoProfile", "-NonInteractive", "-Command", cmd)
}

// sendMsg usa o msg.exe como fallback básico.
//
// O msg.exe não aceita ícones, portanto Icon é ignorado neste caminho. O
// prazo de expiração é informado em segundos inteiros (arredondado para
// cima, com no mínimo 1 segundo).
func sendMsg(n *Notification) error {
	args, err := msgArgs(n)
	if err != nil {
		return err
	}

	return runCommand("msg", args...)
}
