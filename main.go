package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekkinox/yo/tui"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	app := tea.NewProgram(tui.NewTui())

	if _, err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
