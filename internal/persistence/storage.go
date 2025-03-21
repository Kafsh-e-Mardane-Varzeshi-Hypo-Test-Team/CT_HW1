package persistence

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/internal"
)

func Load(filename string) (*internal.Manager, error) {
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

func Save(filename string, jsonData []byte) error {
	err := os.WriteFile(filename, jsonData, 0644) // 0644 is the file permission (read/write for owner, read for others)
	if err != nil {
		return err
	}
	return nil
}
