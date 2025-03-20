package tui

import (
	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Key Bindings
type downloadsKeyMap struct {
	Navigation key.Binding
	Delete     key.Binding
	Pause      key.Binding
	Retry      key.Binding
	Quit       key.Binding
}

func (k downloadsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k downloadsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation, k.Quit},
		{k.Delete, k.Pause, k.Retry},
	}
}

type DownloadsTab struct {
	manager      *models.Manager
	downloads    []models.Download
	table        table.Model
	help         help.Model
	keys         downloadsKeyMap
	footerString string
}

func NewDownloadsTab(manager *models.Manager) DownloadsTab {
	Downloads := models.GetDownloads()
	columns := []table.Column{
		{Title: "URL", Width: 30},
		{Title: "Queue", Width: 20},
		{Title: "Status", Width: 15},
		{Title: "Transfer Rate", Width: 15},
		{Title: "Progress", Width: 10},
	}
	rows := []table.Row{}
	for _, download := range Downloads {
		rows = append(rows, []string{
			download.Url,
			download.Queue,
			download.Status,
			download.TransferRate,
			fmt.Sprintf("%#6.2f%%", download.Progress*100),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	help := help.New()
	help.ShowAll = true
	help.FullSeparator = " \t "

	return DownloadsTab{
		manager:   manager,
		downloads: Downloads,
		table:     t,
		help:      help,
		keys: downloadsKeyMap{
			Navigation: key.NewBinding(
				key.WithKeys("up", "down", "left", "right"),
				key.WithHelp("↑/↓/←/→", "navigate"),
			),
			Delete: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete"),
			),
			Pause: key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp("p", "pause/resume"),
			),
			Retry: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "retry"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("ctrl+c/esc", "quit"),
			),
		},
		footerString: "",
	}
}

func (m DownloadsTab) Init() tea.Cmd { return nil }

func (m DownloadsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Navigation):
		case key.Matches(msg, m.keys.Pause):
			// case "p":
			// 	row := m.table.Cursor()
			// 	// case switch for selected row
			// 	// if status is downloading, pause it
			// 	// if status is paused, resume it
		case key.Matches(msg, m.keys.Retry):
			// case "r":
			// 	row := m.table.Cursor()
			// 	// case switch for selected row
			// 	// if status is failed, retry it
		case key.Matches(msg, m.keys.Delete):
			// case "d":
			// 	row := m.table.Cursor()
			// 	// case switch for selected row
			// 	// delete the download
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m DownloadsTab) View() string {
	row := m.table.Cursor()
	status := m.downloads[row].Status

	// Update the help view
	switch status {
	case "Downloading":
		m.keys.Retry.SetEnabled(false)
		m.keys.Pause.SetEnabled(true)
	case "Paused":
		m.keys.Retry.SetEnabled(false)
		m.keys.Pause.SetEnabled(true)
	case "Failed":
		m.keys.Retry.SetEnabled(true)
		m.keys.Pause.SetEnabled(false)
	case "Completed":
		m.keys.Retry.SetEnabled(false)
		m.keys.Pause.SetEnabled(false)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		baseStyle.Render(m.table.View()),
		noStyle.Render(m.footerString),
		helpStyle.Render(m.help.View(m.keys)),
	)
}
