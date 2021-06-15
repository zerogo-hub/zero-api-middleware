package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zambodylimit "github.com/zerogo-hub/zero-api-middleware/bodylimit"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := app.New()

	// curl -X POST -d '{"id":123}' http://127.0.0.1:8877
	// {"code":"413","message":"request entity too large"}
	a.Post("/", helloworldHandle)

	a.Use(zambodylimit.New(zeroapi.MethodPost, 1))

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}
