package utils

import (
    "crypto/md5"
    "fmt"
    "io"
)

// CryptMD5 变长参数md5
func CryptMD5(args ...string) string {
    w := md5.New()
    for _, arg := range args {
        _, _ = io.WriteString(w, arg)
    }
    return fmt.Sprintf("%x", w.Sum(nil))
}
