package sign

// Option ..
type Option struct {
	// Enable 是否开启检测，默认 true
	Enable bool
	// SignName sign 的参数名称, 一般为 sign, 与 middleware/must-param 中一致
	SignName string
}

func defaultOption() Option {
	return Option{
		Enable:   true,
		SignName: "sign",
	}
}
