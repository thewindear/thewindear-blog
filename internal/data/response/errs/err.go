package errs

import (
    "fmt"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "net/http"
)

type Err struct {
    // HttpCode 返回的http码
    HttpCode uint
    // Message 业务描述错误信息
    Message string `json:"message"`
    // Errors HttpCode 422 错误返回错误字段
    Errors []ErrField `json:"errors,omitempty"`
    // Ori 原错误
    Err error
}

// Error 这里会打印原错误信息以及业务错误信息
func (ce *Err) Error() string {
    if utils.ErrNotEmpty(ce.Err) {
        return fmt.Sprintf("[original error]: %s [message]: %s", ce.Err, ce.Message)
    } else {
        return fmt.Sprintf("[message]: %s", ce.Message)
    }
}

// StatusForbidden 资源被禁止访问
func StatusForbidden(message string) *Err {
    return &Err{
        HttpCode: http.StatusForbidden,
        Message:  message,
    }
}

// Unauthorized 认证失败
func Unauthorized(message string) *Err {
    return &Err{
        HttpCode: http.StatusUnauthorized,
        Message:  message,
    }
}

// Conflict 冲突
func Conflict(message string) *Err {
    return &Err{
        HttpCode: http.StatusConflict,
        Message:  message,
    }
}

// BadRequest 解析错误
func BadRequest(message string, oriErr error) *Err {
    return &Err{
        HttpCode: http.StatusBadRequest,
        Message:  message,
        Errors:   nil,
        Err:      oriErr,
    }
}

// DefaultServerError 默认无法处理的错误返回
func DefaultServerError(oriErr error) *Err {
    return &Err{
        HttpCode: http.StatusInternalServerError,
        Message:  "服务异常,请稍后再试",
        Errors:   nil,
        Err:      oriErr,
    }
}
