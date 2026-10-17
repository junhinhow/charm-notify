package notify

import (
	"errors"
	"fmt"
	"time"
)

func ExampleSend() {
	// Notificação simples.
	if err := Send("Olá", "Mundo"); err != nil {
		fmt.Println("erro:", err)
	}

	// Notificação com opções. O suporte de ícone, som e timeout varia por
	// plataforma — consulte notify.PlatformCapabilities().
	err := Send(
		"Build concluído",
		"Todos os testes passaram",
		WithIcon("dialog-information"),
		WithTimeout(10*time.Second),
		WithSound(true),
	)
	if err != nil {
		fmt.Println("erro:", err)
	}
}

func ExamplePlatform() {
	fmt.Println(Platform())
	fmt.Println(Supported())
	fmt.Println(PlatformCapabilities())
}

func ExampleWithTimeout() {
	n := &Notification{Timeout: DefaultTimeout}

	// Prazo válido: a notificação expira depois de 2s.
	WithTimeout(2 * time.Second)(n)
	fmt.Println(n.Timeout)

	// Valores não positivos usam o prazo padrão da plataforma.
	WithTimeout(-1)(n)
	fmt.Println(n.Timeout)
}

func ExampleCommandError() {
	err := Send("", "mensagem")

	// Erros sentinela podem ser testados com errors.Is.
	if errors.Is(err, ErrEmptyTitle) {
		fmt.Println("informe um título")
	}

	// Falhas de comando carregam o executável e o stderr capturado.
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		fmt.Println(cmdErr.Command, cmdErr.Stderr)
	}
}
