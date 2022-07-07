package test_services

import (
    "github.com/thewindear/thewindear-blog/internal/data/params"
    "github.com/thewindear/thewindear-blog/internal/service"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "github.com/thewindear/thewindear-blog/tests"
    "testing"
)

var oauthService service.IOAuthService

func init() {
    tests.InitApp()
    oauthService = service.NewOAuthService()
}

func TestCreateOAuthPasswordAccount(t *testing.T) {
    param := &params.CreateOAuthPassword{
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
