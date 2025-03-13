package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AddDownloadTab struct {
	URL string
	// handle input
}

func NewAddDownloadTab() AddDownloadTab {
	return AddDownloadTab{}
}

func (m AddDownloadTab) Init() tea.Cmd {
	return nil
}

func (m AddDownloadTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m AddDownloadTab) View() string {
	return "Add Download Tab"
}
