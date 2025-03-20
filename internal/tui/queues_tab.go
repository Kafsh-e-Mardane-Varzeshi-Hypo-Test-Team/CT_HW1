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

// Key Bindings
type queuesKeyMap struct {
	Navigation key.Binding
	Delete     key.Binding
	Edit       key.Binding
	NewQueue   key.Binding
	Quit       key.Binding
}

func (k queuesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k queuesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation, k.Quit},
		{k.NewQueue, k.Edit, k.Delete},
	}
}

type QueuesTab struct {
	Queues []models.Queue
	table  table.Model
	help   help.Model
	keys   queuesKeyMap
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

	help := help.New()
	help.ShowAll = true
	help.FullSeparator = " \t "

	return QueuesTab{
		Queues: Queues,
		table:  t,
		help:   help,
		keys: queuesKeyMap{
			Navigation: key.NewBinding(
				key.WithKeys("up", "down", "left", "right"),
				key.WithHelp("↑/↓/←/→", "navigate"),
			),
			Delete: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete"),
			),
			NewQueue: key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "new queue"),
			),
			Edit: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c", "esc", "q"),
				key.WithHelp("ctrl+c/esc", "quit"),
			),
		},
	}
}

func (m QueuesTab) Init() tea.Cmd { return nil }

func (m QueuesTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	m.table.SetHeight(msg.Height - 3)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Navigation):
		case key.Matches(msg, m.keys.Delete):

		case key.Matches(msg, m.keys.NewQueue):

		case key.Matches(msg, m.keys.Edit):

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m QueuesTab) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		baseStyle.Render(m.table.View()),
		helpStyle.Render(m.help.View(m.keys)),
	)
}
