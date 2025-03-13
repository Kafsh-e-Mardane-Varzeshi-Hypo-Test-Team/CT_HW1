package tui

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

type QueuesTab struct {
	Downloads []models.Download
}

func NewQueuesTab() QueuesTab {
	return QueuesTab{}
}

func (m QueuesTab) Init() tea.Cmd { return nil }

func (m QueuesTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m QueuesTab) View() string {
	return "Queues Tab"
}
