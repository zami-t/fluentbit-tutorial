package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type MyFluentLogger struct {
	*fluent.Fluent
}

func NewMyFluentLogger() (*MyFluentLogger, error) {
	logger, _ := fluent.New(fluent.Config{
		FluentHost:          "fluent-bit-sidecar",
		FluentPort:          24224,
		Timeout:             3 * time.Second,
		ReadTimeout:         3 * time.Second,
		WriteTimeout:        15 * time.Second,
		BufferLimit:         8192,
		MaxRetry:            13,
		MaxRetryWait:        60000,
		Async:               true,
		ForceStopAsyncSend:  false,
		AsyncResultCallback: asyncResultCallback,
		RequestAck:          true, // Fluent Bit の Require_ack_response と対になる
	})

	return &MyFluentLogger{logger}, nil
}

func (l *MyFluentLogger) Close() {
	l.Fluent.Close()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	logger, err := NewMyFluentLogger()
	if err != nil {
		println("Failed to create logger:", err.Error())
		return
	}
	defer logger.Close()

LOOP:
	for {
		select {
		case <-time.After(5 * time.Second):
			println("send message")
			logger.sendMessage("app.log", map[string]string{"message": "Hello, Fluent Bit!"})
		case <-c:
			println("ctl+c pressed")
			break LOOP
		}
	}
}

func (l *MyFluentLogger) sendMessage(tag string, message map[string]string) {
	l.Post(tag, message)
}

func asyncResultCallback(data []byte, err error) {
	if err != nil {
		println("Failed to send log", "data: ", string(data), "error: ", err.Error())
	} else {
		println("Successfully sent log:", string(data))
	}
}

