package app

import (
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/thewindear/thewindear-blog/internal/model"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var Config *Conf
var DB *gorm.DB

// Run 传入配置文件路径启动应用如果存在错误返回 error 错误
// 外部打印 error 错误信息并停止应用
func Run(configPath string) (err error) {
    err = parseConfig(configPath)
    if err != nil {
        return
    }
    err = initDB()
    if err != nil {
        return
    }
    if Config.Database.Migrate {
        if err = autoMigrate(); err != nil {
            return err
        }
    }
    return nil
}

// parseConfig 用于解析应用配置文件
func parseConfig(configPath string) (err error) {
    _, err = toml.DecodeFile(configPath, &Config)
    if err != nil {
        err = fmt.Errorf("解析配置文件失败: %v", err)
        return
    }
    return
}

func initDB() (err error) {
    dialector := mysql.Open(Config.Database.DSN())
    DB, err = gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        err = fmt.Errorf("初始化数据连接失败: %v", err)
        return
    }
    sqlDB, err := DB.DB()
    if err != nil {
        err = fmt.Errorf("设置数据连接失败: %v", err)
        return
    }
    Config.Database.Setting(sqlDB)
    return nil
}

func autoMigrate() error {
    err := DB.AutoMigrate(
        &model.Account{},
        &model.Token{},
        &model.User{},
    )
    if err != nil {
        err = fmt.Errorf("迁移model出错: %v", err)
    }
    return err
}
