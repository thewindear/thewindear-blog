package model

import (
    "github.com/thewindear/thewindear-blog/internal/utils"
    "github.com/thewindear/thewindear-blog/pkg/oauth2"
    "gorm.io/gorm"
)

const (
    // NormalUserRule 普通用户角色
    NormalUserRule = 1
    // SuperUserRule 超级管理员角色
    SuperUserRule = 2
)

type User struct {
    gorm.Model
    Uid        string `gorm:"unique;size:12;comment:UID"`
    GithubName string `gorm:"unique;size:30;comment:github name"`
    Email      string `gorm:"unique;size:30;comment:邮箱地址"`
    Phone      string `gorm:"unique;size:11;comment:手机号"`
    Nickname   string `gorm:"unique;size:30;comment:昵称"`
    Exp        uint   `gorm:"default:0;not null;comment:经验值"`
    Gender     uint8  `gorm:"default:0;not null;comment:性别:1-保密,2-女,3-男"`
    Avatar     string `gorm:"size:255;not null;comment:头像"`
    Location   string `gorm:"size:50;not null;comment:地理位置"`
    Intro      string `gorm:"size:255;not null;comment:个人介绍"`
    Rule       uint8  `gorm:"size:1;not null;default:1;comment:超级管理员:1-正常用户,2-超级管理员"`
}

// IsSuperAdmin 是否为超级管理员
func (u *User) IsSuperAdmin() bool {
    return u.Rule == SuperUserRule
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.Uid = utils.Uid()
    if u.Nickname == "" && u.Email != "" {
        u.Nickname = utils.GetEmailUsername(u.Email)
    }
    u.Rule = NormalUserRule
    return nil
}

func NewDefaultUser(email string) *User {
    return &User{
        Email: email,
    }
}

// NewOAuth2User 通过OAuth2User返回的用户信息生成User模型
func NewOAuth2User(oauth2User *oauth2.UserInfo) *User {
    var user = &User{
        Nickname: oauth2User.Nickname,
        Avatar:   oauth2User.Avatar,
    }
    switch oauth2User.From {
    case oauth2.IOAuth2Github:
        user.GithubName = oauth2User.Username
        break
    }
    return user
}
