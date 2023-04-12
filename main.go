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
	flag.BoolVar(&exec, "exec", false, "Exec mode")
	flag.BoolVar(&chat, "chat", false, "Chat mode")
	flag.Parse()

	args := flag.Args()

	mode := ui.ReplMode
	if len(args) > 0 {
		mode = ui.CliMode
	}

	if _, err := tea.NewProgram(ui.NewUi(mode, strings.Join(args, " "))).Run(); err != nil {
		log.Fatal(err)
	}
}
