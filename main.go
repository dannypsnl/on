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
			prompt.KeyBind{Key: prompt.ControlA, Fn: h.onControlA},
			prompt.KeyBind{Key: prompt.ControlC, Fn: h.removeLastElementFromContext},
		),
	)
	p.Run()
}
