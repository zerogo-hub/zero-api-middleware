// Package bodylimit 限制请求体大小
package bodylimit

import (
	"net/http"

	zeroapi "github.com/zerogo-hub/zero-api"
)

// New 限制请求体大小
//
// limit: 字节数(bytes)，不可以大于该值
//
// 8bit = 1bytes
// 1024bytes = 1kb
// 1024kb = 1mb
func New(method string, limit int64) zeroapi.Handler {
	return func(ctx zeroapi.Context) {
		m := ctx.Method()
		if m != method {
			return
		}

		l := ctx.Request().ContentLength
		if l == -1 || (l == 0 && ctx.Request().Body != nil) {
			ctx.SetHTTPCode(http.StatusBadRequest)
			if _, err := ctx.Message(http.StatusBadRequest, "bad request"); err != nil {
				ctx.App().Logger().Errorf("set message failed, err: %s", err.Error())
			}
			ctx.App().Logger().Warnf("bad request, ctx.Request().ContentLength: %d, method: %s, ip: %s", l, method, ctx.IP())
			ctx.Stopped()
			return
		}
		if l > limit {
			ctx.SetHTTPCode(http.StatusRequestEntityTooLarge)
			if _, err := ctx.Message(http.StatusRequestEntityTooLarge, "request entity too large"); err != nil {
				ctx.App().Logger().Errorf("set message failed, err: %s", err.Error())
			}
			ctx.App().Logger().Warnf("request entity too large, limit: %d, ctx.Request().ContentLength: %d, ip: %s", limit, l, ctx.IP())
			ctx.Stopped()
			return
		}
	}
}
