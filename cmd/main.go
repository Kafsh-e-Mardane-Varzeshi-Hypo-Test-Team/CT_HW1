package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/persistence"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/tui"
)

const filename string = "data.json"

func main() {
	manager, err := persistence.Load(filename)
	if err != nil {
		log.Fatalln(err)
	}
	manager.Start()

	go autoSave(manager)

	p := tea.NewProgram(tui.NewMainView(manager))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func autoSave(manager *models.Manager) error {
	ticker := time.NewTicker(time.Duration(30 * time.Second))
	defer ticker.Stop()

	for range ticker.C {
		jsonData, err := manager.GetJson()
		if err != nil {
			return err
		}

		persistence.Save(filename, jsonData)
	}
	return nil
}
