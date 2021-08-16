package mustparam

import (
	"net/http"

	zeroapi "github.com/zerogo-hub/zero-api"
)

// New 对请求中的必要参数进行验证
func New(opts ...Option) zeroapi.Handler {
	opt := defaultOption()
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(ctx zeroapi.Context) {
		params := ctx.QueryAll()
		if params == nil {
			ctx.Stopped()
			ctx.SetHTTPCode(http.StatusBadRequest)
			ctx.App().Logger().Errorf("params is null, method: %s, path: %s", ctx.Method(), ctx.Path())
			return
		}

		for _, field := range opt.fields {
			param, ok := params[field.Name]

			if !ok || len(param) == 0 {
				ctx.Stopped()
				ctx.SetHTTPCode(http.StatusBadRequest)
				ctx.App().Logger().Errorf("miss param: %s, method: %s, path: %s", field.Name, ctx.Method(), ctx.Path())
				return
			}

			if len(param[0]) != field.Size {
				ctx.Stopped()
				ctx.SetHTTPCode(http.StatusBadRequest)
				ctx.App().Logger().Errorf("field size wrong, field.name: %s, required size: %d, current siz: %d, method: %s, path: %s", field.Name, field.Size, len(param[0]), ctx.Method(), ctx.Path())
				return
			}
		}
	}
}
