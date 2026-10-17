package notify

import "errors"

// Erros retornados pela pacote. Use errors.Is para verificá-los.
var (
	// ErrEmptyTitle indica que o título da notificação não foi informado.
	ErrEmptyTitle = errors.New("notify: título não pode ser vazio")

	// ErrEmptyMessage indica que a mensagem da notificação não foi informada.
	ErrEmptyMessage = errors.New("notify: mensagem não pode ser vazia")

	// ErrMultilineText indica que o título ou a mensagem continham quebras
	// de linha. Quebras de linha são rejeitadas porque viram comandos
	// multilinha nas plataformas que montam scripts (AppleScript/PowerShell).
	ErrMultilineText = errors.New("notify: texto não pode conter quebras de linha")

	// ErrUnsupportedPlatform indica que a plataforma atual não possui
	// mecanismo de notificação desktop.
	ErrUnsupportedPlatform = errors.New("notify: notificações não suportadas nesta plataforma")

	// ErrCommandNotFound indica que o executável responsável pela
	// notificação não foi encontrado no PATH.
	ErrCommandNotFound = errors.New("notify: executável não encontrado no PATH")
)
