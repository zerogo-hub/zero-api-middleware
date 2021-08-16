package mustparam

// Option ..
type Option struct {
	fields []Field
}

// Field ..
type Field struct {
	// name 参数名称，如 id
	name string
	// 参数长度
	size int
}

func defaultOption() Option {
	// 默认需要 id, timestamp, nonce, sign 四个字段
	return Option{
		fields: []Field{
			// 时间戳，秒
			{name: "timestamp", size: 10},
			// 随机字符串，32 位
			{name: "nonce", size: 32},
			// 签名，默认使用 github.com/zerogo-hub/zero-api-middleware/sign 签名方式
			{name: "sign", size: 64},
		},
	}
}
