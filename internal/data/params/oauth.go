package params

import (
    "github.com/thewindear/thewindear-blog/internal/app"
    "github.com/thewindear/thewindear-blog/internal/model"
)

// LoginOauthPasswordParam 帐号密码登录参数
type LoginOauthPasswordParam struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Captcha  string `json:"captcha,omitempty"`
}

func (lp *LoginOauthPasswordParam) CheckPassword(account model.Account) bool {
    return app.Config.Crypt.SaltPassword(lp.Password) == account.Password
}
