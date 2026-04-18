// app/main.go
package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
}

func main() {
	logDir := "/var/log/app"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logDir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	for {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "info",
			Message:   "hello from go app",
			Service:   "my-service",
		}
		if err := encoder.Encode(entry); err != nil {
			log.Println("write error:", err)
		}
		time.Sleep(3 * time.Second)
	}
}
