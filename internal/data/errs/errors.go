package errs

import (
    "fmt"
)

type Err struct {
    // HttpCode 返回的http码
    HttpCode uint
    // Message 业务描述错误信息
    Message string
    // Ori 原错误
    Err error
}

// Error 这里会打印原错误信息以及业务错误信息
func (ce *Err) Error() string {
    return fmt.Sprintf("[original error]: %s, [message]: %s", ce.Err, ce.Message)
}
