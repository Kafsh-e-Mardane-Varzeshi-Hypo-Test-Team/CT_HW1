package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

// Key Bindings
type keyMap struct {
	Next       key.Binding
	Prev       key.Binding
	Select     key.Binding
	Submit     key.Binding
	Cancel     key.Binding
	ToggleHelp key.Binding
	Quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleHelp, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev},               // Navigation keys
		{k.Select, k.Submit, k.Cancel}, // Actions
		{k.ToggleHelp, k.Quit},         // Help and quit
	}
}

var keys = keyMap{
	Next: key.NewBinding(
		key.WithKeys("tab", "down"),
		key.WithHelp("tab/↓", "next field"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab/↑", "previous field"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm download"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

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

// fields
type AddDownloadTabField int

const (
	urlField AddDownloadTabField = iota
	filenameField
	queueField
	confirmField
	cancelField
)

// AddDownloadTab Model
type AddDownloadTab struct {
	focusIndex    AddDownloadTabField
	urlInput      textinput.Model
	filenameInput textinput.Model
	queues        list.Model
	choices       []string
	selectedQueue int
	listExpanded  bool
	help          help.Model
	keys          keyMap
	showHelp      bool
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

	help := help.New()
	help.ShowAll = true

	return AddDownloadTab{
		urlInput:      urlInput,
		filenameInput: filenameInput,
		queues:        queueList,
		choices:       availableQueues,
		selectedQueue: 0,
		focusIndex:    0,
		help:          help,
		keys:          keys,
		showHelp:      true,
	}
}

func (m AddDownloadTab) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddDownloadTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ToggleHelp):
			m.showHelp = !m.showHelp
		}
	}
	switch m.focusIndex {
	case urlField:
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
	case filenameField:
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
	case queueField:
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
	case confirmField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				url := m.urlInput.Value()
				filename := m.filenameInput.Value()
				queue := m.choices[m.selectedQueue]
				id, err := models.AddDownload(url, filename, queue)
				if err == nil {
					fmt.Println("Download added with ID:", id)
				} else {
					fmt.Println("Error adding download:", err)
				}
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 4)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	case cancelField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				return NewMainView(), nil
			case "tab", "down":
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
	case urlField:
		m.urlInput.Focus()
		m.urlInput.PromptStyle = focusedStyle
		m.urlInput.TextStyle = focusedStyle
	case filenameField:
		m.filenameInput.Focus()
		m.filenameInput.PromptStyle = focusedStyle
		m.filenameInput.TextStyle = focusedStyle
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
		helpStyle.Render(m.help.View(m.keys)),
	)

	return docStyle.Render(form)
}
