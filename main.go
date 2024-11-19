package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"hichamlamine.github.io/sonik/model"
)

func main() {
	p := tea.NewProgram(model.Model{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
