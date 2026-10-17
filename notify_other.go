//go:build !windows && !linux && !darwin

package notify

import (
	"fmt"
	"runtime"
)

// sendPlatform retorna erro em plataformas não suportadas.
func sendPlatform(_ *Notification) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedPlatform, runtime.GOOS)
}
