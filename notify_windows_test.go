//go:build windows

package notify

import (
	"errors"
	"strings"
	"testing"
)

func TestSendPlataformaUneOsDoisErrosDoWindows(t *testing.T) {
	installMissingLookPath(t)

	err := sendPlatform(&Notification{Title: "T", Message: "M"})
	if err == nil {
		t.Fatal("esperava erro quando nenhum dos dois caminhos funciona")
	}

	if !strings.Contains(err.Error(), "BurntToast") {
		t.Errorf("erro deveria mencionar o BurntToast, obteve %q", err.Error())
	}
	if !strings.Contains(err.Error(), "msg.exe") {
		t.Errorf("erro deveria mencionar o msg.exe, obteve %q", err.Error())
	}
	if !errors.Is(err, ErrCommandNotFound) {
		t.Errorf("esperava ErrCommandNotFound, obteve %v", err)
	}
}
