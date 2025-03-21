package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal"
)

func main() {
	manager, err := LoadManager("hello.txt")
	if err != nil {
		log.Fatalln(err)
	}
	manager.Start(make(chan struct{}))

}

func LoadManager(filename string) (*internal.Manager, error) {
	file, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return internal.NewManager(), nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	m := &internal.Manager{}
	if err := decoder.Decode(m); err != nil {
		return nil, err
	}

	return m, nil
}
