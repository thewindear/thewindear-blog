package test_services

import (
    "github.com/thewindear/thewindear-blog/internal/app"
    "github.com/thewindear/thewindear-blog/internal/data/params"
    "github.com/thewindear/thewindear-blog/internal/model"
    "github.com/thewindear/thewindear-blog/internal/service"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "github.com/thewindear/thewindear-blog/pkg/oauth2"
    "github.com/thewindear/thewindear-blog/tests"
    "testing"
)

var oauthService service.IOAuthService
var user *model.User

func init() {
    tests.InitApp()
    oauthService = service.NewOAuthService()
    app.DB.Model(user).Where("id = ?", 5).First(&user)
}

func TestCreateOAuthPasswordAccount(t *testing.T) {
    param := &params.PutOAuthPassword{
        Username: "thewindear@outlook.com",
        Password: utils.CryptMD5("laiwenbang"),
    }
    res, err := oauthService.CreateOAuthPasswordAccount(param)
    if utils.ErrNotEmpty(err) {
        t.Fatalf("注册账号失败原因: %s", err)
    }
    t.Logf("注册账号成功 生成 jwt token: \n %s", res.Token)
}

func TestOAuthPasswordLogin(t *testing.T) {
    param := &params.OAuthPassword{
        Username: "thewindear@outlook.com",
        Password: utils.CryptMD5("laiwenbang"),
    }
    res, err := oauthService.OAuthLoginPassword(param)
    if utils.ErrNotEmpty(err) {
        t.Fatalf("登录失败: %s", err)
    }
    t.Logf("登录成功 生成 jwt token: \n %s", res.Token)
}

func TestOAuthGithubLogin(t *testing.T) {
    res, err := oauthService.OAuth2App("github", "eb0fe46426ce4f50e964")
    if utils.ErrNotEmpty(err) {
        t.Fatalf("登录失败: %s", err)
    }
    t.Logf("登录成功 生成 jwt token: \n %s", res.Token)
}

func TestOAuthBindEmail(t *testing.T) {
    err := oauthService.BindEmail(user, "375500819@qq.com")
    if utils.ErrNotEmpty(err) {
        t.Fatalf("绑定邮箱失败: %s", err)
    } else {
        t.Logf("绑定邮箱成功 现在的邮箱是: %s", user.Email)
    }
}

func TestOAuthUnBindEmail(t *testing.T) {
    err := oauthService.UnBindEmail(user)
    if utils.ErrNotEmpty(err) {
        t.Fatalf("解绑邮箱失败: %v", err)
    } else {
        t.Logf("解绑成功当前邮箱为空: %s", user.Email)
    }
}

func TestOAuthUnBindOAuth2App(t *testing.T) {
    err := oauthService.UnBindOAuth2App(user, oauth2.IOAuth2Github)
    if utils.ErrNotEmpty(err) {
        t.Fatalf("解绑绑定github失败: %v", err)
    } else {
        t.Logf("解绑github成功: %s", user.GithubName)
    }
}
