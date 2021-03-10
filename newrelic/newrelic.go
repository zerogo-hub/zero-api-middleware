// Package newrelic 监控
package newrelic

import (
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	zeroapi "github.com/zerogo-hub/zero-api"
)

// New newrelic 监控
// 在 newrelic.com 创建账号，选择 golang app 监控
// 获取 appname 和 license
func New(appname string, license string) zeroapi.Handler {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appname),
		newrelic.ConfigLicense(license),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		return func(ctx zeroapi.Context) {}
	}

	return func(ctx zeroapi.Context) {
		name := ctx.Method() + " " + ctx.Path()
		txn := app.StartTransaction(name)

		txn.SetWebRequestHTTP(ctx.Request())

		w := txn.SetWebResponse(ctx.Response().Writer())
		ctx.Response().SetWriter(w)

		ctx.AppendEnd(func() error {
			txn.End()
			return nil
		})
	}
}
