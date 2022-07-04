package app

import (
    "database/sql"
    "fmt"
    "github.com/thewindear/thewindear-blog/internal/model"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "time"
)

type database struct {
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

// NewDatabase 创建一个数据库实例
func newDatabase(dbConf *database) (*gorm.DB, error) {
    dialector := mysql.Open(Config.Database.DSN())
    db, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        err = fmt.Errorf("初始化数据连接失败: %v", err)
        return nil, err
    }
    sqlDB, err := db.DB()
    if err != nil {
        err = fmt.Errorf("设置数据连接失败: %v", err)
        return nil, err
    }
    dbConf.Setting(sqlDB)
    return db, nil
}

func autoMigrate(db *gorm.DB, isMigrate bool) error {
    if isMigrate {
        err := db.AutoMigrate(
            &model.Account{},
            &model.Token{},
            &model.User{},
        )
        if err != nil {
            err = fmt.Errorf("迁移model出错: %v", err)
        }
        return err
    }
    return nil

}
