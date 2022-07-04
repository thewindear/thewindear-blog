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

func TestOAuthLoginPassword(t *testing.T) {
    param := &params.LoginOauthPasswordParam{
        Username: "thewindear@outlook.com",
        Password: utils.CryptMD5("laiwenbang"),
    }
    token, err := oauthService.LoginPassword(param)
    if err != nil {
        t.Fatalf("login failure error: %s", err)
    } else {
        t.Logf("login success: %v", token)
    }
}
