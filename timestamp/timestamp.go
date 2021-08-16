package timestamp

import (
	"net/http"
	"strconv"

	zeroapi "github.com/zerogo-hub/zero-api"

	zerotime "github.com/zerogo-hub/zero-helper/time"
)

// New 验证请求中的时间戳与服务端相比，是否相差太大
// 由插件 middleware/must-param 确保 Field 存在
func New(opts ...Option) zeroapi.Handler {

	opt := defaultOption()
	if len(opts) > 0 {
		opt = opts[0]
		if opt.Enable && opt.Diff <= 0 {
			panic("opt.Diff must bigger than zero")
		}
	}

	return func(ctx zeroapi.Context) {
		if opt.Enable {
			s := ctx.Query(opt.Field)

			timestamp, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				ctx.Stopped()
				ctx.SetHTTPCode(http.StatusBadRequest)
				ctx.App().Logger().Errorf("no valid timestamp: %d, method: %s, path: %s", timestamp, ctx.Method(), ctx.Path())
				return
			}

			now := zerotime.Now()

			if timestamp > now+opt.Diff || timestamp < now-opt.Diff {
				ctx.Stopped()
				ctx.SetHTTPCode(http.StatusBadRequest)
				ctx.App().Logger().Errorf("invalid timestamp: %d, method: %s, path: %s", timestamp, ctx.Method(), ctx.Path())
				return
			}
		}
	}
}
