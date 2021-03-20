package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zamthrottle "github.com/zerogo-hub/zero-api-middleware/throttle"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("method: %s, `ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", ctx.Method(), pid, pid)
}

func main() {
	a := zeroapi.New()

	a.Use(zamthrottle.New(&zamthrottle.Config{
		// 每分钟请求两次
		PerMin: 2,
		Burst:  1,
	}))

	a.Get("/", helloworldHandle)

	a.Post("/", helloworldHandle)

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}

// curl -i -X GET http://127.0.0.1:8877
// 第一次请求
// X-Ratelimit-Limit: 2
// X-Ratelimit-Remaining: 1
// X-Ratelimit-Reset: 30
//
// 第二次请求
// X-Ratelimit-Limit: 2
// X-Ratelimit-Remaining: 0
// X-Ratelimit-Reset: 52

// curl -i -X POST http://127.0.0.1:8877
