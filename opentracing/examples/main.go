package main

// 对于每一个请求，服务端都生成新的 span

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zamopentracing "github.com/zerogo-hub/zero-api-middleware/opentracing"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := zeroapi.New()

	a.Use(zamopentracing.New("test-1", "0.0.0.0:5778", a))

	a.Get("/", helloworldHandle)

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}
