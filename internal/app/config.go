package app

import (
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "strconv"
    "time"
)

type (
    conf struct {
        Database    *database
        Application *application
        Crypt       *crypt
    }
    application struct {
        Name        string //服务名
        Port        int    //端口号
        Domain      string //域名
        TokenExpire uint   `toml:"token-expire"` //token过期时长
    }
    crypt struct {
        Password string //密码加盐
        Token    string //token加盐
    }
)

func (s *application) TokenExpireSeconds() time.Duration {
    return time.Duration(s.TokenExpire) * time.Second
}

func (s *application) ListenHost() string {
    return ":" + strconv.Itoa(s.Port)
}

// newConfig 用于解析应用配置文件
func newConfig(configPath string) (*conf, error) {
    var conf conf
    _, err := toml.DecodeFile(configPath, &conf)
    if err != nil {
        err = fmt.Errorf("解析配置文件失败: %v", err)
        return nil, err
    }
    return &conf, nil
}

// SaltPassword 给原密码加上salt
func (c *crypt) SaltPassword(pwd string) string {
    return utils.CryptMD5(pwd, c.Password)
}
