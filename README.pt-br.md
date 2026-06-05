# charm-notify

Notificacoes desktop multiplataforma em Go.

[Read in English](README.md)

## Recursos

- **Windows** — BurntToast (PowerShell) com fallback para `msg.exe`
- **macOS** — Notification Center nativo via `osascript`
- **Linux** — `notify-send` (libnotify)
- Opcoes funcionais: icone, timeout, som
- Zero dependencias (apenas stdlib)

## Instalacao

```bash
go get github.com/junhinhow/charm-notify
```

## Uso

```go
package main

import (
    "log"
    "time"

    notify "github.com/junhinhow/charm-notify"
)

func main() {
    // Notificacao simples
    err := notify.Send("Ola", "Mundo")
    if err != nil {
        log.Fatal(err)
    }

    // Com opcoes
    err = notify.Send("Build Completo", "Todos os testes passaram",
        notify.WithIcon("dialog-information"),
        notify.WithTimeout(10*time.Second),
        notify.WithSound(true),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Detalhes por Plataforma

### Windows

Requer o modulo PowerShell [BurntToast](https://github.com/Windos/BurntToast) para notificacoes toast. Usa `msg.exe` como fallback se BurntToast nao estiver disponivel.

```powershell
Install-Module -Name BurntToast -Force
```

### macOS

Usa AppleScript via `osascript`. Sem dependencias adicionais.

### Linux

Requer `notify-send` do libnotify:

```bash
# Debian/Ubuntu
sudo apt install libnotify-bin

# Arch
sudo pacman -S libnotify

# Fedora
sudo dnf install libnotify
```

## Licenca

MIT
