package notify

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// execCommand permite injetar o executor de comandos em testes.
// Em produção aponta para exec.Command.
var execCommand = exec.Command

// lookPath permite injetar a busca de executáveis em testes.
// Em produção aponta para exec.LookPath.
var lookPath = exec.LookPath

// CommandError descreve a falha na execução de um comando externo.
//
// O conteúdo capturado do stderr é incluído na mensagem porque várias
// ferramentas de notificação (osascript, powershell, msg) escrevem os
// diagnósticos mais úteis justamente nessa saída.
type CommandError struct {
	// Command nome do executável que falhou.
	Command string
	// Args argumentos passados ao executável.
	Args []string
	// Stderr conteúdo capturado da saída de erro (pode ser vazio).
	Stderr string
	// Err erro original devolvido pelo sistema operacional.
	Err error
}

// Error implementa a interface error.
func (e *CommandError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "comando %q falhou", e.Command)
	if e.Stderr != "" {
		fmt.Fprintf(&b, ": %s", e.Stderr)
	} else {
		b.WriteString(": sem saída de erro")
	}
	if e.Err != nil {
		fmt.Fprintf(&b, ": %v", e.Err)
	}
	return b.String()
}

// Unwrap expõe o erro original para errors.Is/errors.As.
func (e *CommandError) Unwrap() error { return e.Err }

// runCommand executa um comando externo capturando o stderr.
//
// O executável é antes localizado com exec.LookPath para devolver uma
// mensagem útil quando ele não está instalado.
func runCommand(name string, args ...string) error {
	path, err := lookPath(name)
	if err != nil {
		return &CommandError{
			Command: name,
			Args:    args,
			Err:     fmt.Errorf("%w: %q", ErrCommandNotFound, name),
		}
	}

	cmd := execCommand(path, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return &CommandError{
			Command: name,
			Args:    args,
			Stderr:  strings.TrimSpace(stderr.String()),
			Err:     err,
		}
	}

	return nil
}
