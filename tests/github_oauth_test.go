package tests

import (
    "github.com/thewindear/thewindear-blog/internal/service/oauth2"
    "log"
    "testing"
)

var oauth *oauth2.OAuthGithub

func init() {
    oauth = &oauth2.OAuthGithub{
        ClientSecret: "74a44b1cae9f39e9c9acb3f83fc461f1382bdcd7",
        ClientId:     "afcbad44b4686623123e",
    }
}

func TestGithubOAuthRedirectUri(t *testing.T) {
    log.Println(oauth.RedirectUri("http://localhost/oauth/github/authorize/callback", "abcdefg"))
}

func TestCode2AccessToken(t *testing.T) {
    accessToken, err := oauth.Code2AccessToken("4c94a905abf37f2051a7")
    if err != nil {
        t.Errorf("code 2 access token failure: %s", err.Error())
    } else {
        t.Log(accessToken)
    }
}

func TestAccessToken2Userinfo(t *testing.T) {
    userinfo, err := oauth.AccessToken2UserInfo("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")
    if err != nil {
        t.Errorf("access token to userinfo failure: %s", err.Error())
    } else {
        t.Log(userinfo)
    }
}

func TestUsername2Userinfo(t *testing.T) {
    userinfo, err := oauth.Username2Userinfo("Saherallail7")
    if err != nil {
        t.Errorf("username to userinfo failure: %s", err.Error())
    } else {
        t.Log(userinfo)
    }
}
