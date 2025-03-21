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
	manager      *models.Manager
	queues       []*models.QueueInfo
	table        table.Model
	help         help.Model
	keys         queuesKeyMap
	addQueueTab  tea.Model
	editQueueTab tea.Model
	addingQueue  bool
	editingQueue bool
	footerString string
}

func NewQueuesTab(manager *models.Manager) QueuesTab {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Target Directory", Width: 30},
		{Title: "Max Parallel", Width: 15},
		{Title: "Speed Limit", Width: 15},
		{Title: "Start Time", Width: 10},
		{Title: "End Time", Width: 10},
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

	queuesTab := QueuesTab{
		manager:      manager,
		queues:       nil,
		table:        t,
		addQueueTab:  NewAddQueueTab(manager),
		addingQueue:  false,
		editQueueTab: NewEditQueueTab(manager, &models.QueueInfo{}),
		editingQueue: false,
		help:         help,
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
		footerString: "",
	}

	queuesTab.updateRows()

	return queuesTab
}

func (m QueuesTab) Init() tea.Cmd { return nil }

func (m QueuesTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.updateRows()

	if m.editingQueue {
		switch msg := msg.(type) {
		case CloseChildMsg:
			m.editingQueue = false
			m.footerString = ""
			m.updateRows()
			return m, nil
		default:
			m.editQueueTab, cmd = m.editQueueTab.Update(msg)
			return m, cmd
		}
	} else if m.addingQueue {
		switch msg := msg.(type) {
		case CloseChildMsg:
			m.addingQueue = false
			m.footerString = ""
			m.updateRows()
			return m, nil
		default:
			m.addQueueTab, cmd = m.addQueueTab.Update(msg)
			return m, cmd
		}
	} else {
		switch msg := msg.(type) {
		// case tea.WindowSizeMsg:
		// 	m.table.SetHeight(msg.Height - 3)
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Navigation):
			case key.Matches(msg, m.keys.Delete):
				if m.table.Cursor() >= 0 && m.table.Cursor() < len(m.queues) {
					m.manager.RemoveQueue(m.queues[m.table.Cursor()].Name)
					m.updateRows()
				}
			case key.Matches(msg, m.keys.NewQueue):
				m.addingQueue = true
				cmd = m.addQueueTab.Init()
			case key.Matches(msg, m.keys.Edit):
				if m.table.Cursor() >= 0 && m.table.Cursor() < len(m.queues) {
					m.editingQueue = true
					q := m.queues[m.table.Cursor()]
					m.editQueueTab = NewEditQueueTab(m.manager, q)
					cmd = m.editQueueTab.Init()
				}
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			}
		}
		var cmds tea.Cmd

		if m.table.Cursor() < 0 || m.table.Cursor() >= len(m.queues) {
			m.table.SetCursor(0)
		}

		m.table, cmds = m.table.Update(msg)
		return m, tea.Batch(cmd, cmds)
	}
}

func (m QueuesTab) View() string {
	if m.editingQueue {
		return m.editQueueTab.View()
	} else if m.addingQueue {
		return m.addQueueTab.View()
	} else {
		if len(m.queues) == 0 {
			m.keys.Delete.SetEnabled(false)
			m.keys.Edit.SetEnabled(false)
		} else {
			m.keys.Delete.SetEnabled(true)
			m.keys.Edit.SetEnabled(true)
		}

		return lipgloss.JoinVertical(
			lipgloss.Left,
			borderedStyle.Render(m.table.View()),
			noStyle.Render(m.footerString),
			helpStyle.Render(m.help.View(m.keys)),
		)
	}
}

func (m *QueuesTab) updateRows() {
	m.queues = m.manager.GetQueueList()
	rows := []table.Row{}
	for _, queue := range m.queues {
		var sp string
		if queue.SpeedLimit == 0 {
			sp = "∞"
		} else {
			sp = speedString(float64(queue.SpeedLimit))
		}
		rows = append(rows, []string{
			queue.Name,
			queue.TargetDirectory,
			fmt.Sprintf("%d", queue.MaxParallel),
			sp,
			queue.StartTime.Format("15:04"),
			queue.EndTime.Format("15:04"),
		})
	}

	m.table.SetRows(rows)
}
