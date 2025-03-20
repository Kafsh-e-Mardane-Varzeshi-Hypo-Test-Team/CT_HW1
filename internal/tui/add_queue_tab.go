// - AddQueueTab
package tui

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Key Bindings
type addQueueKeyMap struct {
	Next       key.Binding
	Prev       key.Binding
	Navigation key.Binding
	Select     key.Binding
	Cancel     key.Binding
}

func (k addQueueKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Cancel}
}

func (k addQueueKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Navigation}, // Navigation keys
		{k.Select, k.Cancel},           // Actions
	}
}

type addQueueTabField int

const (
	nameField addQueueTabField = iota
	targetDirectoryField
	maxParallelField
	speedLimitField
	startTimeField
	endTimeField
	confirmAddQueueField
	cancelAddQueueField
)

// AddQueueTab Model
type AddQueueTab struct {
	manager        *models.Manager
	focusIndex     addQueueTabField
	nameInput      textinput.Model
	targetDirInput textinput.Model
	maxParallel    textinput.Model
	speedLimit     textinput.Model
	startTime      textinput.Model
	endTime        textinput.Model
	help           help.Model
	keys           addQueueKeyMap
	footerMessage  string
}

func NewAddQueueTab(manager *models.Manager) AddQueueTab {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter queue name"
	nameInput.Focus()
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle
	nameInput.Cursor.Style = cursorStyle

	targetDirInput := textinput.New()
	targetDirInput.Placeholder = "Enter target directory"
	targetDirInput.PromptStyle = noStyle
	targetDirInput.TextStyle = noStyle
	targetDirInput.Cursor.Style = cursorStyle

	maxParallel := textinput.New()
	maxParallel.Placeholder = "Enter max parallel downloads"
	maxParallel.PromptStyle = noStyle
	maxParallel.TextStyle = noStyle
	maxParallel.Cursor.Style = cursorStyle

	speedLimit := textinput.New()
	speedLimit.Placeholder = "Enter speed limit"
	speedLimit.PromptStyle = noStyle
	speedLimit.TextStyle = noStyle
	speedLimit.Cursor.Style = cursorStyle

	startTime := textinput.New()
	startTime.Placeholder = "Enter start time"
	startTime.PromptStyle = noStyle
	startTime.TextStyle = noStyle
	startTime.Cursor.Style = cursorStyle

	endTime := textinput.New()
	endTime.Placeholder = "Enter end time"
	endTime.PromptStyle = noStyle
	endTime.TextStyle = noStyle
	endTime.Cursor.Style = cursorStyle

	help := help.New()
	help.ShowAll = true
	help.FullSeparator = " \t "

	return AddQueueTab{
		manager:        manager,
		nameInput:      nameInput,
		targetDirInput: targetDirInput,
		maxParallel:    maxParallel,
		speedLimit:     speedLimit,
		startTime:      startTime,
		endTime:        endTime,
		help:           help,
		keys: addQueueKeyMap{
			Next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next field"),
			),
			Prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "previous field"),
			),
			Navigation: key.NewBinding(
				key.WithKeys("up", "down", "left", "right"),
				key.WithHelp("↑/↓/←/→", "navigate"),
			),
			Select: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			Cancel: key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("ctrl+c/esc", "cancel"),
			),
		},
		footerMessage: "",
	}
}

func (m AddQueueTab) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddQueueTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case nameField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 7)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
		m.nameInput, cmd = m.nameInput.Update(msg)
	case targetDirectoryField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 7)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
		m.targetDirInput, cmd = m.targetDirInput.Update(msg)
	case maxParallelField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 7)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
		m.maxParallel, cmd = m.maxParallel.Update(msg)
	case speedLimitField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 7)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
		m.speedLimit, cmd = m.speedLimit.Update(msg)
	case startTimeField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 7)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
		m.startTime, cmd = m.startTime.Update(msg)
	case endTimeField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab", "down":
				m.focusIndex = min(m.focusIndex+1, 7)
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
		m.endTime, cmd = m.endTime.Update(msg)
	case confirmAddQueueField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// name := m.nameInput.Value()
				// targetDir := m.targetDirInput.Value()
				// maxParallel := m.maxParallel.Value()
				// speedLimit := m.speedLimit.Value()
				// startTime := m.startTime.Value()
				// endTime := m.endTime.Value()

				// var err error
				// if name == "" {
				// 	err = fmt.Errorf("Name cannot be empty")
				// } else {
				// 	err = m.manager.AddQueue(name, targetDir, maxParallel, speedLimit, startTime, endTime)
				// }

				// if err == nil {
				// 	m.footerMessage = "Queue added successfully."
				// 	m.nameInput.SetValue("")
				// 	m.targetDirInput.SetValue("")
				// 	m.maxParallel.SetValue("")
				// 	m.speedLimit.SetValue("")
				// 	m.startTime.SetValue("")
				// 	m.endTime.SetValue("")
				// 	m.focusIndex = 0
				// } else {
				// 	m.footerMessage = "Error adding queue:" + err.Error()
				// }
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			case "down":
				m.focusIndex = confirmAddQueueField
			case "tab", "right":
				m.focusIndex = cancelAddQueueField
				cmd = tea.Cmd(textinput.Blink)
			}
		}
	case cancelAddQueueField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			case "up":
				m.focusIndex = endTimeField
			case "left", "shift+tab":
				m.focusIndex = confirmAddQueueField
				cmd = tea.Cmd(textinput.Blink)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			}
		}
	}
	m.updateFocus()

	return m, cmd
}

func (m *AddQueueTab) updateFocus() {
	m.nameInput.Blur()
	m.targetDirInput.Blur()
	m.maxParallel.Blur()
	m.speedLimit.Blur()
	m.startTime.Blur()
	m.endTime.Blur()

	m.nameInput.PromptStyle = noStyle
	m.nameInput.TextStyle = noStyle
	m.targetDirInput.PromptStyle = noStyle
	m.targetDirInput.TextStyle = noStyle
	m.maxParallel.PromptStyle = noStyle
	m.maxParallel.TextStyle = noStyle
	m.speedLimit.PromptStyle = noStyle
	m.speedLimit.TextStyle = noStyle
	m.startTime.PromptStyle = noStyle
	m.startTime.TextStyle = noStyle
	m.endTime.PromptStyle = noStyle
	m.endTime.TextStyle = noStyle

	switch m.focusIndex {
	case nameField:
		m.nameInput.Focus()
		m.nameInput.PromptStyle = focusedStyle
		m.nameInput.TextStyle = focusedStyle
	case targetDirectoryField:
		m.targetDirInput.Focus()
		m.targetDirInput.PromptStyle = focusedStyle
		m.targetDirInput.TextStyle = focusedStyle
	case maxParallelField:
		m.maxParallel.Focus()
		m.maxParallel.PromptStyle = focusedStyle
		m.maxParallel.TextStyle = focusedStyle
	case speedLimitField:
		m.speedLimit.Focus()
		m.speedLimit.PromptStyle = focusedStyle
		m.speedLimit.TextStyle = focusedStyle
	case startTimeField:
		m.startTime.Focus()
		m.startTime.PromptStyle = focusedStyle
		m.startTime.TextStyle = focusedStyle
	case endTimeField:
		m.endTime.Focus()
		m.endTime.PromptStyle = focusedStyle
		m.endTime.TextStyle = focusedStyle
	}
}

func (m AddQueueTab) View() string {
	blurredConfirm := blurredStyle.Render("[ Confirm ]")
	blurredCancel := blurredStyle.Render("[ Cancel ]")

	if m.focusIndex == confirmAddQueueField {
		blurredConfirm = focusedStyle.Render("[ Confirm ]")
	} else if m.focusIndex == cancelAddQueueField {
		blurredCancel = focusedStyle.Render("[ Cancel ]")
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
						noStyle.Render("Name: "),
						m.nameInput.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("Target Directory: "),
						m.targetDirInput.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("Max Parallel Downloads: "),
						m.maxParallel.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("Speed Limit: "),
						m.speedLimit.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("Start Time: "),
						m.startTime.View(),
					),
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						noStyle.Render("End Time: "),
						m.endTime.View(),
					),
				),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					blurredConfirm,
					blurredCancel,
				),
			),
		),

		noStyle.Render(m.footerMessage),
		helpStyle.Render(m.help.View(m.keys)),
	)

	return docStyle.Render(form)
}

func (m *AddQueueTab) resetForm() {
	m.nameInput.SetValue("")
	m.targetDirInput.SetValue("")
	m.maxParallel.SetValue("")
	m.speedLimit.SetValue("")
	m.startTime.SetValue("")
	m.endTime.SetValue("")
	m.focusIndex = 0
	m.footerMessage = ""
}

type CloseChildMsg struct{}
