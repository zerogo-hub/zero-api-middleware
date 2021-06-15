package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zambasic "github.com/zerogo-hub/zero-api-middleware/auth/basic"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)

	// 预置一对账号密码
	a.Use(zambasic.New(nil, map[string]string{"foo": "bar"}))

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}
