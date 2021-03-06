// Package jwt ...
package jwt

import (
	"errors"
	"net/http"
	"strings"

	zeroapi "github.com/zerogo-hub/zero-api"
	zerojwt "github.com/zerogo-hub/zero-helper/jwt"
)

// TokenHandler 获取 jwt token 值
type TokenHandler func(ctx zeroapi.Context) (string, error)

// FailedHandler 检查失败时的回调
type FailedHandler func(ctx zeroapi.Context, err error)

// New jwt 验证
//
// jwt 使用 zerojwt.NewJWT() 创建
// onToken 获取 jwt token，默认从 header 中获取
// onFailed 失败时的回调，有默认的回调
func New(jwt zerojwt.JWT, onToken TokenHandler, onFailed FailedHandler) zeroapi.Handler {
	if jwt == nil {
		panic("jwt cant be nil")
	}

	if onToken == nil {
		onToken = TokenFromHeader
	}

	if onFailed == nil {
		onFailed = OnFailed
	}

	return func(ctx zeroapi.Context) {
		tokenValue, err := onToken(ctx)
		if err != nil {
			onFailed(ctx, err)
			return
		}

		if len(tokenValue) == 0 {
			return
		}

		payload, err := jwt.Verify(tokenValue)
		if err != nil {
			onFailed(ctx, err)
			return
		}

		// 将 m 拷贝到 ctx 中
		for k, v := range payload {
			ctx.SetValue(k, v)
		}
	}
}

// TokenFromHeader 从 Header 获取 jwt token 值
func TokenFromHeader(ctx zeroapi.Context) (string, error) {
	h := ctx.Header("Authorization")
	if len(h) == 0 {
		return "", nil
	}

	// Authorization bearer {token}
	headers := strings.Split(h, " ")
	if len(headers) != 2 || strings.ToLower(headers[0]) != "bearer" {
		return "", errors.New("Authorization bearer {token}")
	}

	return headers[1], nil
}

// TokenFromParam 从请求参数中获取 jwt token 值
func TokenFromParam(ctx zeroapi.Context, param string) (string, error) {
	return ctx.Query(param), nil
}

// OnFailed 失败时的回调
func OnFailed(ctx zeroapi.Context, err error) {
	ctx.App().Logger().Errorf("jwt check failed: %v", err)
	ctx.SetHTTPCode(http.StatusUnauthorized)
	ctx.Stopped()
}
