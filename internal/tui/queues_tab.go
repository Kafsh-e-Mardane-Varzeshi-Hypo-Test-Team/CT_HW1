package tui

import (
	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QueuesTab struct {
	Queues []models.Queue
	table  table.Model
}

func NewQueuesTab() QueuesTab {
	Queues := models.GetQueues()
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Target Directory", Width: 30},
		{Title: "Max Parallel Downloads", Width: 20},
		{Title: "Speed Limit", Width: 15},
		{Title: "Start Time", Width: 10},
		{Title: "End Time", Width: 10},
	}
	rows := []table.Row{}
	for _, queue := range Queues {
		rows = append(rows, []string{
			queue.Name,
			queue.TargetDirectory,
			fmt.Sprintf("%d", queue.MaxParallelDownloads),
			queue.SpeedLimit,
			queue.StartTime.Format("15:04:05"),
			queue.EndTime.Format("15:04:05"),
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

	return QueuesTab{
		Queues: Queues,
		table:  t,
	}
}

func (m QueuesTab) Init() tea.Cmd { return nil }

func (m QueuesTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	m.table.SetHeight(msg.Height - 3)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m QueuesTab) View() string {
	// TODO: add help text
	return baseStyle.Render(m.table.View()) + "\n"
}
