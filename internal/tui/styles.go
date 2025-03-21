package tui

import "github.com/charmbracelet/lipgloss"

// General styles
var (
	borderedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	noStyle = lipgloss.NewStyle()
)

// Form styles
var (
	focusedStyle = noStyle.Foreground(lipgloss.Color("205"))
	blurredStyle = noStyle.Foreground(lipgloss.Color("240"))

	selectedItemStyle = focusedStyle.PaddingLeft(2)
	itemStyle         = noStyle.PaddingLeft(4)
	cursorStyle       = focusedStyle

	helpStyle = blurredStyle.Margin(1, 0, 0, 0)

	blurredButton  = noStyle.Foreground(lipgloss.Color("231"))
	focusedConfirm = focusedStyle.Render("[ Confirm ]")
	blurredConfirm = blurredButton.Render("[ Confirm ]")
	focusedCancel  = focusedStyle.Render("[ Cancel ]")
	blurredCancel  = blurredButton.Render("[ Cancel ]")
)

// Tab item styles
var (
	tabBorder        = lipgloss.RoundedBorder()
	docStyle         = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().
				Border(tabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 2).
				AlignHorizontal(lipgloss.Center)
	activeTabStyle = inactiveTabStyle.
			Foreground(highlightColor).
			Padding(0, 2).
			AlignHorizontal(lipgloss.Center)
	// Background(highlightColor).
	// BorderBackground(highlightColor).
)
