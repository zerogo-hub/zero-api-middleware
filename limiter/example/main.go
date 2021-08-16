package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	app "github.com/zerogo-hub/zero-api/app"

	zamtime "github.com/zerogo-hub/zero-helper/time"

	zamlimiter "github.com/zerogo-hub/zero-api-middleware/limiter"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)

	// 10 秒内允许 2 个请求
	a.Use(zamlimiter.New(zamtime.Second(10), 2))

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}

// 默认有 timestamp, nonce, sign 三个必要字段
// 请求: curl -i -X GET http://127.0.0.1:8877
