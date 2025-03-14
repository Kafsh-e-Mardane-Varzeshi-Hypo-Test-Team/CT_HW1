package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	cursorStyle       = focusedStyle
	noStyle           = lipgloss.NewStyle()
	helpStyle         = blurredStyle.Margin(1, 0, 0, 0)

	focusedConfirm = focusedStyle.Render("[ Confirm ]")
	blurredConfirm = blurredStyle.Render("[ Confirm ]")
	focusedCancel  = focusedStyle.Render("[ Cancel ]")
	blurredCancel  = blurredStyle.Render("[ Cancel ]")
)

// List Item Delegate
type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// AddDownloadTab Model
type AddDownloadTab struct {
	focusIndex    int
	urlInput      textinput.Model
	filenameInput textinput.Model
	queues        list.Model
	choices       []string
	selectedQueue int
	listExpanded  bool
}

func NewAddDownloadTab() AddDownloadTab {
	urlInput := textinput.New()
	urlInput.Placeholder = "Enter file URL"
	urlInput.Focus()
	urlInput.PromptStyle = focusedStyle
	urlInput.TextStyle = focusedStyle
	urlInput.Cursor.Style = cursorStyle

	filenameInput := textinput.New()
	filenameInput.Placeholder = "(Optional) Enter output filename"
	filenameInput.PromptStyle = noStyle
	filenameInput.TextStyle = noStyle
	filenameInput.Cursor.Style = cursorStyle

	queues := models.GetQueues()
	availableQueues := []string{}
	for _, q := range queues {
		availableQueues = append(availableQueues, q.Name)
	}
	items := []list.Item{}
	for _, q := range availableQueues {
		items = append(items, item(q))
	}

	queueList := list.New(items, itemDelegate{}, 30, min(3, len(availableQueues)))
	queueList.SetShowTitle(false)
	queueList.SetShowStatusBar(false)
	queueList.SetFilteringEnabled(false)
	queueList.SetHeight(8)

	return AddDownloadTab{
		urlInput:      urlInput,
		filenameInput: filenameInput,
		queues:        queueList,
		choices:       availableQueues,
		selectedQueue: 0,
		focusIndex:    0,
	}
}

func (m AddDownloadTab) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddDownloadTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "enter", "down":
				m.focusIndex = min(m.focusIndex+1, 4)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c":
				return m, tea.Quit
			default:
				m.urlInput, cmd = m.urlInput.Update(msg)
			}
		}
	case 1:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "enter", "down":
				m.focusIndex = min(m.focusIndex+1, 4)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c":
				return m, tea.Quit
			default:
				m.filenameInput, cmd = m.filenameInput.Update(msg)
			}
		}
	case 2:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if m.listExpanded {
				switch msg.String() {
				case "enter":
					m.selectedQueue = m.queues.Index()
					m.listExpanded = false
					return m, nil
				case "esc":
					m.listExpanded = false
				default:
					m.queues, cmd = m.queues.Update(msg)
					return m, cmd
				}
			} else {
				switch msg.String() {
				case "tab", "down":
					m.focusIndex = min(m.focusIndex+1, 4)
				case "up", "shift+tab":
					m.focusIndex = max(m.focusIndex-1, 0)
				case "enter":
					m.listExpanded = true
					return m, nil
				case "ctrl+c":
					return m, tea.Quit
				}
			}
		}
	case 3, 4:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "enter", "down":
				m.focusIndex = min(m.focusIndex+1, 4)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}
	m.updateFocus()

	return m, cmd
}

func (m *AddDownloadTab) updateFocus() {
	m.urlInput.Blur()
	m.filenameInput.Blur()
	m.urlInput.PromptStyle = noStyle
	m.urlInput.TextStyle = noStyle
	m.filenameInput.PromptStyle = noStyle
	m.filenameInput.TextStyle = noStyle

	switch m.focusIndex {
	case 0:
		m.urlInput.Focus()
		m.urlInput.PromptStyle = focusedStyle
		m.urlInput.TextStyle = focusedStyle
	case 1:
		m.filenameInput.Focus()
		m.filenameInput.PromptStyle = focusedStyle
		m.filenameInput.TextStyle = focusedStyle
	case 2:
		// Blinking effect for queue selection
	}
}

func (m AddDownloadTab) View() string {
	var queueDisplay string
	if m.listExpanded {
		queueDisplay = m.queues.View()
	} else {
		if m.focusIndex == 2 {
			queueDisplay = focusedStyle.Render(m.choices[m.selectedQueue])
		} else {
			queueDisplay = blurredStyle.Render(m.choices[m.selectedQueue])
		}
	}

	buttonConfirm := blurredConfirm
	buttonCancel := blurredCancel
	if m.focusIndex == 3 {
		buttonConfirm = focusedConfirm
	} else if m.focusIndex == 4 {
		buttonCancel = focusedCancel
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			"URL: ",
			m.urlInput.View(),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			"Filename: ",
			m.filenameInput.View(),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			"Destination Queue: ",
			queueDisplay,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			buttonConfirm,
			buttonCancel,
		),
	)

	return docStyle.Render(form)
}

// Utility function to limit queue list size
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
