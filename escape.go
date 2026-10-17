package notify

import (
	"fmt"
	"strings"
)

// validateSingleLine garante que o texto não contenha quebras de linha.
//
// Quebras de linha são perigosas nas plataformas que montam um script
// (AppleScript e PowerShell): elas transformam o payload em um comando
// multilinha. O campo é usado na mensagem de erro para indicar qual dos
// textos foi recusado.
func validateSingleLine(field, value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("%w: %s", ErrMultilineText, field)
	}
	return nil
}

// escapeAppleScript escapa um texto para uso dentro de um literal AppleScript
// delimitado por aspas duplas.
//
// Dentro de literais AppleScript o caractere de escape é a barra invertida,
// portanto ela precisa ser duplicada antes das aspas duplas.
func escapeAppleScript(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `"`, `\"`)
}

// escapePowerShell escapa um texto para uso dentro de um literal PowerShell
// delimitado por aspas simples.
//
// Dentro de literais de aspas simples do PowerShell a única caractere que
// precisa ser escapado é a própria aspas simples, que é duplicada. A barra
// invertida não tem significado especial nesses literais.
func escapePowerShell(s string) string {
	return strings.ReplaceAll(s, `'`, `''`)
}
