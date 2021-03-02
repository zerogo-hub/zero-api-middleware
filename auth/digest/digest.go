// Package digest 摘要认证，与简单认证相比，其不在网络中明文传送账号和密码
package digest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	zeroapi "github.com/zerogo-hub/zero-api"
	"github.com/zerogo-hub/zero-helper/crypto"
	"github.com/zerogo-hub/zero-helper/random"
)

// Config 配置
type Config struct {
	// Password 获取密码信息
	Password Handler

	// Users 记录用户账号密码
	Users map[string]string

	// Expires 有效时间
	Expires time.Duration

	// Realm basic realm
	Realm string

	// Digest ..
	Digest string

	// digestlen ..
	digestlen int
}

// Handler 获取账号对应的密码
type Handler func(account string) (string, error)

var defaultConfig = &Config{
	Expires: time.Duration(1 * time.Hour),
	Realm:   strconv.Quote("Authorization Required"),
	Digest:  "Digest ",
}

// New 摘要认证
func New(config *Config) zeroapi.Handler {
	c := defaultConfig
	c.init(config)

	return func(ctx zeroapi.Context) {
		m := ctx.Method()
		h := ctx.Header("Authorization")
		if c.verify(m, h) {
			return
		}

		// 验证不通过，下发 401 请求
		c.failed(ctx)
	}
}

func (c *Config) init(config *Config) {
	if config == nil || (config.Password == nil && config.Users == nil) {
		panic("all be nil")
	}

	if config != nil {
		c.Password = config.Password
		if config.Expires > 0 {
			c.Expires = config.Expires
		}
		if config.Realm != "" {
			c.Realm = config.Realm
		}
		if config.Digest != "" {
			c.Digest = config.Digest
		}
	}

	c.digestlen = len(c.Digest)
	c.Users = config.Users
}

func (c *Config) verify(method, header string) bool {
	// header 示例
	// Digest username="foo", realm="\"Authorization Required\"", nonce="101410811111041093470976875717970", uri="/", algorithm=MD5, response="b5a487d54704ee73a150a2e000cc6da5", qop=auth, nc=00000003, cnonce="4d110299aebb8979"
	fmt.Println(header)

	if header == "" || len(header) < c.digestlen+1 {
		return false
	}

	h := header[:c.digestlen]
	if h != c.Digest {
		return false
	}

	// 解析 "Digest "之后形如 m=n 参数
	s := strings.SplitN(header[c.digestlen:], ",", -1)
	m := make(map[string]string, len(s))

	for _, pair := range s {
		pair := strings.TrimSpace(pair)
		if i := strings.Index(pair, "="); i < 0 {
			m[pair] = ""
		} else {
			k := pair[:i]
			v := pair[i+1:]
			// 去掉 "
			if v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
			m[k] = v
		}
	}

	// algorithm 默认 MD5
	if _, exist := m["algorithm"]; !exist {
		m["algorithm"] = "MD5"
	}

	// 参数检测
	if m["algorithm"] != "MD5" || m["qop"] != "auth" {
		return false
	}

	// 取出密码计算摘要信息
	password, err := c.password(m["username"])
	if err != nil {
		return false
	}

	// RFC 2617
	//
	// HA1 = MD5(A1) = MD5(username:realm:password)
	//
	// 当 qop = "auth"时:
	// HA2 = MD5(A2) = MD5(method:digestURI)
	// response = MD5(HA1:nonce:nc:cnonce:qop:HA2)
	//
	// 当 qop = "auth-int" 时 ...

	a1 := strings.Join([]string{m["username"], c.Realm, password}, ":")
	ha1 := crypto.Md5(a1)

	a2 := strings.Join([]string{method, m["uri"]}, ":")
	ha2 := crypto.Md5(a2)

	response := crypto.Md5(strings.Join([]string{
		ha1, m["nonce"], m["nc"], m["cnonce"], m["qop"], ha2,
	}, ":"))

	return m["response"] == response
}

func (c *Config) failed(ctx zeroapi.Context) {
	nonce := random.String(16)
	askHeader := fmt.Sprintf(`Digest realm="%s",nonce="%s",opaque="",algorithm=MD5,qop="auth"`,
		c.Realm, nonce,
	)
	ctx.SetHeader("WWW-Authenticate", askHeader)
	ctx.SetHTTPCode(http.StatusUnauthorized)

	ctx.Stopped()
}

func (c *Config) password(account string) (string, error) {
	if c.Users != nil {
		if p := c.Users[account]; p != "" {
			return p, nil
		}
	}

	if c.Password != nil {
		return c.Password(account)
	}

	return "", errors.New("password not found")
}
