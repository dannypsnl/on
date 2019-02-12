package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	ctxs := []string{os.Args[1]}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("on(%s)\n", prettyContext(ctxs))
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed at read command, error: %s", err)
		}
		input = input[:len(input)-1]
		restCmds := strings.Split(input, " ")
		cmd := exec.Command(ctxs[0], append(ctxs[1:], restCmds...)...)
		result, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		fmt.Println(string(result))
	}
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
