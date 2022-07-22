package sign

import (
	"bytes"
	"errors"
	"net/http"
	"sort"
	"sync"

	zeroapi "github.com/zerogo-hub/zero-api"
	zerocrypto "github.com/zerogo-hub/zero-helper/crypto"
)

// New 签名验证
// 由插件 middleware/must_param 确保 sign 存在
//
// SecretKey 签名密钥
//
func New(secretKey string, opts ...Option) zeroapi.Handler {
	opt := defaultOption()
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(ctx zeroapi.Context) {
		if opt.Enable {
			if err := checkSign(opt.SignName, secretKey, ctx.QueryAll()); err != nil {
				ctx.Stopped()
				ctx.SetHTTPCode(http.StatusBadRequest)
				ctx.App().Logger().Errorf("check sign failed, method: %s, path: %s", ctx.Method(), ctx.Path())
			}
		}
	}
}

// checkSign 计算签名值
func checkSign(signName string, secretKey string, values map[string][]string) error {
	if len(values) == 0 {
		return errors.New("no param")
	}

	if len(values[signName]) == 0 {
		return errors.New("miss sign param")
	}

	sign := values[signName][0]
	if sign == "" {
		return errors.New("sign is empty")
	}

	calcSign, err := calcSign(secretKey, signName, values)
	if err != nil {
		return errors.New("calc sign failed")
	}

	if calcSign != sign {
		return errors.New("sign check failed")
	}

	return nil
}

// calcSign 计算签名
func calcSign(secretKey, signName string, values map[string][]string) (string, error) {
	// 所有参数按照字母顺序从小到大排列
	// 所有参数形成如 key1=value1key2=value2 的形式
	size := len(values)
	keys := make([]string, 0, size)
	for key := range values {
		// sign 不参与签名
		if key != signName {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	b := buffer()
	defer releaseBuffer(b)
	b.Reset()

	for _, key := range keys {
		if key == "" {
			continue
		}

		b.WriteString(key)
		b.WriteByte('=')
		vvs := values[key]
		for idx, vv := range vvs {
			// 相同的参数使用 , 连接
			// 例如: a=1&a=2&a=3，会构造成字符串: a=1,2,3
			b.WriteString(vv)
			if idx != len(vvs)-1 {
				b.WriteByte(',')
			}
		}
	}

	signStr := b.String()
	return calcWithHmacSha256(secretKey, signStr), nil
}

// calcWithHmacSha256 使用 HmacSha256 进行签名
func calcWithHmacSha256(secretKey, signStr string) string {
	return zerocrypto.HmacSha256(signStr, secretKey)
}

var bufferPool *sync.Pool

func buffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

func releaseBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}

func init() {
	bufferPool = &sync.Pool{}
	bufferPool.New = func() interface{} {
		return &bytes.Buffer{}
	}
}
