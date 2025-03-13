package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type tab int

const tabCount = 3

const (
	downloads tab = iota
	queues
	addDownload
)

type MainView struct {
	currentTab     tab
	downloadTab    tea.Model
	queueTab       tea.Model
	addDownloadTab tea.Model
}

func NewMainView() MainView {
	return MainView{
		currentTab:     downloads,
		downloadTab:    NewDownloadsTab(),
		queueTab:       NewQueuesTab(),
		addDownloadTab: NewAddDownloadTab(),
	}
}

func (m MainView) Init() tea.Cmd {
	return nil
}

func (m MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentTab {
	case downloads:
		downloadTab, cmd := m.downloadTab.Update(msg)
		m.downloadTab = downloadTab
		return m, cmd
	case queues:
		queueTab, cmd := m.queueTab.Update(msg)
		m.queueTab = queueTab
		return m, cmd
	case addDownload:
		addDownloadTab, cmd := m.addDownloadTab.Update(msg)
		m.addDownloadTab = addDownloadTab
		return m, cmd
	}

	return m, nil
}

func (m MainView) View() string {
	// TODO: add tab name to the view

	switch m.currentTab {
	case downloads:
		return m.downloadTab.View()
	case queues:
		return m.queueTab.View()
	case addDownload:
		return m.addDownloadTab.View()
	}

	return ""
}

// // cmd

// // msg
