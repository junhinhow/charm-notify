//go:build linux

package notify

import (
	"fmt"
	"os/exec"
	"strconv"
)

// sendPlatform envia notificacao no Linux via notify-send.
func sendPlatform(n *Notification) error {
	_, err := exec.LookPath("notify-send")
	if err != nil {
		return fmt.Errorf("notify: notify-send nao encontrado — instale libnotify (ex: apt install libnotify-bin)")
	}

	args := []string{n.Title, n.Message}

	if n.Icon != "" {
		args = append(args, "--icon="+n.Icon)
	}

	if n.Timeout > 0 {
		ms := strconv.FormatInt(n.Timeout.Milliseconds(), 10)
		args = append(args, "--expire-time="+ms)
	}

	return exec.Command("notify-send", args...).Run()
}
