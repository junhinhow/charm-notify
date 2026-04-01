# charm-notify

Cross-platform desktop notifications from Go.

[Leia em Portugues](README.pt-br.md)

## Features

- **Windows** — BurntToast (PowerShell) with `msg.exe` fallback
- **macOS** — Native Notification Center via `osascript`
- **Linux** — `notify-send` (libnotify)
- Functional options: icon, timeout, sound
- Zero dependencies (stdlib only)

## Install

```bash
go get github.com/junhinhow/charm-notify
```

## Usage

```go
package main

import (
    "log"
    "time"

    notify "github.com/junhinhow/charm-notify"
)

func main() {
    // Simple notification
    err := notify.Send("Hello", "World")
    if err != nil {
        log.Fatal(err)
    }

    // With options
    err = notify.Send("Build Complete", "All tests passed",
        notify.WithIcon("dialog-information"),
        notify.WithTimeout(10*time.Second),
        notify.WithSound(true),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Platform Details

### Windows

Requires [BurntToast](https://github.com/Windos/BurntToast) PowerShell module for toast notifications. Falls back to `msg.exe` if BurntToast is not available.

```powershell
Install-Module -Name BurntToast -Force
```

### macOS

Uses AppleScript via `osascript`. No additional dependencies.

### Linux

Requires `notify-send` from libnotify:

```bash
# Debian/Ubuntu
sudo apt install libnotify-bin

# Arch
sudo pacman -S libnotify

# Fedora
sudo dnf install libnotify
```

## License

MIT
