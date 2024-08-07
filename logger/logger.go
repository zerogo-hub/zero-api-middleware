// Package logger 请求日志，输出`一次请求`的一些信息，即使发生异常(`panic`)，也能输出
package logger

import (
	"bytes"
	"strconv"
	"sync"
	"time"

	zeroapi "github.com/zerogo-hub/zero-api"
	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

// Config 配置
type Config struct {
	// IP 是否打印 IP，默认 true
	IP bool

	// Code 是否打印 HTTP Code，默认 true
	Code bool

	// Cost 是否打印花费的时间，默认 true
	Cost bool

	// Extend 日志扩展，拼接该函数的返回结果
	Extend func(ctx zeroapi.Context) string
}

var defaultConfig = &Config{
	IP:   true,
	Code: true,
	Cost: true,
}

var (
	ipFlag     = []byte("ip: ")
	sep        = []byte(", ")
	methodFlag = []byte("method: ")
	pathFlag   = []byte("path: ")
	codeFlag   = []byte("code: ")
	costFlag   = []byte("cost: ")
)

// New 记录每一条请求的信息
func New(config ...*Config) zeroapi.Handler {
	var c *Config
	if len(config) > 0 {
		c = config[0]
	} else {
		c = defaultConfig
	}

	return func(ctx zeroapi.Context) {
		start := time.Now()

		ctx.AppendEnd(func() error {
			buff := buffer()
			defer releaseBuffer(buff)

			if c.IP {
				buff.Write(ipFlag)
				buff.Write(zerobytes.StringToBytes(ctx.IP()))
				buff.Write(sep)
			}

			buff.Write(methodFlag)
			buff.Write(zerobytes.StringToBytes(ctx.Method()))
			buff.Write(sep)
			buff.Write(pathFlag)
			buff.Write(zerobytes.StringToBytes(ctx.Path()))
			buff.Write(sep)

			if c.Code {
				buff.Write(codeFlag)
				buff.Write(zerobytes.StringToBytes(strconv.Itoa(ctx.HTTPCode())))
				buff.Write(sep)
			}

			if c.Cost {
				cost := time.Since(start)
				buff.Write(costFlag)
				buff.Write(zerobytes.StringToBytes(cost.String()))
				buff.Write(sep)
			}

			if c.Extend != nil {
				buff.Write(zerobytes.StringToBytes(c.Extend(ctx)))
			}

			ctx.App().Logger().Info(buff.String())
			return nil
		})
	}
}

var bufferPool *sync.Pool

// buffer 从池中获取 buffer
func buffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

// releaseBuffer 将 buff 放入池中
func releaseBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}

func init() {
	bufferPool = &sync.Pool{}
	bufferPool.New = func() interface{} {
		return &bytes.Buffer{}
	}
}
