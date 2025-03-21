package tui

import (
	"fmt"
	"time"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	downloads    []*models.DownloadInfo
	table        table.Model
	help         help.Model
	keys         downloadsKeyMap
	footerString string
}

func NewDownloadsTab(manager *models.Manager) DownloadsTab {
	columns := []table.Column{
		{Title: "URL", Width: 30},
		{Title: "Queue", Width: 20},
		{Title: "Status", Width: 15},
		{Title: "Transfer Rate", Width: 15},
		{Title: "Progress", Width: 10},
	}

	rows := []table.Row{}

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

	downloadsTab := DownloadsTab{
		manager:   manager,
		downloads: nil,
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

	downloadsTab.updateRows()

	return downloadsTab
}

func (m DownloadsTab) Init() tea.Cmd {
	return tickUpdate()
}

func (m DownloadsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.updateRows()

	switch msg := msg.(type) {
	case updateMsg:
		m.footerString = fmt.Sprint(m.footerString, "1")
		return m, tickUpdate()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Navigation):
		case key.Matches(msg, m.keys.Pause):
			if m.table.Cursor() >= 0 && m.table.Cursor() < len(m.downloads) {
				dl := m.downloads[m.table.Cursor()]
				switch dl.Status {
				case models.InProgress:
					m.manager.PauseDownload(dl.ID)
					m.updateRows()
				case models.Paused:
					m.manager.ResumeDownload(dl.ID)
					m.updateRows()
				}
			}
		case key.Matches(msg, m.keys.Retry):
			if m.table.Cursor() >= 0 && m.table.Cursor() < len(m.downloads) {
				dl := m.downloads[m.table.Cursor()]
				if dl.Status == models.Failed {
					m.manager.ResumeDownload(dl.ID)
					m.updateRows()
				}
			}
		case key.Matches(msg, m.keys.Delete):
			if m.table.Cursor() >= 0 && m.table.Cursor() < len(m.downloads) {
				dl := m.downloads[m.table.Cursor()]
				m.manager.RemoveDownload(dl.ID)
				m.updateRows()
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m DownloadsTab) View() string {
	if len(m.table.Rows()) != 0 {
		row := m.table.Cursor()
		status := m.downloads[row].Status

		// Update the help view
		switch status {
		case models.InProgress:
			m.keys.Retry.SetEnabled(false)
			m.keys.Pause.SetEnabled(true)
		case models.Paused:
			m.keys.Retry.SetEnabled(false)
			m.keys.Pause.SetEnabled(true)
		case models.Failed:
			m.keys.Retry.SetEnabled(true)
			m.keys.Pause.SetEnabled(false)
		case models.Completed:
			m.keys.Retry.SetEnabled(false)
			m.keys.Pause.SetEnabled(false)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		borderedStyle.Render(m.table.View()),
		noStyle.Render(m.footerString),
		helpStyle.Render(m.help.View(m.keys)),
	)
}

func (m *DownloadsTab) updateRows() {
	m.footerString = fmt.Sprint(m.footerString, ".")
	m.downloads = m.manager.GetDownloadList()
	rows := []table.Row{}
	for _, download := range m.downloads {
		status := download.Status

		var statusString string
		switch status {
		case models.InProgress:
			statusString = "Downloading"
		case models.Paused:
			statusString = "Paused"
		case models.Completed:
			statusString = "Completed"
		case models.Failed:
			statusString = "Failed"
		}

		rows = append(rows, []string{
			download.URL,
			download.QueueName,
			statusString,
			speedString(download.TransferRate),
			fmt.Sprintf("%#6.2f%%", download.Progress*100),
		})
	}

	m.table.SetRows(rows)
}

func speedString(speed int64) string {
	if speed < 1024 {
		return fmt.Sprintf("%d B/s", speed)
	} else if speed < 1024*1024 {
		return fmt.Sprintf("%d KB/s", speed/1024)
	} else {
		return fmt.Sprintf("%d MB/s", speed/1024/1024)
	}
}

// update loop
type updateMsg struct{}

func tickUpdate() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return updateMsg{}
	})
}
