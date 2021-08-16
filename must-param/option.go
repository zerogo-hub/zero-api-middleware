package mustparam

// Option ..
type Option struct {
	// Fields 必要参数列表
	Fields []Field
}

// Field ..
type Field struct {
	// name 参数名称，如 id
	Name string
	// 参数长度
	Size int
}

func defaultOption() Option {
	// 默认需要 id, timestamp, nonce, sign 四个字段
	return Option{
		Fields: []Field{
			// 时间戳，秒
			{Name: "timestamp", Size: 10},
			// 随机字符串，32 位
			{Name: "nonce", Size: 32},
			// 签名，默认使用 github.com/zerogo-hub/zero-api-middleware/sign 签名方式
			{Name: "sign", Size: 64},
		},
	}
}
