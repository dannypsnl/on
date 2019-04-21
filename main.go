package main

import (
	"os"

	prompt "github.com/c-bata/go-prompt"
)

func main() {
	h := NewHistoryWithContext(os.Args[1:])

	p := prompt.New(
		h.executor,
		h.completer,
		prompt.OptionLivePrefix(h.livePrefix),
		prompt.OptionAddKeyBind(
			prompt.KeyBind{Key: prompt.ControlA, Fn: func(*prompt.Buffer) {
				h.startWaitingAppendingContext()
			}},
			prompt.KeyBind{Key: prompt.ControlC, Fn: func(*prompt.Buffer) {
				h.removeLastElementFromContext()
			}},
		),
	)
	p.Run()
}
