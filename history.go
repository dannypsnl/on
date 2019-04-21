package main

import (
	"fmt"
	"os/exec"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

func NewHistoryWithContext(contexts []string) *History {
	return &History{
		exist:    make(map[string]map[string]bool),
		suggests: make(map[string][]prompt.Suggest),
		contexts: contexts,
	}
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
	return prompt.FilterFuzzy(h.getSuggests(), d.GetWordBeforeCursor(), true)
}

func (h *History) getSuggests() []prompt.Suggest {
	return h.suggests[h.curContext()]
}

func (h *History) executor(command string) {
	if h.waitingNewContext {
		h.updateContext(command)
		h.waitingNewContext = false
		return
	}

	fullCommand := strings.Join(append(h.contexts, command), " ")
	err := runPipeline(cutCommand(fullCommand)...)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	h.addCommandIntoSuggests(suggestCommand(command))
}

func cutCommand(command string) []*exec.Cmd {
	pipeCmdText := strings.Split(command, "|")
	pipeCmds := make([]*exec.Cmd, 0)
	for _, pipeCmd := range pipeCmdText {
		pipeCmds = append(pipeCmds, newCommand(pipeCmd))
	}
	return pipeCmds
}

func suggestCommand(commandText string) string {
	cs := make([]string, 0)
	for _, e := range strings.Split(commandText, " ") {
		if e != "" {
			cs = append(cs, e)
		}
	}
	return strings.Join(cs, " ")
}

func (h *History) addCommandIntoSuggests(command string) {
	if command == "" {
		return
	}
	curCtx := h.curContext()
	if h.exist[curCtx] == nil {
		h.exist[curCtx] = make(map[string]bool)
	}
	if !h.exist[curCtx][command] {
		h.suggests[curCtx] = append(h.suggests[curCtx], prompt.Suggest{Text: command})
	}
	h.exist[curCtx][command] = true
}

func (h *History) updateContext(command string) {
	newCtxs := strings.Split(command, " ")
	for _, c := range newCtxs {
		if c != "" {
			h.contexts = append(h.contexts, c)
		}
	}
}

func (h *History) livePrefix() (string, bool) {
	return fmt.Sprintf("on(%s)> ", h.curContext()), true
}

func (h *History) curContext() string {
	return strings.Join(h.contexts, " ")
}
