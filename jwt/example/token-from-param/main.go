package main

import (
	zeroapi "github.com/zerogo-hub/zero-api"
	zamjwt "github.com/zerogo-hub/zero-api-middleware/jwt"
	app "github.com/zerogo-hub/zero-api/app"
	zerojwt "github.com/zerogo-hub/zero-helper/jwt"
)

var jwt zerojwt.JWT

func init() {
	jwt = zerojwt.NewJWT()
}

func createtoken(ctx zeroapi.Context) {
	// 用户自定义数据，将存储在 jwt token 中
	values := map[string]interface{}{
		"id":   10001,
		"name": "JWT",
	}

	// 模拟登陆成功
	// 生成一个 jwt token 用于测试，默认有效期 5 分钟
	token, err := jwt.Token(values)
	if err != nil {
		ctx.Textf("create token failed: %s", err.Error())
		return
	}

	ctx.Textf("please visit http://127.0.0.1:8877/check?token=%s", token)
}

func checktoken(ctx zeroapi.Context) {
	// 输出存储在 jwt token 中的内容
	id := ctx.Value("id")
	name := ctx.Value("name")

	ctx.Textf("id: %v, name: %v", id, name)
}

func main() {
	a := app.New()

	// 创建 jwt token
	a.Get("/", createtoken)

	// 检查 jwt token
	a.Get("/check", zamjwt.New(jwt, onToken, nil), checktoken)

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Errorf("app run failed, err: %s", err.Error())
	}
}

func onToken(ctx zeroapi.Context) (string, error) {
	return zamjwt.TokenFromParam(ctx, "token")
}
