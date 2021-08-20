package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zamlogger "github.com/zerogo-hub/zero-api-middleware/logger"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func hellopanic(ctx zeroapi.Context) {
	// 即使发生异常，也能正常打印请求日志
	panic("hello panic")
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)

	a.Get("/panic", hellopanic)

	a.Use(zamlogger.New())

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Errorf("app run failed, err: %s", err.Error())
	}
}
