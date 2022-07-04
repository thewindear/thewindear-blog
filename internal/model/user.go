package model

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Uid        string `gorm:"unique;size:12;not null;comment:UID"`
    GithubName string `gorm:"unique;size:30;not null;comment:github name"`
    Email      string `gorm:"unique;size:30;not null;comment:邮箱地址"`
    Phone      string `gorm:"unique;size:11;not null;comment:手机号"`
    Nickname   string `gorm:"unique;size:30;not null;comment:昵称"`
    Exp        uint   `gorm:"default:0;not null;comment:经验值"`
    Gender     uint8  `gorm:"default:0;not null;comment:性别:1-保密,2-女,3-男"`
    Avatar     string `gorm:"size:255;not null;comment:头像"`
    Location   string `gorm:"size:50;not null;comment:地理位置"`
    Intro      string `gorm:"size:255;not null;comment:个人介绍"`
}
