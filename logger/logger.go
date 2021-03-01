package logger

import (
	"bytes"
	"strconv"
	"sync"
	"time"

	zeroapi "github.com/zerogo-hub/zero-api"
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
				buff.WriteString("ip: ")
				buff.WriteString(ctx.IP())
				buff.WriteString(", ")
			}

			buff.WriteString("method: ")
			buff.WriteString(ctx.Method())
			buff.WriteString(", ")
			buff.WriteString("path: ")
			buff.WriteString(ctx.Path())
			buff.WriteString(", ")

			if c.Code {
				buff.WriteString("code: ")
				buff.WriteString(strconv.Itoa(ctx.HTTPCode()))
				buff.WriteString(", ")
			}

			if c.Cost {
				cost := time.Since(start)
				buff.WriteString("cost: ")
				buff.WriteString(cost.String())
				buff.WriteString(", ")
			}

			if c.Extend != nil {
				buff.WriteString(c.Extend(ctx))
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
