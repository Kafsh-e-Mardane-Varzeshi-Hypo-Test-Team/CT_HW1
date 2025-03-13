// package main

// import (
// 	"log"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// func main() {
// 	m := NewModel()
// 	// NewProgram() w/ initial model and program options
// 	p := tea.NewProgram(m)
// 	// Run the program
// 	_, err := p.Run()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }

// type tab int

// const (
// 	downloads tab = iota
// 	queues
// 	addDownload
// )

// // Model: app state
// type MainModel struct {
// 	currentTab     tab
// 	downloadTab    tea.Model
// 	queueTab       tea.Model
// 	addDownloadTab tea.Model
// }

// type DownloadTab struct {
// 	Downloads []Download
// }

// type QueueTab struct {
// 	Queues []Queue
// }

// // NewModel: initial model
// func NewModel() MainModel {
// 	return MainModel{
// 		currentTab: downloads,
// 	}
// }

// // Init: kick off the event loop
// func (m MainModel) Init() tea.Cmd {
// 	return nil
// }

// // Update: handle msgs
// func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "q":
// 			return m, tea.Quit
// 		// switch tabs with numbers
// 		case "1":
// 			m.currentTab = downloads
// 		case "2":
// 			m.currentTab = queues
// 		case "3":
// 			m.currentTab = addDownload
// 		}
// 	}

// 	// handle msgs for each tab
// 	switch m.currentTab {
// 	case downloads:
// 		downloadTab, cmd := m.downloadTab.Update(msg)
// 		m.downloadTab = downloadTab
// 		return m, cmd
// 	case queues:
// 		queueTab, cmd := m.queueTab.Update(msg)
// 		m.queueTab = queueTab
// 		return m, cmd
// 	case addDownload:
// 		addDownloadTab, cmd := m.addDownloadTab.Update(msg)
// 		m.addDownloadTab = addDownloadTab
// 		return m, cmd
// 	}
// 	return m, nil
// }

// // View: return a string based on the state
// func (m MainModel) View() string {
// 	switch m.currentTab {
// 	case downloads:
// 		return m.downloadTab.View()
// 	case queues:
// 		return m.queueTab.View()
// 	case addDownload:
// 		return m.addDownloadTab.View()
// 	}
// 	return ""
// }

// // cmd

// // msg

package main

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(app.New())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
