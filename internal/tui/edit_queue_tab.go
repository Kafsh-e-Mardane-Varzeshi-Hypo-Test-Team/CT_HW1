package tui

import (
	"fmt"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type editQueueKeyMap struct {
	Next       key.Binding
	Prev       key.Binding
	Navigation key.Binding
	Select     key.Binding
	Cancel     key.Binding
}

func (k editQueueKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Cancel}
}

func (k editQueueKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Navigation},
		{k.Select, k.Cancel},
	}
}

type EditQueueField int

const (
	editTargetDirectoryField EditQueueField = iota
	editMaxParallelField
	editSpeedLimitField
	editStartTimeField
	editEndTimeField
	editConfirmQueueField
	editCancelQueueField
)

type EditQueueTab struct {
	manager        *models.Manager
	queueName      string
	focusIndex     EditQueueField
	targetDirInput textinput.Model
	maxParallel    textinput.Model
	speedLimit     textinput.Model
	startTime      textinput.Model
	endTime        textinput.Model
	help           help.Model
	keys           editQueueKeyMap
	footerMessage  string
}

func NewEditQueueTab(manager *models.Manager, queueInfo *models.QueueInfo) EditQueueTab {
	name := queueInfo.Name

	targetDirInput := textinput.New()
	targetDirInput.Placeholder = "Enter target directory"
	targetDirInput.SetValue(queueInfo.TargetDirectory)
	targetDirInput.PromptStyle = noStyle
	targetDirInput.TextStyle = noStyle
	targetDirInput.Cursor.Style = cursorStyle

	maxParallel := textinput.New()
	maxParallel.Placeholder = "Enter max parallel downloads (integer)"
	maxParallel.SetValue(fmt.Sprint(queueInfo.MaxParallel))
	maxParallel.PromptStyle = noStyle
	maxParallel.TextStyle = noStyle
	maxParallel.Cursor.Style = cursorStyle

	speedLimit := textinput.New()
	speedLimit.Placeholder = "Enter speed limit (Bytes per second) (0 for no limit)"
	speedLimit.SetValue(fmt.Sprint(queueInfo.SpeedLimit))
	speedLimit.PromptStyle = noStyle
	speedLimit.TextStyle = noStyle
	speedLimit.Cursor.Style = cursorStyle

	startTime := textinput.New()
	startTime.Placeholder = "Enter start time (HH:MM)"
	startTime.SetValue(queueInfo.StartTime.Format("15:04"))
	startTime.PromptStyle = noStyle
	startTime.TextStyle = noStyle
	startTime.Cursor.Style = cursorStyle

	endTime := textinput.New()
	endTime.Placeholder = "Enter end time (HH:MM)"
	endTime.SetValue(queueInfo.EndTime.Format("15:04"))
	endTime.PromptStyle = noStyle
	endTime.TextStyle = noStyle
	endTime.Cursor.Style = cursorStyle

	help := help.New()
	help.ShowAll = true
	help.FullSeparator = " \t "

	return EditQueueTab{
		manager:        manager,
		focusIndex:     editTargetDirectoryField,
		queueName:      name,
		targetDirInput: targetDirInput,
		maxParallel:    maxParallel,
		speedLimit:     speedLimit,
		startTime:      startTime,
		endTime:        endTime,
		help:           help,
		keys: editQueueKeyMap{
			Next:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
			Prev:       key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "previous field")),
			Navigation: key.NewBinding(key.WithKeys("up", "down", "left", "right"), key.WithHelp("↑/↓/←/→", "navigate")),
			Select:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			Cancel:     key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("ctrl+c/esc", "cancel")),
		},
		footerMessage: "",
	}
}

func (m EditQueueTab) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditQueueTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case editTargetDirectoryField:
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
	case editMaxParallelField:
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
	case editSpeedLimitField:
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
	case editStartTimeField:
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
	case editEndTimeField:
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
	case editConfirmQueueField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				name := m.queueName
				targetDir := m.targetDirInput.Value()
				maxParallel := m.maxParallel.Value()
				speedLimit := m.speedLimit.Value()
				startTime := m.startTime.Value()
				endTime := m.endTime.Value()

				queueInfo, err := makeQueueInfo(name, targetDir, maxParallel, speedLimit, startTime, endTime)

				if err != nil {
					m.footerMessage = err.Error()
					return m, nil
				}

				err = m.manager.UpdateQueue(queueInfo)
				if err != nil {
					m.footerMessage = err.Error()
					return m, nil
				}
				m.footerMessage = "Queue updated successfully."
				return m, func() tea.Msg { return CloseChildMsg{} }
			case "up", "shift+tab":
				m.focusIndex = max(m.focusIndex-1, 0)
			case "ctrl+c", "esc":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			case "down":
				m.focusIndex = editConfirmQueueField
			case "tab", "right":
				m.focusIndex = editCancelQueueField
				cmd = tea.Cmd(textinput.Blink)
			}
		}
	case editCancelQueueField:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.resetForm()
				return m, func() tea.Msg { return CloseChildMsg{} }
			case "up":
				m.focusIndex = editEndTimeField
			case "left", "shift+tab":
				m.focusIndex = editConfirmQueueField
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

func (m *EditQueueTab) updateFocus() {
	m.targetDirInput.Blur()
	m.maxParallel.Blur()
	m.speedLimit.Blur()
	m.startTime.Blur()
	m.endTime.Blur()

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
	case editTargetDirectoryField:
		m.targetDirInput.Focus()
		m.targetDirInput.PromptStyle = focusedStyle
		m.targetDirInput.TextStyle = focusedStyle
	case editMaxParallelField:
		m.maxParallel.Focus()
		m.maxParallel.PromptStyle = focusedStyle
		m.maxParallel.TextStyle = focusedStyle
	case editSpeedLimitField:
		m.speedLimit.Focus()
		m.speedLimit.PromptStyle = focusedStyle
		m.speedLimit.TextStyle = focusedStyle
	case editStartTimeField:
		m.startTime.Focus()
		m.startTime.PromptStyle = focusedStyle
		m.startTime.TextStyle = focusedStyle
	case editEndTimeField:
		m.endTime.Focus()
		m.endTime.PromptStyle = focusedStyle
		m.endTime.TextStyle = focusedStyle
	}
}

func (m EditQueueTab) View() string {
	var speedLimit string

	blurredConfirm := blurredStyle.Render("[ Confirm ]")
	blurredCancel := blurredStyle.Render("[ Cancel ]")

	if m.focusIndex == editConfirmQueueField {
		blurredConfirm = focusedStyle.Render("[ Confirm ]")
	} else if m.focusIndex == editCancelQueueField {
		blurredCancel = focusedStyle.Render("[ Cancel ]")
	}

	// speedLimit, if empty or value is 0, add no limit
	if m.speedLimit.Value() == "0" {
		speedLimit = m.speedLimit.View() + " (no limit)"
	} else {
		speedLimit = m.speedLimit.Value()
		if speedLimit == "" {
			speedLimit = m.speedLimit.View()
		} else {
			sp, err := strconv.ParseInt(speedLimit, 10, 64)
			if err != nil {
				speedLimit = fmt.Sprint(m.speedLimit.View(), " (invalid)")
			} else {
				speedLimit = fmt.Sprint(m.speedLimit.View(), "B/s (", speedString(float64(sp)), ")")
			}
		}
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
						m.queueName,
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
						speedLimit,
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

func (m *EditQueueTab) resetForm() {
	m.targetDirInput.SetValue("")
	m.maxParallel.SetValue("")
	m.speedLimit.SetValue("")
	m.startTime.SetValue("")
	m.endTime.SetValue("")
	m.focusIndex = 0
	m.footerMessage = ""
}
