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

// ErrNotEmpty 判断err不为空
func ErrNotEmpty(err error) bool {
    return !errors.Is(err, nil)
}

// Uid 用户id
func Uid() string {
    return nanoId(1, 8)
}

// PostId 文章id
func PostId() string {
    return nanoId(1, 10)
}

// MakeToken 随机生成token
func MakeToken() string {
    return nanoId(2, 16)
}

// GetEmailUsername 获取邮箱用户名
func GetEmailUsername(email string) string {
    usernameAndDomain := strings.Split(email, "@")
    return usernameAndDomain[0]
}

func nanoId(level, size int) (token string) {
    seed := "0123456789qazxswedcvfrtgbnhyyujmkiolp"
    if level >= 1 {
        seed += "QAZXSWEDCVFRTGBNHYYUJMKIOLP-"
    }
    if level >= 2 {
        seed += "!@#$%^&*()_+"
    }
    token, _ = gonanoid.Generate(seed, size)
    return token
}
