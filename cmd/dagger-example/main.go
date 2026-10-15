package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// 初始化 Dagger 客户端
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// 获取当前目录
	src := client.Host().Directory(".")

	// 定义容器并在其中编译 Go 项目（工作目录用 /src，避免落在镜像默认的 $GOPATH /go 下）
	output, err := client.Container().
		From("golang:1.27").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "./..."}). // 运行单元测试
		Stdout(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
