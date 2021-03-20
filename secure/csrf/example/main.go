package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	"github.com/zerogo-hub/zero-api-middleware/secure/csrf"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("method: %s, `ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", ctx.Method(), pid, pid)
}

func main() {
	a := zeroapi.New()

	a.Use(csrf.New(nil))

	a.Get("/", helloworldHandle)

	a.Post("/", helloworldHandle)

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}

// GET 请求会添加 csrf token
// curl -i -X GET http://127.0.0.1:8877
//
// HTTP/1.1 200 OK
// Set-Cookie: csrfToken=b6ac200843fab18ad8fe19befcf81e9b|1616254411|5ecac4d6e715a0b74fc7026bfcd25703; Max-Age=86400; HttpOnly
// Vary: Cookie
// Date: Sat, 20 Mar 2021 15:33:31 GMT
// Content-Length: 87
// Content-Type: text/plain; charset=utf-8
//
// POST 请求在 header 中携带 token
// curl -i -X POST -H "x-csrf-token:b6ac200843fab18ad8fe19befcf81e9b" --cookie "csrfToken=b6ac200843fab18ad8fe19befcf81e9b|1616254411|5ecac4d6e715a0b74fc7026bfcd25703" http://127.0.0.1:8877
//
// HTTP/1.1 200 OK
// Date: Sat, 20 Mar 2021 15:34:15 GMT
// Content-Length: 88
// Content-Type: text/plain; charset=utf-8
