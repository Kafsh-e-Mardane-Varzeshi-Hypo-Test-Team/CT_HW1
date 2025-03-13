package tui

import (
	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type DownloadsTab struct {
	Downloads []models.Download
	table     table.Model
}

func NewDownloadsTab() DownloadsTab {
	Downloads := models.GetDownloads()
	columns := []table.Column{
		{Title: "URL", Width: 30},
		{Title: "Queue", Width: 15},
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

	return DownloadsTab{
		Downloads: Downloads,
		table:     t,
	}
}

func (m DownloadsTab) Init() tea.Cmd { return nil }

func (m DownloadsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	m.table.SetHeight(msg.Height - 3)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m DownloadsTab) View() string {
	// TODO: add help text
	return baseStyle.Render(m.table.View()) + "\n"
}
