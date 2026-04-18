package main

import (
	"fmt"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func main() {
	logger, err := fluent.New(fluent.Config{
		FluentPort: 24224,
		FluentHost: "localhost",
	})
	if err != nil {
		fmt.Println(err)
	}
	defer logger.Close()
	tag := "myapp.access"
	data := map[string]string{
		"foo":  "bar",
		"hoge": "hoge",
	}
	err = logger.Post(tag, data)
	if err != nil {
		panic(err)
	}

	err = logger.PostWithTime(tag, time.Now(), data)
	if err != nil {
		panic(err)
	}
}
