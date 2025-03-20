package tui

import (
	"strings"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const tabCount = 3

const (
	addDownload tab = iota
	downloads
	queues
)

var (
	tabBorder        = lipgloss.RoundedBorder()
	docStyle         = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().
				Border(tabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 2).
				AlignHorizontal(lipgloss.Center)
	activeTabStyle = inactiveTabStyle.
			Foreground(highlightColor).
			Padding(0, 2).
			AlignHorizontal(lipgloss.Center)
	// Background(highlightColor).
	// BorderBackground(highlightColor).
)

type MainView struct {
	currentTab     tab
	manager        *models.Manager
	downloadTab    tea.Model
	queueTab       tea.Model
	addDownloadTab tea.Model
}

func NewMainView(manager *models.Manager) MainView {
	return MainView{
		currentTab:     downloads,
		manager:        manager,
		downloadTab:    NewDownloadsTab(manager),
		queueTab:       NewQueuesTab(manager),
		addDownloadTab: NewAddDownloadTab(manager),
	}
}

func (m MainView) Init() tea.Cmd {
	return nil
}

func (m MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Send keypress to the active tab first
	switch m.currentTab {
	case downloads:
		m.downloadTab, cmd = m.downloadTab.Update(msg)
	case queues:
		m.queueTab, cmd = m.queueTab.Update(msg)
	case addDownload:
		m.addDownloadTab, cmd = m.addDownloadTab.Update(msg)
	}

	if cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// TODO: left and right have different meanings in different tabs (e.g. in addDownloadTab, left and right are used to move cursor)
		case "left":
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "right":
			if m.currentTab < tabCount-1 {
				m.currentTab++
			}
		// TODO: q means different things in different tabs (e.g. in addDownloadTab, q is a character that can be typed)
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m MainView) View() string {
	builder := strings.Builder{}

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

	builder.WriteString(row)
	builder.WriteString("\n")
	builder.WriteString(content)

	return docStyle.Render(builder.String())
}

// // cmd

// // msg
