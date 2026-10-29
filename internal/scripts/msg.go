package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		msg := strings.TrimSpace(scanner.Text())
		if msg == "" {
			continue
		}

		cmd := exec.Command(
			"go",
			"test",
			"-count=1",
			"-tags=integration",
			"./internal",
			"-run", "^TestPing",
			"-args", msg,
		)
		cmd.Env = append(os.Environ(), "MSG="+msg)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		go func() {
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
			}
		}()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "stdin:", err)
		os.Exit(1)
	}
}
