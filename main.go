package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	prompt "github.com/c-bata/go-prompt"
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

func (h *History) generateCommand(restCmd string) *exec.Cmd {
	cmdText := strings.Join(append(h.contexts, restCmd), " ")
	return newCommand(cmdText)
}

func (h *History) executor(command string) {
	restOfCmd := strings.Split(command, " ")
	if h.waitingNewContext {
		h.updateContext(restOfCmd)
		h.waitingNewContext = false
		return
	}

	pipeCmdText := strings.Split(command, "|")
	pipeCmds := make([]*exec.Cmd, 0)
	pipeCmds = append(pipeCmds, h.generateCommand(pipeCmdText[0]))
	for _, pipeCmd := range pipeCmdText[1:] {
		pipeCmds = append(pipeCmds, newCommand(pipeCmd))
	}
	err := runPipeline(pipeCmds...)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	h.addCommandIntoSuggests(suggestCommand(command))
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

func runPipeline(cmdPipe ...*exec.Cmd) error {
	if len(cmdPipe) <= 0 {
		return fmt.Errorf("no commands for pipeline")
	}
	if len(cmdPipe) == 1 {
		onlyCmd := cmdPipe[0]
		onlyCmd.Stderr = os.Stderr
		onlyCmd.Stdout = os.Stdout
		onlyCmd.Stdin = os.Stdin
		return onlyCmd.Run()
	}
	prevCmd := cmdPipe[0]
	for _, curCmd := range cmdPipe[1:] {
		curCmd.Stdin, _ = prevCmd.StdoutPipe()
		err := prevCmd.Start()
		if err != nil {
			return err
		}
		prevCmd = curCmd
	}
	lastCmd := cmdPipe[len(cmdPipe)-1]
	lastCmd.Stdout = os.Stdout
	lastCmd.Stderr = os.Stderr
	return lastCmd.Run()
}

func newCommand(commandText string) *exec.Cmd {
	cmd := make([]string, 0)
	for _, elem := range strings.Split(commandText, " ") {
		if elem != "" {
			cmd = append(cmd, elem)
		}
	}
	return exec.Command(cmd[0], cmd[1:]...)
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

func (h *History) livePrefix() (string, bool) {
	return fmt.Sprintf("on(%s)> ", h.curContext()), true
}

func (h *History) curContext() string {
	return strings.Join(h.contexts, " ")
}
