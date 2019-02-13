package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

func main() {
	h := History{
		exist:    make(map[string]bool),
		suggests: make([]prompt.Suggest, 0),
		contexts: os.Args[1:],
	}
	LivePrefix = prettyContext(h.contexts)
	p := prompt.New(
		h.executor,
		h.completer,
		prompt.OptionLivePrefix(changeLivePrefix),
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
	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(string(res))

	LivePrefix = prettyContext(h.contexts)
	h.Add(t)
}

func (h *History) updateContext(newCtxs []string) {
	h.contexts = append(h.contexts, newCtxs...)
	LivePrefix = prettyContext(h.contexts)
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

var LivePrefix string

func changeLivePrefix() (string, bool) {
	return fmt.Sprintf("on(%s)> ", LivePrefix), true
}
