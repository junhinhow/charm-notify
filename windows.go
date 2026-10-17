package notify

import (
	"fmt"
	"math"
	"time"
)

// timeoutSeconds converte um prazo para segundos inteiros.
//
// O msg.exe só aceita segundos inteiros, então o valor é arredondado para
// cima e nunca fica abaixo de 1 segundo. Valores não positivos significam
// "usar o prazo padrão do sistema" e são devolvidos como zero.
func timeoutSeconds(d time.Duration) int {
	if d <= 0 {
		return 0
	}

	seconds := int(math.Ceil(d.Seconds()))
	if seconds < 1 {
		seconds = 1
	}

	return seconds
}

// burntToastCommand monta o comando PowerShell que usa o módulo BurntToast.
//
// O BurntToast aceita ícone (-Icon) e som (-Sound), mas não aceita prazo de
// expiração: essa opção só é respeitada no fallback msg.exe.
func burntToastCommand(n *Notification) (string, error) {
	if err := validateSingleLine("título", n.Title); err != nil {
		return "", err
	}
	if err := validateSingleLine("mensagem", n.Message); err != nil {
		return "", err
	}

	cmd := fmt.Sprintf(
		"New-BurntToastNotification -Text '%s','%s'",
		escapePowerShell(n.Title),
		escapePowerShell(n.Message),
	)

	if n.Icon != "" {
		cmd += fmt.Sprintf(" -Icon '%s'", escapePowerShell(n.Icon))
	}

	if n.Sound {
		cmd += " -Sound Default"
	}

	return cmd, nil
}

// msgArgs monta os argumentos usados pelo msg.exe (fallback do Windows).
//
// O msg.exe não aceita ícones, portanto a opção Icon é ignorada neste
// caminho. O prazo de expiração é informado em segundos inteiros.
func msgArgs(n *Notification) ([]string, error) {
	if err := validateSingleLine("título", n.Title); err != nil {
		return nil, err
	}
	if err := validateSingleLine("mensagem", n.Message); err != nil {
		return nil, err
	}

	args := []string{"*"}

	if seconds := timeoutSeconds(n.Timeout); seconds > 0 {
		args = append(args, fmt.Sprintf("/TIME:%d", seconds))
	}

	args = append(args, fmt.Sprintf("%s: %s", n.Title, n.Message))

	return args, nil
}
