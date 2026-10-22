package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed go.mod
var content []byte

// Version 传入 go.mod 内容 解析依赖版本
func Version(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "go "); ok {
			return strings.TrimSpace(after)
		}
	}
	return ""
}

func main() {
	fmt.Println(Version(content))
	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() {
			time.Sleep(time.Second)
			fmt.Println(v)
			done <- true
		}()
	}

	// wait for all goroutines to complete before exiting
	for range values {
		<-done
	}
}
