package main

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	manager := models.NewManager()
	p := tea.NewProgram(tui.NewMainView(manager))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
