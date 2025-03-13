package app

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tabView tui.TabView
}

func New() Model {
	return Model{
		tabView: tui.NewTabView(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.tabView.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// assert type before assigning to m.tabView
	returnValue, cmd := m.tabView.Update(msg)
	m.tabView = returnValue.(tui.TabView)
	return m, cmd
}

func (m Model) View() string {
	return m.tabView.View()
}
