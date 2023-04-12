package main

import (
	"flag"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/ekkinox/yo/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var exec, chat bool
	flag.BoolVar(&exec, "e", false, "exec prompt mode")
	flag.BoolVar(&chat, "c", false, "chat prompt mode")
	flag.Parse()

	args := flag.Args()

	runMode := ui.ReplMode
	if len(args) > 0 {
		runMode = ui.CliMode
	}

	promptMode := ui.DefaultPromptMode
	if exec && !chat {
		promptMode = ui.ExecPromptMode
	} else if !exec && chat {
		promptMode = ui.ChatPromptMode
	}

	ui := ui.NewUi(runMode, promptMode, strings.Join(args, " "))

	if _, err := tea.NewProgram(ui).Run(); err != nil {
		log.Fatal(err)
	}
}
