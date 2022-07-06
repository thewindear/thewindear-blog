package model

import (
    "gorm.io/gorm"
)

const (
    // PlatformAccount 本系统帐号
    PlatformAccount = 1
    // PlatformGithub github帐号
    PlatformGithub = 1
    // TokenStatusEnable token启用
    TokenStatusEnable = 2
    // TokenStatusDisable token被禁用
    TokenStatusDisable = 1
)

// Token 令牌表
type Token struct {
    gorm.Model
    AccountId uint   `gorm:"not null;default:0;comment:帐号id;index:idx_uniq_bind,unique"`
    UserId    uint   `gorm:"not null;default:0;comment:用户id;index:idx_uniq_bind,unique"`
    Token     string `gorm:"unique;not null;size:32;comment:令牌"`
    FromId    string `gorm:"unique;not null;size:60;comment:平台唯一id"`
    Platform  uint8  `gorm:"not null;default:0;comment:平台:1-帐号,2-github"`
    Status    uint8  `gorm:"not null;default:2;comment:状态:1-禁用,2-启用"`
}

// IsEnable 是否启用
func (t *Token) IsEnable() bool {
    return t.Status == TokenStatusEnable
}

func NewAccountToken(token, fromId string, accountId, userId uint) *Token {
    return &Token{
        AccountId: accountId,
        UserId:    userId,
        Token:     token,
        FromId:    fromId,
        Platform:  PlatformAccount,
        Status:    TokenStatusEnable,
    }
}
