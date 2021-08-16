package timestamp

// Option ..
type Option struct {
	// Field 时间戳字段名称，默认为 timestamp
	Field string
	// Diff 相差的时间，单位 秒，默认 10 秒
	Diff int64
	// Enable 是否开启检测，默认 true
	Enable bool
}

func defaultOption() Option {
	return Option{
		Field:  "timestamp",
		Diff:   10,
		Enable: true,
	}
}
