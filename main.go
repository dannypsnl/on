package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	prompt "github.com/c-bata/go-prompt"
	"github.com/kr/pty"
)

func main() {
	h := History{
		exist:    make(map[string]map[string]bool),
		suggests: make(map[string][]prompt.Suggest),
		contexts: os.Args[1:],
	}
	h.exist[h.curContext()] = make(map[string]bool)
	h.exist[h.curContext()][""] = true
	p := prompt.New(
		h.executor,
		h.completer,
		prompt.OptionLivePrefix(h.livePrefix),
		prompt.OptionAddKeyBind(
			prompt.KeyBind{Key: prompt.ControlA, Fn: h.onControlA},
			prompt.KeyBind{Key: prompt.ControlC, Fn: h.removeLastElementFromContext},
		),
	)
	p.Run()
}

type History struct {
	exist    map[string]map[string]bool
	suggests map[string][]prompt.Suggest
	contexts []string

	waitingNewContext bool
}

func (h *History) onControlA(*prompt.Buffer) {
	h.waitingNewContext = true
}
func (h *History) removeLastElementFromContext(*prompt.Buffer) {
	if len(h.contexts) > 0 {
		h.contexts = h.contexts[:len(h.contexts)-1]
	}
}

func (h *History) completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterFuzzy(h.suggests[h.curContext()], d.GetWordBeforeCursor(), true)
}

func (h *History) executor(command string) {
	restCmd := strings.Split(command, " ")
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

	listenOnStdin := make(chan os.Signal, 1)
	signal.Notify(listenOnStdin)
	go func() {
		for signalC := range listenOnStdin {
			tty.WriteString(signalC.String())
		}
	}()
	io.Copy(os.Stdout, tty)
	signal.Stop(listenOnStdin)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("error: %s\n", err)
	} else {
		h.addCommandIntoSuggests(command)
	}
}

func (h *History) addCommandIntoSuggests(command string) {
	curCtx := h.curContext()
	if h.exist[curCtx] == nil {
		h.exist[curCtx] = make(map[string]bool)
	}
	if !h.exist[curCtx][command] {
		h.suggests[curCtx] = append(h.suggests[curCtx], prompt.Suggest{Text: command})
	}
	h.exist[curCtx][command] = true
}

func (h *History) updateContext(newCtxs []string) {
	for _, c := range newCtxs {
		if c != "" {
			h.contexts = append(h.contexts, c)
		}
	}
}

func prettyContext(ctxs []string) string {
	if len(ctxs) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, ctx := range ctxs {
		sb.WriteString(ctx)
		sb.WriteRune(' ')
	}
	s := sb.String()
	return s[:len(s)-1]
}

func (h *History) livePrefix() (string, bool) {
	return fmt.Sprintf("on(%s)> ", h.curContext()), true
}

func (h *History) curContext() string {
	return prettyContext(h.contexts)
}
