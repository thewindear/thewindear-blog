package oauth2

type IOAuth2 interface {
    // RedirectUri 生成跳转至授权uri
    RedirectUri(callback, state string) string
    // Code2AccessToken 通过授权的code返回对应 AccessAccessToken
    Code2AccessToken(code string) (*AccessToken, error)
    // AccessToken2UserInfo 通过accessToken获取用户信息
    AccessToken2UserInfo(accessToken string) (*UserInfo, error)
    // Username2Userinfo 通过用户名获取用户信息
    Username2Userinfo(username string) (*UserInfo, error)
}

// AccessToken 用户令牌
type AccessToken struct {
    Token        string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpireIn     uint   `json:"expire_in"`
}

// UserInfo 用户基本信息
type UserInfo struct {
    Username string
    FirstId  string
    SecondId string
    Nickname string
    Avatar   string
    HomePage string
    Email    string
    // 原始字段
    Origin map[string]interface{}
}
