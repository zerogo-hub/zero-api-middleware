// Package basic 简单认证，如果之前没有认证，会返回`401`，要求用户输入用户名和密码
package basic

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	zeroapi "github.com/zerogo-hub/zero-api"
)

// Config 配置
type Config struct {
	users map[string]*user

	// Check 验证账号和密码是否正确
	Check Handler

	// Expires 有效时间
	Expires time.Duration

	// Realm basic realm
	Realm string

	// Basic ..
	Basic string

	// basicLen "Basic "的长度
	basicLen int
}

var defaultConfig = &Config{
	Expires: time.Duration(1 * time.Hour),
	Realm:   strconv.Quote("Authorization Required"),
	Basic:   "Basic ",
}

// Handler 验证账号和密码是否正确
type Handler func(account, password string) bool

type user struct {
	account string
	expire  time.Time
	// logged 曾经是否登录过
	logged bool
}

// New 简单认证
// accounts 内置一些账号密码，可以为 nil
func New(config *Config, accounts map[string]string) zeroapi.Handler {
	c := defaultConfig
	c.init(config, accounts)

	return func(ctx zeroapi.Context) {
		h := ctx.Header("Authorization")
		if c.verify(h) {
			return
		}

		// 验证不通过，下发 401 请求
		c.failed(ctx)
	}
}

func (c *Config) init(config *Config, accounts map[string]string) {
	if (config == nil || config.Check == nil) && accounts == nil {
		panic("all be nil")
	}

	if config != nil {
		c.Check = config.Check
		if config.Expires > 0 {
			c.Expires = config.Expires
		}
		if config.Realm != "" {
			c.Realm = config.Realm
		}
		if config.Basic != "" {
			c.Basic = config.Basic
		}
	}

	c.basicLen = len(c.Basic)

	if len(accounts) > 0 {
		c.users = make(map[string]*user, len(accounts))
		for account, password := range accounts {
			if account == "" {
				continue
			}
			// header = Basic Zm9vOmJhcg==
			header := c.Basic + base64.StdEncoding.EncodeToString([]byte(account+":"+password))
			c.users[header] = &user{
				account: account,
				logged:  false,
			}
		}
	}
}

func (c *Config) verify(header string) bool {
	if header == "" || len(header) < c.basicLen+1 {
		return false
	}

	if header[:c.basicLen] != c.Basic {
		return false
	}

	if c.users != nil {
		user := c.users[header]
		// 在已有已有记录中查询
		if user != nil {
			if !user.logged {
				// 尚未登录过
				user.logged = true
				user.expire = time.Now().Add(c.Expires)
				return true
			} else if time.Now().After(user.expire) {
				// 登录超时
				delete(c.users, header)
				return false
			}

			return true
		}

		if c.Check != nil {
			// 解码，取出 account 和 password
			a, err := base64.StdEncoding.DecodeString(header[c.basicLen:])
			if err == nil {
				b := string(a)
				for i := 0; i < len(b); i++ {
					if b[i] == ':' {
						if c.Check(b[:i], b[i+1:]) {
							// 验证通过
							return true
						}
						return false
					}
				}
			}
		}

		return false
	}

	return false
}

func (c *Config) failed(ctx zeroapi.Context) {
	askHeader := c.Basic + " realm=" + c.Realm
	ctx.SetHeader("WWW-Authenticate", askHeader)
	ctx.SetHTTPCode(http.StatusUnauthorized)
	ctx.Stopped()
}
