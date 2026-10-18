# charm-notify

Notificações desktop multiplataforma em Go.

[Read in English](README.md)

## Recursos

- **Windows** — BurntToast (PowerShell) com fallback para `msg.exe`
- **macOS** — Notification Center nativo via `osascript`
- **Linux** — `notify-send` (libnotify)
- Opções funcionais: ícone, timeout, som (o suporte varia por plataforma — veja a tabela abaixo)
- Os erros identificam o comando que falhou e carregam o `stderr` capturado
- Zero dependências (apenas stdlib)

## Instalação

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
    // Notificação simples
    err := notify.Send("Ola", "Mundo")
    if err != nil {
        log.Fatal(err)
    }

    // Com opções
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

## Suporte por plataforma

Nem toda opção é respeitada em todo sistema operacional. A tabela abaixo
mostra o que cada backend realmente faz:

| Opção | Linux (`notify-send`) | macOS (`osascript`) | Windows (BurntToast) | Windows (fallback `msg.exe`) |
|---|---|---|---|---|
| `WithIcon` | ✅ `--icon` | ❌ o Notification Center não aceita ícones personalizados | ✅ `-Icon` | ❌ |
| `WithSound` | ❌ o `notify-send` não tem opção de som | ✅ `sound name "default"` | ✅ `-Sound Default` | ❌ |
| `WithTimeout` | ✅ `--expire-time` (milissegundos) | ❌ o `osascript` não tem timeout | ❌ | ✅ `/TIME` (segundos inteiros) |

No Windows o toast é tentado primeiro com o BurntToast e cai para o `msg.exe`;
cada caminho respeita um conjunto diferente de opções.

É possível consultar o suporte da plataforma atual em tempo de execução:

```go
fmt.Println(notify.Platform())            // "linux", "darwin", "windows" ou ""
fmt.Println(notify.Supported())            // true quando a plataforma tem backend
fmt.Println(notify.PlatformCapabilities()) // ex.: {Icon:true Sound:false Timeout:true}
```

### Semântica do timeout

`WithTimeout(d)` se comporta assim:

- `d > 0` — a notificação expira depois de `d`.
- `d <= 0` — a notificação usa o padrão da plataforma; nenhuma flag de
  expiração explícita é enviada. Valores negativos são normalizados para zero.

O fallback `msg.exe` do Windows só aceita segundos inteiros, então
`WithTimeout(500*time.Millisecond)` vira **1 segundo** lá (antes um bug de
truncamento transformava em 5 segundos).

### Restrições de texto

O título e a mensagem são obrigatórios e **não podem conter quebras de linha**.
`Send` devolve um erro que envolve `notify.ErrMultilineText` quando existem,
porque uma quebra de linha transformaria o payload em um comando multilinha nas
plataformas que montam um script (AppleScript/PowerShell).

### Erros

Toda falha identifica o comando que falhou e inclui o `stderr` capturado, em
vez de um `exit status 1` seco:

```go
var cmdErr *notify.CommandError
if errors.As(err, &cmdErr) {
    log.Printf("comando %q falhou: %s", cmdErr.Command, cmdErr.Stderr)
}
```

Os erros sentinela são `notify.ErrEmptyTitle`, `notify.ErrEmptyMessage`,
`notify.ErrMultilineText`, `notify.ErrUnsupportedPlatform` e
`notify.ErrCommandNotFound`; use `errors.Is` para verificá-los.

No Windows, quando o BurntToast e o `msg.exe` falham, os dois erros são
juntados com `errors.Join` para que a causa real (por exemplo, "o BurntToast
não está instalado") nunca seja perdida.

## Detalhes por Plataforma

### Windows

Requer o módulo PowerShell [BurntToast](https://github.com/Windos/BurntToast) para notificações toast. Usa `msg.exe` como fallback se o BurntToast não estiver disponível.

```powershell
Install-Module -Name BurntToast -Force
```

### macOS

Usa AppleScript via `osascript`. Sem dependências adicionais.

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

## Licença

MIT
