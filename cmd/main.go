package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	logger "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/logger"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/models"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/persistence"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/tui"
)

const filename string = "internal/persistence/data.json"

func main() {
	go logger.StartLoggingToFile()
	manager, err := persistence.Load(filename)
	if err != nil {
		log.Fatalln(err)
	}
	saveState(manager)
	manager.Start()

	autoSave := func() {
		ticker := time.NewTicker(time.Duration(30 * time.Second))
		defer ticker.Stop()

		for range ticker.C {
			saveState(manager)
		}
	}
	go autoSave()

	p := tea.NewProgram(tui.NewMainView(manager))
	if _, err := p.Run(); err != nil {
		panic(err)
	}

	manager.Stop()
	saveState(manager)
}

func saveState(manager *models.Manager) error {
	jsonData, err := manager.GetJson()
	if err != nil {
		log.Println(err)
	}

	persistence.Save(filename, jsonData)
	return nil
}
