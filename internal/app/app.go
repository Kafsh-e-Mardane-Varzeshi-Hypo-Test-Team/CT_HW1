package app

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tabView tui.MainView
}

func New() Model {
	return Model{
		tabView: tui.NewMainView(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.tabView.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	returnValue, cmd := m.tabView.Update(msg)
	castValue, ok := returnValue.(tui.MainView)
	if !ok {
		panic("type assertion to tui.TabView failed")
	}
	m.tabView = castValue
	return m, cmd
}

func (m Model) View() string {
	return m.tabView.View()
}
