package xcto

import (
	zeroapi "github.com/zerogo-hub/zero-api"
)

// New 添加信息头 X-Content-Type-Options: nosniff
func New() zeroapi.Handler {

	return func(ctx zeroapi.Context) {
		ctx.AddHeader("X-Content-Type-Options", "nosniff")
	}
}
