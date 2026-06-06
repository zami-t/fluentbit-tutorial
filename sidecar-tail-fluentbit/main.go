package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	logDir  = "/var/log/app"
	logFile = "app.log"
)

type LogEntry struct {
	Message string `json:"message"`
}

func main() {
	// ログディレクトリを用意（ボリュームマウント済み想定だが念のため）
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.Fatalf("failed to create log dir: %v", err)
	}

	f, err := os.OpenFile(
		filepath.Join(logDir, logFile),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()

	// シグナルハンドリング（SIGTERMはdocker stop時に送られる）
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("app started, writing logs to", filepath.Join(logDir, logFile))

LOOP:
	for {
		select {
		case <-ticker.C:
			entry := LogEntry{
				Message: "Hello, Fluent Bit!",
			}
			if err := writeJSONLine(f, entry); err != nil {
				log.Printf("failed to write log: %v", err)
				continue
			}
			log.Println("wrote log entry")
		case sig := <-sigCh:
			log.Printf("signal received: %v, shutting down", sig)
			break LOOP
		}
	}
}

// JSON Lines形式で1行書き込み + fsync
func writeJSONLine(f *os.File, entry LogEntry) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if _, err := f.Write(b); err != nil {
		return err
	}
	// tail側が即座に読めるように flush。負荷が高い場合は外す判断もあり。
	return f.Sync()
}

