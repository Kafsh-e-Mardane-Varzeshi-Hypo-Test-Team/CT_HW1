package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TabView struct {
	activeTab int
	tabs      []tea.Model
	tabBar    []string
}

func NewTabView() TabView {
	return TabView{
		activeTab: 0,
		tabs: []tea.Model{
			NewDownloadsTab(),
			// NewQueuesTab(),
			// NewAddTab(),
		},
		tabBar: []string{
			"Downloads",
			// "Queues",
			// "Add",
		},
	}
}

func (m TabView) Init() tea.Cmd {
	return nil
}

func (m TabView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.activeTab > 0 {
				m.activeTab--
			}
		case "right":
			if m.activeTab < len(m.tabs)-1 {
				m.activeTab++
			}
		}
	}

	// Update the active tab
	var cmd tea.Cmd
	m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
	return m, cmd
}

func (m TabView) View() string {
	m.activeTab = 0
	return m.tabs[m.activeTab].(DownloadsTab).View()
}
