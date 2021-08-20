package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zamsign "github.com/zerogo-hub/zero-api-middleware/sign"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)

	signKey := "123456"
	a.Use(zamsign.New(signKey))

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Errorf("app run failed, err: %s", err.Error())
	}
}

// 默认有 timestamp, nonce, sign 三个必要字段
// 正确的命令: curl -i -X GET http://127.0.0.1:8877?timestamp=1629014688&nonce=c49acbad9673f6eac0a28dcfb90277de&sign=a01d70c47fd6f07e476e5c851bccfc7590abea9f8aaf97cdbb4dc5c8959808cf
