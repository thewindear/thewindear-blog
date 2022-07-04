package service

import (
    "github.com/thewindear/thewindear-blog/internal/app"
    "github.com/thewindear/thewindear-blog/internal/data/params"
    "github.com/thewindear/thewindear-blog/internal/data/response"
    "github.com/thewindear/thewindear-blog/internal/model"
    "gorm.io/gorm"
)

type IOAuthService interface {
    // LoginPassword 通过帐号密码登录
    LoginPassword(param *params.LoginOauthPasswordParam) (*response.OAuthLoginPasswordResponse, error)
}

type oauthService struct{}

// LoginPassword 通过帐号密码获取账户信息并生成jwt token返回
func (o *oauthService) LoginPassword(param *params.LoginOauthPasswordParam) (*response.OAuthLoginPasswordResponse, error) {
    var account model.Account
    err := app.DB.Model(account).Where("username = ?", param.Username).First(&account).Error
    if err != nil {
        if gorm.ErrRecordNotFound == err {
            return nil, nil
        }
    }

    if !param.CheckPassword(account) {

    }
    return nil, nil
}

// NewOAuthService 返回IOAuthService实例
func NewOAuthService() IOAuthService {
    return &oauthService{}
}
