package oauth2

import (
    "encoding/json"
    "fmt"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
)

const (
    timeout               = 10
    githubRedirectBaseUri = "https://github.com/login/oauth/authorize"
    githubAccessToken     = "https://github.com/login/oauth/access_token"
    githubUser            = "https://api.github.com/user"
    githubUsers           = "https://api.github.com/users"
)

type OAuthGithub struct {
    ClientId     string
    ClientSecret string
}

// RedirectUri 生成github授权url跳转地址
func (o *OAuthGithub) RedirectUri(callbackUri, state string) string {
    param := url.Values{}
    param.Add("client_id", o.ClientId)
    param.Add("redirect_uri", callbackUri)
    param.Add("scope", "user:read")
    param.Add("state", state)
    param.Add("allow_signup", "true")
    return githubRedirectBaseUri + "?" + param.Encode()
}

// Code2AccessToken 通过授权的code返回对应 AccessAccessToken
func (o *OAuthGithub) Code2AccessToken(code string) (accessToken *AccessToken, err error) {
    param := url.Values{}
    param.Add("client_id", o.ClientId)
    param.Add("client_secret", o.ClientSecret)
    param.Add("code", code)
    api := githubAccessToken + "?" + param.Encode()
    req, _ := http.NewRequest(http.MethodGet, api, nil)
    result, err := o.request(req)
    if utils.ErrNotEmpty(err) {
        return
    }
    accessToken = &AccessToken{
        Token: result["access_token"].(string),
    }
    return

}

// AccessToken2UserInfo 通过accessToken获取用户信息
func (o *OAuthGithub) AccessToken2UserInfo(accessToken string) (userinfo *UserInfo, err error) {
    req, _ := http.NewRequest(http.MethodGet, githubUser, nil)
    req.Header.Set("Authorization", "token "+accessToken)
    result, err := o.request(req)
    if utils.ErrNotEmpty(err) {
        return
    }
    userinfo = &UserInfo{
        Username: result["login"].(string),
        FirstId:  fmt.Sprintf("%.0f", result["id"]),
        SecondId: result["node_id"].(string),
        Nickname: result["name"].(string),
        Avatar:   result["avatar_url"].(string),
        HomePage: result["html_url"].(string),
        Origin:   result,
    }
    if email, ok := result["email"]; ok && !utils.IsNull(email) {
        userinfo.Email = email.(string)
    }
    return userinfo, nil
}

// Username2Userinfo github通过使用username获取用户信息
func (o *OAuthGithub) Username2Userinfo(username string) (userinfo *UserInfo, err error) {
    req, _ := http.NewRequest(http.MethodGet, githubUsers+"/"+username, nil)
    result, err := o.request(req)
    if utils.ErrNotEmpty(err) {
        return
    }
    userinfo = &UserInfo{
        Username: result["login"].(string),
        FirstId:  fmt.Sprintf("%.0f", result["id"]),
        SecondId: result["node_id"].(string),
        Nickname: result["name"].(string),
        Avatar:   result["avatar_url"].(string),
        HomePage: result["html_url"].(string),
        Origin:   result,
    }
    return userinfo, nil
}

func (o *OAuthGithub) request(req *http.Request) (map[string]interface{}, error) {
    client := http.DefaultClient
    client.Timeout = timeout * time.Second
    req.Header.Set("Accept", "application/vnd.github+json")
    resp, err := client.Do(req)
    if utils.ErrNotEmpty(err) {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        err = fmt.Errorf("request [%s] failure http code: %d", req.URL.String(), resp.StatusCode)
        return nil, err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var tmp map[string]interface{}
    _ = json.Unmarshal(body, &tmp)
    if _, ok := tmp["error"]; ok {
        err = fmt.Errorf("request [%s] failure error: %s reason: %s", req.URL.String(), tmp["error"], tmp["error_description"])
        return nil, err
    }
    return tmp, nil
}
