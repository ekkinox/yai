package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ekkinox/yo/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	app := tea.NewProgram(ui.NewUi())

	if _, err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
