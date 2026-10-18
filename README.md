# charm-notify

Cross-platform desktop notifications from Go.

[Leia em Português](README.pt-br.md)

## Features

- **Windows** — BurntToast (PowerShell) with `msg.exe` fallback
- **macOS** — Native Notification Center via `osascript`
- **Linux** — `notify-send` (libnotify)
- Functional options: icon, timeout, sound (support varies per platform — see the table below)
- Errors identify the failing command and carry its captured `stderr`
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

## Platform support

Not every option is honoured on every operating system. The table below shows
what each backend actually does:

| Option | Linux (`notify-send`) | macOS (`osascript`) | Windows (BurntToast) | Windows (`msg.exe` fallback) |
|---|---|---|---|---|
| `WithIcon` | ✅ `--icon` | ❌ Notification Center has no custom icons | ✅ `-Icon` | ❌ |
| `WithSound` | ❌ `notify-send` has no sound option | ✅ `sound name "default"` | ✅ `-Sound Default` | ❌ |
| `WithTimeout` | ✅ `--expire-time` (milliseconds) | ❌ `osascript` has no timeout | ❌ | ✅ `/TIME` (whole seconds) |

On Windows the toast is attempted with BurntToast first and falls back to
`msg.exe`; each path honours a different subset of the options.

You can inspect the support of the current platform at runtime:

```go
fmt.Println(notify.Platform())            // "linux", "darwin", "windows" or ""
fmt.Println(notify.Supported())            // true when the platform has a backend
fmt.Println(notify.PlatformCapabilities()) // e.g. {Icon:true Sound:false Timeout:true}
```

### Timeout semantics

`WithTimeout(d)` behaves as follows:

- `d > 0` — the notification expires after `d`.
- `d <= 0` — the notification uses the platform default; no explicit expiry
  flag is sent. Negative values are normalised to zero.

The Windows `msg.exe` fallback only accepts whole seconds, so
`WithTimeout(500*time.Millisecond)` becomes **1 second** there (previously a
truncation bug turned it into 5 seconds).

### Text constraints

The title and the message are mandatory and **must not contain line breaks**.
`Send` returns an error wrapping `notify.ErrMultilineText` when they do,
because a line break would turn the payload into a multi-line command on the
platforms that build a script (AppleScript/PowerShell).

### Errors

Every failure identifies the command that failed and includes its captured
`stderr` instead of a bare `exit status 1`:

```go
var cmdErr *notify.CommandError
if errors.As(err, &cmdErr) {
    log.Printf("command %q failed: %s", cmdErr.Command, cmdErr.Stderr)
}
```

The sentinel errors are `notify.ErrEmptyTitle`, `notify.ErrEmptyMessage`,
`notify.ErrMultilineText`, `notify.ErrUnsupportedPlatform` and
`notify.ErrCommandNotFound`; use `errors.Is` to test them.

On Windows, when both BurntToast and `msg.exe` fail, the two errors are
joined with `errors.Join` so the real cause (for example, "BurntToast is not
installed") is never lost.

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
