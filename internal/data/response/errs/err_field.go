package errs

// ErrField 错误字段类型
type ErrField struct {
    // Field 字段名
    Field string `json:"field"`
    // Message 错误说明
    Message string `json:"message"`
    // Code 表单验证错误码
    Code string `json:"code"`
}
