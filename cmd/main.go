package main

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	manager := models.NewManager()
	manager.Queues = map[string]*models.Queue{
		"Queue1": {
			Name: "Queue1",
		},
	}
	manager.Downloads = []*models.Download{
		{
			URL:          "https://speed.hetzner.de/100MB.bin",
			Queue:        "Queue1",
			Status:       models.InProgress,
			TransferRate: 123,
			Progress:     0.234,
		},
		{
			URL:          "https://speed.hetzner.de/100MB.bin",
			Queue:        "Queue1",
			Status:       models.InProgress,
			TransferRate: 123123,
			Progress:     0.134,
		},
		{
			URL:          "https://speed.hetzner.de/100MB.bin",
			Queue:        "Queue1",
			Status:       models.InProgress,
			TransferRate: 123123123,
			Progress:     0.334,
		},
		{
			URL:          "https://speed.hetzner.de/100MB.bin",
			Queue:        "Queue1",
			Status:       models.InProgress,
			TransferRate: 123123123123,
			Progress:     0.634,
		},
	}
	p := tea.NewProgram(tui.NewMainView(manager))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
