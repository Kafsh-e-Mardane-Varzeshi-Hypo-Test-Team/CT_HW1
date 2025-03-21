package main

import (
	"log"
	"time"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal/persistence"
)

const filename string = "data.json"

func main() {
	manager, err := persistence.Load(filename)
	if err != nil {
		log.Fatalln(err)
	}
	manager.Start()

	go autoSave(manager)
}

func autoSave(manager *internal.Manager) error {
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
