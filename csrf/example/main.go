package main

import (
	zeroapi "github.com/zerogo-hub/zero-api"
	zamcsrf "github.com/zerogo-hub/zero-api-middleware/csrf"
	app "github.com/zerogo-hub/zero-api/app"
)

func createTokenHandle(ctx zeroapi.Context) {
	ctx.Textf("http://127.0.0.1:8877/verify")
}

func verifyTokenHandle(ctx zeroapi.Context) {
	ctx.Text("verify success")
}

func main() {
	a := app.New()

	a.Use(zamcsrf.New())

	a.Get("/", createTokenHandle)
	a.Get("/verify", zamcsrf.Verify(), verifyTokenHandle)

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Errorf("app run failed, err: %s", err.Error())
	}
}
