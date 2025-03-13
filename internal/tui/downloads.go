package tui

import tea "github.com/charmbracelet/bubbletea"

type DownloadsTab struct{}

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
