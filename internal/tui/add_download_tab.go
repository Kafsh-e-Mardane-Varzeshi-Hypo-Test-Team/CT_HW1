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

// Key Bindings
type addDownloadKeyMap struct {
	next       key.Binding
	prev       key.Binding
	navigation key.Binding
	selectOpt  key.Binding
	quit       key.Binding
}

func (k addDownloadKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k addDownloadKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.next, k.prev, k.navigation}, // Navigation keys
		{k.selectOpt, k.quit},          // Actions
	}
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
type addDownloadTabField int

const (
	urlField addDownloadTabField = iota
	filenameField
	queueField
	confirmDownloadField
	cancelDownloadField
)

// AddDownloadTab Model
type AddDownloadTab struct {
	manager       *models.Manager
	focusIndex    addDownloadTabField
	urlInput      textinput.Model
	filenameInput textinput.Model
	queues        list.Model
	choices       []string
	selectedQueue int
	listExpanded  bool
	help          help.Model
	keys          addDownloadKeyMap
	footerMessage string
}

func NewAddDownloadTab(manager *models.Manager) AddDownloadTab {
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
	help.FullSeparator = " \t "

	return AddDownloadTab{
		manager:       manager,
		urlInput:      urlInput,
		filenameInput: filenameInput,
		queues:        queueList,
		choices:       availableQueues,
		selectedQueue: 0,
		focusIndex:    0,
		help:          help,
		keys: addDownloadKeyMap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next field"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "previous field"),
			),
			navigation: key.NewBinding(
				key.WithKeys("up", "down", "left", "right"),
				key.WithHelp("↑/↓/←/→", "navigate"),
			),
			selectOpt: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("ctrl+c/esc", "quit"),
			),
		},
		footerMessage: "",
	}
}

func (m AddDownloadTab) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddDownloadTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case urlField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
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
			case "tab", "down":
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
	case confirmDownloadField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				url := m.urlInput.Value()
				filename := m.filenameInput.Value()
				queue := m.choices[m.selectedQueue]

				var err error
				if url == "" {
					err = fmt.Errorf("URL cannot be empty")
				} else {
					err = m.manager.AddDownload(url, filename, queue)
				}

				if err == nil {
					m.footerMessage = "Download added successfully."
					m.urlInput.SetValue("")
					m.filenameInput.SetValue("")
					m.selectedQueue = 0
					m.focusIndex = 0
				} else {
					m.footerMessage = "Error adding download:" + err.Error()
				}
			case "tab", "right":
				m.focusIndex = min(m.focusIndex+1, 4)
				cmd = tea.Cmd(textinput.Blink)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	case cancelDownloadField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// reset fields
				m.urlInput.SetValue("")
				m.filenameInput.SetValue("")
				m.selectedQueue = 0
				m.focusIndex = 0
				m.footerMessage = ""
			case "tab", "down":
			case "shift+tab", "left":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "up":
				m.focusIndex = queueField
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

	var footerHelpText string

	if m.listExpanded {
		queueDisplay = m.queues.View()
		footerHelpText = ""
	} else {
		if m.focusIndex == 2 {
			queueDisplay = focusedStyle.Render(m.choices[m.selectedQueue])
		} else {
			queueDisplay = noStyle.Render(m.choices[m.selectedQueue])
		}
		footerHelpText = m.help.View(m.keys)
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
		borderedStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("URL: "),
						m.urlInput.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("Filename: "),
						m.filenameInput.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("Destination Queue: "),
						queueDisplay,
					),
				),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					buttonConfirm,
					buttonCancel,
				),
			),
		),

		noStyle.Render(m.footerMessage),
		helpStyle.Render(footerHelpText),
	)

	return docStyle.Render(form)
}
