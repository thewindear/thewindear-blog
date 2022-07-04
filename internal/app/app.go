package app

import (
    "gorm.io/gorm"
)

var Config *conf
var DB *gorm.DB

// Run 传入配置文件路径启动应用如果存在错误返回 error 错误
// 外部打印 error 错误信息并停止应用
func Run(configPath string) (err error) {
    if err = Init(configPath); err != nil {
        return err
    }
    if err = newRoute(Config); err != nil {
        return err
    }
    return nil
}

// Init 初始化全局变量
func Init(configPath string) (err error) {
    Config, err = newConfig(configPath)
    if err != nil {
        return err
    }
    DB, err = newDatabase(Config.Database)
    if err != nil {
        return err
    }
    if err = autoMigrate(DB, Config.Database.Migrate); err != nil {
        return err
    }
    return nil
}
