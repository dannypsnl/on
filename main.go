package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	prompt "github.com/c-bata/go-prompt"
	"github.com/kr/pty"
)

func main() {
	h := History{
		exist:    make(map[string]bool),
		suggests: make([]prompt.Suggest, 0),
		contexts: os.Args[1:],
	}
	p := prompt.New(
		h.executor,
		h.completer,
		prompt.OptionLivePrefix(h.livePrefix),
		prompt.OptionAddKeyBind(
			prompt.KeyBind{
				Key: prompt.ControlA,
				Fn:  h.onControlA,
			},
		),
	)
	p.Run()
}

type History struct {
	exist    map[string]bool
	suggests []prompt.Suggest
	contexts []string

	waitingNewContext bool
}

func (h *History) onControlA(*prompt.Buffer) {
	h.waitingNewContext = true
}

func (h *History) Add(command string) {
	if !h.exist[command] {
		h.suggests = append(h.suggests, prompt.Suggest{Text: command})
	}
	h.exist[command] = true
}

func (h *History) completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterFuzzy(h.suggests, d.GetWordBeforeCursor(), true)
}

func (h *History) executor(t string) {
	restCmd := strings.Split(t, " ")
	if h.waitingNewContext {
		h.updateContext(restCmd)
		h.waitingNewContext = false
		return
	}
	cmd := exec.Command(h.contexts[0], append(h.contexts[1:], restCmd...)...)
	tty, err := pty.Start(cmd)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer tty.Close()

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("failed at make raw stdin, error: %s\n", err)
		os.Exit(1)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	go func() { io.Copy(tty, os.Stdin) }()
	io.Copy(os.Stdout, tty)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	h.Add(t)
}

func (h *History) updateContext(newCtxs []string) {
	h.contexts = append(h.contexts, newCtxs...)
}

func prettyContext(ctxs []string) string {
	var sb strings.Builder
	for _, ctx := range ctxs {
		sb.WriteString(ctx)
		sb.WriteRune(' ')
	}
	s := sb.String()
	return s[:len(s)-1]
}

func (h *History) livePrefix() (string, bool) {
	return fmt.Sprintf("on(%s)> ", prettyContext(h.contexts)), true
}
