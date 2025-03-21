package tui

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	addDownload tab = iota
	downloads
	queues
)

type MainView struct {
	currentTab     tab
	manager        *models.Manager
	downloadTab    tea.Model
	queueTab       tea.Model
	addDownloadTab tea.Model
	footerString   string
}

func NewMainView(manager *models.Manager) MainView {
	return MainView{
		currentTab:     downloads,
		manager:        manager,
		downloadTab:    NewDownloadsTab(manager),
		queueTab:       NewQueuesTab(manager),
		addDownloadTab: NewAddDownloadTab(manager),
		footerString:   "",
	}
}

func (m MainView) Init() tea.Cmd {
	if m.currentTab == downloads {
		return tickUpdate()
	}
	return nil
}

func (m MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(updateMsg); ok {
		m.downloadTab, cmd = m.downloadTab.Update(msg)
		return m, cmd
	}

	// Send keypress to the active tab first
	switch m.currentTab {
	case downloads:
		m.downloadTab, cmd = m.downloadTab.Update(msg)
		if cmd != nil {
			return m, cmd
		}
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "left":
				m.currentTab = addDownload
				m.addDownloadTab, cmd = m.addDownloadTab.Update(nil)
			case "right":
				m.currentTab = queues
				m.queueTab, cmd = m.queueTab.Update(nil)
			case "esc", "ctrl+c":
				return m, tea.Quit
			}
		}
	case queues:
		m.queueTab, cmd = m.queueTab.Update(msg)
		if cmd != nil {
			return m, cmd
		}
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "left":
				m.currentTab = downloads
				m.downloadTab, cmd = m.downloadTab.Update(nil)
			case "right":
			case "esc", "ctrl+c":
				return m, tea.Quit
			}
		}
	case addDownload:
		m.addDownloadTab, cmd = m.addDownloadTab.Update(msg)
		if cmd != nil {
			return m, cmd
		}
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "left":
			case "right":
				m.currentTab = downloads
				m.downloadTab, cmd = m.downloadTab.Update(nil)
			case "esc", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m MainView) View() string {
	// Render tabs
	tabs := []string{"Add Download", "Downloads List", "Queues List"}
	var renderedTabs []string

	for i, t := range tabs {
		if m.currentTab == tab(i) {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(t))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	var content string
	switch m.currentTab {
	case downloads:
		content = m.downloadTab.View()
	case queues:
		content = m.queueTab.View()
	case addDownload:
		content = m.addDownloadTab.View()
	}
	return docStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		row,
		content,
		noStyle.Width(100).Render(m.footerString),
	))
}

type CloseChildMsg struct{}
