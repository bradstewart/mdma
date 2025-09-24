package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bradstewart/mdma/internal/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("MDMA - Music DJ Management Application")
		fmt.Println("Version: 0.1.0")
		return
	}

	model := ui.NewModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
