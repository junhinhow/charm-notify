//go:build windows

package notify

import (
	"fmt"
	"os/exec"
	"strings"
)

// sendPlatform envia notificacao no Windows.
// Tenta BurntToast (PowerShell) primeiro, fallback para msg.exe.
func sendPlatform(n *Notification) error {
	err := sendBurntToast(n)
	if err != nil {
		return sendMsg(n)
	}
	return nil
}

// sendBurntToast usa o modulo BurntToast do PowerShell.
func sendBurntToast(n *Notification) error {
	// Escapa aspas simples duplicando-as (padrao PowerShell)
	title := strings.ReplaceAll(n.Title, "'", "''")
	message := strings.ReplaceAll(n.Message, "'", "''")

	cmd := fmt.Sprintf("New-BurntToastNotification -Text '%s','%s'", title, message)

	if n.Sound {
		cmd += " -Sound Default"
	}

	return exec.Command("powershell", "-Command", cmd).Run()
}

// sendMsg usa msg.exe como fallback basico.
func sendMsg(n *Notification) error {
	timeout := int(n.Timeout.Seconds())
	if timeout < 1 {
		timeout = 5
	}

	text := fmt.Sprintf("%s: %s", n.Title, n.Message)
	return exec.Command("msg", "*", fmt.Sprintf("/TIME:%d", timeout), text).Run()
}
