package notify

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestHelperProcess reexecuta o binário de teste como um comando auxiliar.
// Só age quando GO_WANT_HELPER_PROCESS está definida.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if msg := os.Getenv("GO_HELPER_STDERR"); msg != "" {
		fmt.Fprint(os.Stderr, msg)
	}

	code, _ := strconv.Atoi(os.Getenv("GO_HELPER_EXIT"))
	os.Exit(code)
}

// fakeExec substitui o executor de comandos por um processo auxiliar que
// devolve o stderr e o código de saída desejados.
type fakeExec struct {
	command string
	args    []string
}

// installFakeExec injeta o executor falso e devolve uma função que restaura
// o executor original.
func installFakeExec(t *testing.T, stderr string, exitCode int) *fakeExec {
	t.Helper()

	fake := &fakeExec{}
	oldExec, oldLook := execCommand, lookPath

	execCommand = func(name string, args ...string) *exec.Cmd {
		fake.command = name
		fake.args = args

		cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
		cmd.Env = append(os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			"GO_HELPER_STDERR="+stderr,
			"GO_HELPER_EXIT="+strconv.Itoa(exitCode),
		)
		return cmd
	}
	lookPath = func(name string) (string, error) { return name, nil }

	t.Cleanup(func() {
		execCommand = oldExec
		lookPath = oldLook
	})

	return fake
}

// installMissingLookPath faz com que toda busca de executável falhe.
func installMissingLookPath(t *testing.T) {
	t.Helper()

	oldLook := lookPath
	lookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}

	t.Cleanup(func() { lookPath = oldLook })
}

func TestSendTituloVazio(t *testing.T) {
	err := Send("", "mensagem")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("esperava ErrEmptyTitle, obteve %v", err)
	}
}

func TestSendMensagemVazia(t *testing.T) {
	err := Send("titulo", "")
	if !errors.Is(err, ErrEmptyMessage) {
		t.Fatalf("esperava ErrEmptyMessage, obteve %v", err)
	}
}

func TestSendRejeitaQuebrasDeLinha(t *testing.T) {
	tests := []struct {
		nome     string
		titulo   string
		mensagem string
	}{
		{"titulo com \n", "linha1\nlinha2", "mensagem"},
		{"titulo com \r", "linha1\rlinha2", "mensagem"},
		{"mensagem com \n", "titulo", "linha1\nlinha2"},
		{"mensagem com \r", "titulo", "linha1\rlinha2"},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			err := Send(tt.titulo, tt.mensagem)
			if !errors.Is(err, ErrMultilineText) {
				t.Fatalf("esperava ErrMultilineText, obteve %v", err)
			}
			if !strings.Contains(err.Error(), "quebras de linha") {
				t.Errorf("erro deveria explicar o motivo, obteve %q", err.Error())
			}
		})
	}
}

func TestSendIgnoraOpcaoNula(t *testing.T) {
	// Uma opção nula não pode provocar panic.
	installFakeExec(t, "", 0)

	if err := Send("titulo", "mensagem", nil); err != nil {
		t.Fatalf("esperava sucesso, obteve %v", err)
	}
}

func TestWithTimeoutNormalizaValores(t *testing.T) {
	tests := []struct {
		nome     string
		entrada  time.Duration
		esperado time.Duration
	}{
		{"positivo", 10 * time.Second, 10 * time.Second},
		{"sub-segundo", 500 * time.Millisecond, 500 * time.Millisecond},
		{"zero", 0, 0},
		{"negativo", -5 * time.Second, 0},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			n := &Notification{Timeout: DefaultTimeout}
			WithTimeout(tt.entrada)(n)
			if n.Timeout != tt.esperado {
				t.Fatalf("esperava timeout %v, obteve %v", tt.esperado, n.Timeout)
			}
		})
	}
}

func TestTimeoutSecondsArredondaParaCima(t *testing.T) {
	tests := []struct {
		entrada  time.Duration
		esperado int
	}{
		{0, 0},
		{-time.Second, 0},
		{500 * time.Millisecond, 1},
		{time.Second, 1},
		{1500 * time.Millisecond, 2},
		{10 * time.Second, 10},
	}

	for _, tt := range tests {
		if obtido := timeoutSeconds(tt.entrada); obtido != tt.esperado {
			t.Errorf("timeoutSeconds(%v): esperava %d, obteve %d", tt.entrada, tt.esperado, obtido)
		}
	}
}

func TestPlataformaDetectada(t *testing.T) {
	esperado := ""
	switch runtime.GOOS {
	case "linux", "darwin", "windows":
		esperado = runtime.GOOS
	}

	if Platform() != esperado {
		t.Errorf("esperava plataforma %q, obteve %q", esperado, Platform())
	}

	if Supported() != (esperado != "") {
		t.Errorf("Supported() inconsistente com Platform(): %v", Supported())
	}

	capacidades := PlatformCapabilities()
	if esperado == "" {
		if capacidades != (Capabilities{}) {
			t.Errorf("plataforma sem suporte não deveria ter capacidades, obteve %+v", capacidades)
		}
		return
	}

	if !capacidades.Timeout && esperado == "linux" {
		t.Errorf("Linux deveria respeitar o prazo de exibição, obteve %+v", capacidades)
	}
	if !capacidades.Sound && esperado == "darwin" {
		t.Errorf("macOS deveria respeitar o som, obteve %+v", capacidades)
	}
	if !capacidades.Icon && esperado == "linux" {
		t.Errorf("Linux deveria respeitar o ícone, obteve %+v", capacidades)
	}
}

// countUnescapedQuotes conta as aspas duplas que NÃO estão escapadas.
func countUnescapedQuotes(s string) int {
	total := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '"' {
			continue
		}
		barras := 0
		for j := i - 1; j >= 0 && s[j] == '\\'; j-- {
			barras++
		}
		if barras%2 == 0 {
			total++
		}
	}
	return total
}

// payloadsAdversariais são textos que tentam escapar do literal do comando.
var payloadsAdversariais = []string{
	`"`,
	`""`,
	`" ; do shell script "rm -rf /" ; "`,
	`$(whoami)`,
	"`whoami`",
	`'; rm -rf / ; '`,
	`"; rm -rf / ; "`,
	`\`,
	`\\`,
	`\n`,
	`a"b'c$d;e|f&g`,
	"linha1\nlinha2",
	"linha1\rlinha2",
}

func TestEscapeAppleScriptResistePayloadsAdversariais(t *testing.T) {
	for _, payload := range payloadsAdversariais {
		t.Run(payload, func(t *testing.T) {
			escapado := escapeAppleScript(payload)

			if got := countUnescapedQuotes(escapado); got != 0 {
				t.Errorf("escapado contém %d aspas sem escape: %q", got, escapado)
			}

			// Cada barra invertida do payload precisa ter sido duplicada.
			esperado := 2*strings.Count(payload, `\`) + strings.Count(payload, `"`)
			if got := strings.Count(escapado, `\`); got != esperado {
				t.Errorf("esperava %d barras invertidas, obteve %d em %q", esperado, got, escapado)
			}
		})
	}
}

func TestAppleScriptMantemLiteralFechado(t *testing.T) {
	for _, payload := range payloadsAdversariais {
		t.Run(payload, func(t *testing.T) {
			n := &Notification{Title: payload, Message: payload}

			script, err := appleScript(n)
			if strings.ContainsAny(payload, "\r\n") {
				if !errors.Is(err, ErrMultilineText) {
					t.Fatalf("esperava ErrMultilineText, obteve %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			// 4 aspas sem escape: abrindo/fechando mensagem e título.
			if got := countUnescapedQuotes(script); got != 4 {
				t.Errorf("esperava 4 aspas sem escape, obteve %d em %q", got, script)
			}
		})
	}
}

func TestAppleScriptComSom(t *testing.T) {
	script, err := appleScript(&Notification{Title: "T", Message: "M", Sound: true})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !strings.Contains(script, `sound name "default"`) {
		t.Errorf("script deveria conter o som: %q", script)
	}
}

func TestEscapePowerShellResistePayloadsAdversariais(t *testing.T) {
	for _, payload := range payloadsAdversariais {
		t.Run(payload, func(t *testing.T) {
			escapado := escapePowerShell(payload)

			// Toda aspa simples do payload precisa aparecer duplicada.
			esperado := 2 * strings.Count(payload, `'`)
			if got := strings.Count(escapado, `'`); got != esperado {
				t.Errorf("esperava %d aspas simples, obteve %d em %q", esperado, got, escapado)
			}
		})
	}
}

func TestBurntToastComandoMantemLiteralFechado(t *testing.T) {
	for _, payload := range payloadsAdversariais {
		t.Run(payload, func(t *testing.T) {
			n := &Notification{Title: payload, Message: payload}

			cmd, err := burntToastCommand(n)
			if strings.ContainsAny(payload, "\r\n") {
				if !errors.Is(err, ErrMultilineText) {
					t.Fatalf("esperava ErrMultilineText, obteve %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			// 4 aspas simples (delimitadores) + as aspas do payload, duplicadas
			// no título e na mensagem.
			esperado := 4 + 4*strings.Count(payload, `'`)
			if got := strings.Count(cmd, `'`); got != esperado {
				t.Errorf("esperava %d aspas simples, obteve %d em %q", esperado, got, cmd)
			}
		})
	}
}

func TestBurntToastComandoComIconeESom(t *testing.T) {
	cmd, err := burntToastCommand(&Notification{
		Title:   "T",
		Message: "M",
		Icon:    "C:/icons/app.ico",
		Sound:   true,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !strings.Contains(cmd, "-Icon 'C:/icons/app.ico'") {
		t.Errorf("comando deveria conter o ícone: %q", cmd)
	}
	if !strings.Contains(cmd, "-Sound Default") {
		t.Errorf("comando deveria conter o som: %q", cmd)
	}
}

func TestBurntToastComandoSemIconeNemSom(t *testing.T) {
	cmd, err := burntToastCommand(&Notification{Title: "T", Message: "M"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strings.Contains(cmd, "-Icon") {
		t.Errorf("comando não deveria conter ícone: %q", cmd)
	}
	if strings.Contains(cmd, "-Sound") {
		t.Errorf("comando não deveria conter som: %q", cmd)
	}
}

func TestNotifySendArgs(t *testing.T) {
	args, err := notifySendArgs(&Notification{
		Title:   "T",
		Message: "M",
		Icon:    "dialog-information",
		Timeout: 1500 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	esperado := []string{"T", "M", "--icon=dialog-information", "--expire-time=1500"}
	if len(args) != len(esperado) {
		t.Fatalf("esperava %v, obteve %v", esperado, args)
	}
	for i := range esperado {
		if args[i] != esperado[i] {
			t.Errorf("argumento %d: esperava %q, obteve %q", i, esperado[i], args[i])
		}
	}
}

func TestNotifySendArgsSemOpcoes(t *testing.T) {
	args, err := notifySendArgs(&Notification{Title: "T", Message: "M", Timeout: 0})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(args) != 2 {
		t.Errorf("esperava apenas título e mensagem, obteve %v", args)
	}
}

func TestMsgArgsArredondaTimeout(t *testing.T) {
	args, err := msgArgs(&Notification{Title: "T", Message: "M", Timeout: 500 * time.Millisecond})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	esperado := []string{"*", "/TIME:1", "T: M"}
	if len(args) != len(esperado) {
		t.Fatalf("esperava %v, obteve %v", esperado, args)
	}
	for i := range esperado {
		if args[i] != esperado[i] {
			t.Errorf("argumento %d: esperava %q, obteve %q", i, esperado[i], args[i])
		}
	}
}

func TestMsgArgsSemTimeout(t *testing.T) {
	args, err := msgArgs(&Notification{Title: "T", Message: "M"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(args) != 2 {
		t.Errorf("esperava apenas destino e texto, obteve %v", args)
	}
	if args[0] != "*" || args[1] != "T: M" {
		t.Errorf("argumentos inesperados: %v", args)
	}
}

func TestRunCommandCapturaStderr(t *testing.T) {
	installFakeExec(t, "detalhe do erro do comando", 1)

	err := runCommand("osascript", "-e", "display notification")
	if err == nil {
		t.Fatal("esperava erro quando o comando falha")
	}

	var cmdErr *CommandError
	if !errors.As(err, &cmdErr) {
		t.Fatalf("esperava *CommandError, obteve %T (%v)", err, err)
	}
	if cmdErr.Command != "osascript" {
		t.Errorf("esperava comando %q, obteve %q", "osascript", cmdErr.Command)
	}
	if !strings.Contains(cmdErr.Stderr, "detalhe do erro do comando") {
		t.Errorf("stderr deveria ser capturado, obteve %q", cmdErr.Stderr)
	}
	if !strings.Contains(err.Error(), "detalhe do erro do comando") {
		t.Errorf("mensagem de erro deveria conter o stderr, obteve %q", err.Error())
	}
	if !strings.Contains(err.Error(), "osascript") {
		t.Errorf("mensagem de erro deveria identificar o comando, obteve %q", err.Error())
	}
}

func TestRunCommandExecutavelAusente(t *testing.T) {
	installMissingLookPath(t)

	err := runCommand("notify-send", "T", "M")
	if err == nil {
		t.Fatal("esperava erro quando o executável não existe")
	}
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("esperava ErrCommandNotFound, obteve %v", err)
	}
	if !strings.Contains(err.Error(), "notify-send") {
		t.Errorf("mensagem deveria identificar o executável, obteve %q", err.Error())
	}
}

func TestRunCommandSucesso(t *testing.T) {
	installFakeExec(t, "", 0)

	if err := runCommand("notify-send", "T", "M"); err != nil {
		t.Fatalf("esperava sucesso, obteve %v", err)
	}
}

func TestSendUsaExecutavelDaPlataforma(t *testing.T) {
	fake := installFakeExec(t, "", 0)

	esperado := map[string]string{
		"linux":   "notify-send",
		"darwin":  "osascript",
		"windows": "powershell",
	}

	comando, ok := esperado[Platform()]
	if !ok {
		t.Skipf("plataforma sem suporte: %s", runtime.GOOS)
	}

	if err := Send("titulo", "mensagem"); err != nil {
		t.Fatalf("esperava sucesso, obteve %v", err)
	}

	if fake.command != comando {
		t.Errorf("esperava o executável %q, obteve %q", comando, fake.command)
	}
}

func TestSendPlataformaComExecutavelAusente(t *testing.T) {
	installMissingLookPath(t)

	err := sendPlatform(&Notification{Title: "T", Message: "M", Timeout: time.Second})
	if err == nil {
		t.Fatal("esperava erro quando o executável não existe")
	}

	if Platform() == "" {
		if !errors.Is(err, ErrUnsupportedPlatform) {
			t.Fatalf("esperava ErrUnsupportedPlatform, obteve %v", err)
		}
		if !strings.Contains(err.Error(), runtime.GOOS) {
			t.Errorf("erro deveria informar o sistema operacional, obteve %q", err.Error())
		}
		return
	}

	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("esperava ErrCommandNotFound, obteve %v", err)
	}
}

func TestSendPlataformaPropagaErroDeValidacao(t *testing.T) {
	err := sendPlatform(&Notification{Title: "T", Message: "linha1\nlinha2"})
	if !errors.Is(err, ErrMultilineText) {
		t.Fatalf("esperava ErrMultilineText, obteve %v", err)
	}
}
