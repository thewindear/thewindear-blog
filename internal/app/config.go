package app

import (
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "strconv"
)

type (
    conf struct {
        Database *database
        Server   *server
        Crypt    *crypt
    }
    server struct {
        Name string //服务名
        Port int    //端口号
    }
    crypt struct {
        Password string //密码加盐
        Token    string //token加盐
    }
)

func (s *server) ListenHost() string {
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

// SaltToken 给原token加上salt
func (c *crypt) SaltToken(token string) string {
    return utils.CryptMD5(token, c.Token)
}
