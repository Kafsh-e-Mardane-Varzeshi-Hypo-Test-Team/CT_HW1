package tui

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

type DownloadsTab struct {
	Downloads []models.Download
}

func NewDownloadsTab() DownloadsTab {
	return DownloadsTab{}
}

func (m DownloadsTab) Init() tea.Cmd { return nil }

func (m DownloadsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m DownloadsTab) View() string {
	return "Downloads Tab"
}
