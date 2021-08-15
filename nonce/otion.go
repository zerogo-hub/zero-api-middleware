package nonce

// Option ..
type Option struct {
	// Field nonce 字段名称，默认为 nonce
	Field string
	// 在一定时间内，nonce 不可以重复，单位 秒
	// 如果开启了 middleware/timestamp 验证，该时间 >= timestamp.Diff 即可
	Expire string
	// Enable 是否开启检测，默认 true
	Enable bool
	// PrefixNonce 在缓存中的前缀, 如 "nonce:"
	PrefixNonce string
}

func defaultOption() Option {
	return Option{
		Field:       "nonce",
		Expire:      "10",
		Enable:      true,
		PrefixNonce: "nonce:",
	}
}
