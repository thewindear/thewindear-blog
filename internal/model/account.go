package model

import (
    "gorm.io/gorm"
)

const (
    // AccountUnactivated 未激活
    AccountUnactivated = 1
    // AccountActivated 已激活
    AccountActivated = 2
)

// Account 帐号表
type Account struct {
    gorm.Model
    Username string `gorm:"unique;size:36;not null;comment:用户名"`
    Password string `gorm:"size:50;not null;comment:密码"`
    IsActive uint8  `gorm:"default:2;comment:是否被激活:1-未激活,2-已激活"`
}
