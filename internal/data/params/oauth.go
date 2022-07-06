package params

import (
    "github.com/thewindear/thewindear-blog/internal/app"
    "github.com/thewindear/thewindear-blog/internal/model"
)

// OAuthPassword 帐号密码登录参数
type OAuthPassword struct {
    Username string `json:"username" validate:"required,email"`
    Password string `json:"password" validate:"required,size=32"`
}

func (lp *OAuthPassword) CheckPassword(account *model.Account) bool {
    return app.Config.Crypt.SaltPassword(lp.Password) == account.Password
}

// CreateOAuthPassword 创建通过账号密码创建账号
type CreateOAuthPassword struct {
    Username string `json:"username" validate:"required,email"`
    Password string `json:"password" validate:"required,size=32"`
}
