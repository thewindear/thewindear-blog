package app

import (
    "database/sql"
    "fmt"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "time"
)

type (
    Conf struct {
        Database *database
        Server   *server
        Crypt    *crypt
    }
    server struct {
        Name string //服务名
        Port uint   //端口号
    }
    database struct {
        Host     string //主机名
        Port     uint   //端口号
        Username string //用户名
        Password string //密码
        Database string //数据库名
        Params   string //参数
        Idle     int    //空闲数
        Conns    int    //最大连接数
        Lifetime uint   //生命周期
        Migrate  bool   //自动迁移model
    }
    crypt struct {
        Password string //密码加盐
        Token    string //token加盐
    }
)

// DSN 生成数据gorm数据库连接dsn字符串
func (d *database) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", d.Username, d.Password, d.Host, d.Port, d.Database, d.Params)
}

// Setting 配置数据库的优化设置
func (d *database) Setting(sqlDB *sql.DB) {
    sqlDB.SetMaxIdleConns(d.Idle)
    sqlDB.SetMaxOpenConns(d.Conns)
    sqlDB.SetConnMaxLifetime(time.Duration(d.Lifetime))
}

// SaltPassword 给原密码加上salt
func (c *crypt) SaltPassword(pwd string) string {
    return utils.CryptMD5(pwd, c.Password)
}

// SaltToken 给原token加上salt
func (c *crypt) SaltToken(token string) string {
    return utils.CryptMD5(token, c.Token)
}
