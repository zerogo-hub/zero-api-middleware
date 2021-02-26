package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zmlogger "github.com/zerogo-hub/zero-api-middleware/logger"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func hellopanic(ctx zeroapi.Context) {
	panic("hello panic")
}

func main() {
	a := zeroapi.New()

	a.Get("/", helloworldHandle)

	a.Get("/panic", hellopanic)

	a.Use(zmlogger.New())

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}
