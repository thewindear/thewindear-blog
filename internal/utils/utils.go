package utils

import (
    "crypto/md5"
    "errors"
    "fmt"
    gonanoid "github.com/matoous/go-nanoid"
    "gorm.io/gorm"
    "io"
    "strings"
)

// CryptMD5 变长参数md5
func CryptMD5(args ...string) string {
    w := md5.New()
    for _, arg := range args {
        _, _ = io.WriteString(w, arg)
    }
    return fmt.Sprintf("%x", w.Sum(nil))
}

// IsRecordNotFound 是否数据为空
func IsRecordNotFound(err error) bool {
    return errors.Is(err, gorm.ErrRecordNotFound)
}

// NotNull 传入的数据不等于空
func NotNull(data interface{}) bool {
    return data != nil
}

// IsNull 是否为空
func IsNull(data interface{}) bool {
    return data == nil
}

// ErrNotEmpty 判断err不为空
func ErrNotEmpty(err error) bool {
    return !errors.Is(err, nil)
}

// Uid 用户id
func Uid() string {
    return RandomStr(1, 8)
}

// PostId 文章id
func PostId() string {
    return RandomStr(1, 10)
}

// MakeToken 随机生成token
func MakeToken() string {
    return RandomStr(2, 16)
}

// GetEmailUsername 获取邮箱用户名
func GetEmailUsername(email string) string {
    usernameAndDomain := strings.Split(email, "@")
    return usernameAndDomain[0]
}

// RandomStr 随机生成指定等级和长度的字符串
func RandomStr(level, size int) (randStr string) {
    seed := "0123456789qazxswedcvfrtgbnhyyujmkiolp"
    if level >= 1 {
        seed += "QAZXSWEDCVFRTGBNHYYUJMKIOLP-"
    }
    if level >= 2 {
        seed += "!@#$%^&*()_+"
    }
    randStr, _ = gonanoid.Generate(seed, size)
    return randStr
}
