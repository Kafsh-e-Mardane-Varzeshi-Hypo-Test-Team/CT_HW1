package log

import (
	"log"
	"os"
)

func StartLoggingToFile() {
	file, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)
}
