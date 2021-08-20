package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zamnewrelic "github.com/zerogo-hub/zero-api-middleware/newrelic"
	app "github.com/zerogo-hub/zero-api/app"
	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func sleepHandle(ctx zeroapi.Context) {
	zerotime.Sleep(3)
	if _, err := ctx.Text("sleep over"); err != nil {
		ctx.App().Logger().Errorf("set text failed, err: %s", err.Error())
	}
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)
	a.Get("/sleep", sleepHandle)

	appname := "zerogo-test"
	license := "7f577e117380b08a2383a87ba5e417ba5c5aNRAL"
	a.Use(zamnewrelic.New(appname, license))

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Errorf("app run failed, err: %s", err.Error())
	}
}
