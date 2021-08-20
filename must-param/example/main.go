package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zammustparam "github.com/zerogo-hub/zero-api-middleware/must-param"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)

	a.Use(zammustparam.New())

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Errorf("app run failed, err: %s", err.Error())
	}
}

// 默认有 timestamp, nonce, sign 三个必要字段
// 正确的命令: curl -i -X GET http://127.0.0.1:8877?timestamp=1629014688&nonce=c49acbad9673f6eac0a28dcfb90277de&sign=0d7c57a64cf83b086ec8b02e1ceb0fcd
//
// 错误的命令: curl -i -X GET http://127.0.0.1:8877?timestamp=1629014688&sign=0d7c57a64cf83b086ec8b02e1ceb0fcd
