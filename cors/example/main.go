package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zamcors "github.com/zerogo-hub/zero-api-middleware/cors"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := app.New()

	a.Get("/", helloworldHandle)

	a.Use(zamcors.New(&zamcors.Config{
		AccessControlAllowOrigin: []string{"*.abc.com", "abc.com"},
	}))

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}

// 预检通过
// 命令: curl -i -X OPTIONS -H "Origin:api.abc.com" -H "Access-Control-Request-Method:GET" http://127.0.0.1:8877
// 返回:
// HTTP/1.1 204 No Content
// Access-Control-Allow-Methods: HEAD,GET,POST,PUT,PATCH,DELETE
// Access-Control-Allow-Origin: api.abc.com
// Access-Control-Max-Age: 600
//
// 预检未通过
// 命令: curl -i -X OPTIONS -H "Origin:test.com" -H "Access-Control-Request-Method:GET" http://127.0.0.1:8877
// 返回:
// HTTP/1.1 403 Forbidden
// Content-Length: 0
//
// 正常请求通过
// 命令: curl -i -X GET -H "Origin:api.abc.com" -H "Access-Control-Request-Method:GET" http://127.0.0.1:8877
// 返回:
// HTTP/1.1 200 OK
// Access-Control-Allow-Origin: api.abc.com
// Content-Length: 74
// Content-Type: text/plain; charset=utf-8
//
// `ctrl+c` to close, `kill 10625` to shutdown, `kill -USR2 10625` to restart%
