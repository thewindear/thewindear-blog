package service

import (
    "fmt"
    "github.com/golang-jwt/jwt/v4"
    "github.com/thewindear/thewindear-blog/internal/app"
    "github.com/thewindear/thewindear-blog/internal/data/params"
    "github.com/thewindear/thewindear-blog/internal/data/response"
    "github.com/thewindear/thewindear-blog/internal/data/response/errs"
    "github.com/thewindear/thewindear-blog/internal/model"
    "github.com/thewindear/thewindear-blog/internal/utils"
    "github.com/thewindear/thewindear-blog/pkg/oauth2"
    "gorm.io/gorm"
    "time"
)

type IOAuthService interface {
    // OAuthLoginPassword 通过帐号密码登录
    OAuthLoginPassword(param *params.OAuthPassword) (*response.JWTToken, error)
    // CreateOAuthPasswordAccount 通过账号密码创建账户
    CreateOAuthPasswordAccount(param *params.CreateOAuthPassword) (*response.JWTToken, error)
    // OAuth2App 使用第三方OAuth2登录
    OAuth2App(appName, code string) (*response.JWTToken, error)
}

type oauthService struct {
    oauth2 map[string]oauth2.IOAuth2
}

// OAuthLoginPassword 通过帐号密码获取账户信息并生成jwt token返回
func (o *oauthService) OAuthLoginPassword(param *params.OAuthPassword) (res *response.JWTToken, err error) {
    var account *model.Account
    account, err = o.getAccountByUsername(param.Username)
    if !utils.IsNull(err) {
        if utils.IsRecordNotFound(err) {
            err = errs.Unauthorized("账号不存在")
            return
        }
        err = errs.DefaultServerError(err)
        return
    }
    if !param.CheckPassword(account) {
        err = errs.Unauthorized("密码错误")
        return
    }
    if !account.IsActivated() {
        err = errs.StatusForbidden("账号被未被激活")
        return
    }
    var token *model.Token
    token, err = o.getAccountToken(account.ID)
    if !utils.IsNull(err) {
        if utils.IsRecordNotFound(err) {
            err = errs.StatusForbidden("token未生成")
            return
        }
        err = errs.DefaultServerError(err)
        return
    }
    return o.MakeJWTToken(token)
}

// CreateOAuthPasswordAccount 通过用户名密码创建账户
func (o *oauthService) CreateOAuthPasswordAccount(param *params.CreateOAuthPassword) (res *response.JWTToken, err error) {
    var account *model.Account
    account, err = o.getAccountByUsername(param.Username)
    if utils.IsNull(err) && !utils.IsNull(account) {
        err = errs.Conflict("账号已存在")
        return
    }
    if !utils.IsRecordNotFound(err) {
        err = errs.DefaultServerError(err)
        return
    }
    account = model.NewActivatedAccount(param.Username, app.Config.Crypt.SaltPassword(param.Password))
    var user = model.NewDefaultUser(param.Username)
    var token *model.Token
    //1、开启事务
    err = app.DB.Transaction(func(tx *gorm.DB) (err error) {
        //2、创建account
        err = tx.Create(account).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        //3、创建user
        err = tx.Create(user).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        //4、创建token
        token = model.NewAccountToken(o.makeRandomToken(param.Username), param.Username, account.ID, user.ID)
        err = tx.Create(token).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        return nil
    })
    //5、生成jwt
    return o.MakeJWTToken(token)
}

// OAuth2App 使用第三方OAuth2登录
func (o *oauthService) OAuth2App(appName, code string) (res *response.JWTToken, err error) {
    var appOAuth oauth2.IOAuth2
    var ok bool
    if appOAuth, ok = o.oauth2[appName]; !ok || !app.Config.IsOAuthKeyExists(appName) {
        err = errs.BadRequest("不支持此类app授权登录", nil)
        return
    }
    var accessToken *oauth2.AccessToken
    if accessToken, err = appOAuth.Code2AccessToken(code); utils.ErrNotEmpty(err) {
        err = errs.BadRequest("获取授权码失败", err)
        return
    }
    var userinfo *oauth2.UserInfo
    if userinfo, err = appOAuth.AccessToken2UserInfo(accessToken.Token); utils.ErrNotEmpty(err) {
        err = errs.BadRequest("获取用户信息失败", err)
        return
    }
    var token *model.Token
    token, err = o.getPlatformToken(userinfo.FirstId, model.Name2PlatformCode(appName))
    if utils.ErrNotEmpty(err) {
        if utils.IsRecordNotFound(err) {
            //创建用户
            user := model.NewOAuth2User(userinfo)
            err = app.DB.Transaction(func(tx *gorm.DB) (err error) {
                //1、创建用户
                err = tx.Create(user).Error
                if utils.ErrNotEmpty(err) {
                    return
                }
                //2、创建token
                token = model.NewAccountToken(o.makeRandomToken(userinfo.FirstId), userinfo.FirstId, 0, user.ID)
                err = tx.Create(token).Error
                if utils.ErrNotEmpty(err) {
                    return
                }
                return nil
            })
        }
        if utils.ErrNotEmpty(err) {
            err = errs.DefaultServerError(err)
            return
        }
    }
    return o.MakeJWTToken(token)
}

// MakeJWTToken 通过Token模型生成JWTToken
func (o *oauthService) MakeJWTToken(token *model.Token) (res *response.JWTToken, err error) {
    if !token.IsEnable() {
        err = errs.StatusForbidden("token不可用")
        return
    }
    jwtToken, err := o.makeJWT(token)
    if utils.ErrNotEmpty(err) {
        err = errs.DefaultServerError(err)
        return
    }
    res = &response.JWTToken{
        Token: jwtToken,
    }
    return
}

// makeJWT 生成jwt
func (o *oauthService) makeJWT(token *model.Token) (string, error) {
    claims := jwt.RegisteredClaims{
        Issuer:    app.Config.Application.Domain,
        Subject:   "token",
        ID:        token.Token,
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        NotBefore: jwt.NewNumericDate(time.Now()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(app.Config.Application.TokenExpireSeconds())),
    }
    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    jwtTokenStr, err := jwtToken.SignedString([]byte(app.Config.Crypt.Token))
    if err != nil {
        return "", fmt.Errorf("生成token失败: %v", err)
    }
    return jwtTokenStr, nil
}

// makeRandomToken 随机生成一个token
func (o *oauthService) makeRandomToken(fromId string) string {
    return utils.CryptMD5(fromId, utils.MakeToken())
}

// getAccountByUsername 通过用户名查询账户
func (o *oauthService) getAccountByUsername(username string) (account *model.Account, err error) {
    err = app.DB.Model(account).Where("username = ?", username).First(&account).Error
    return
}

// getTokenByAccount 通过账号id查询token
func (o *oauthService) getAccountToken(accountId uint) (token *model.Token, err error) {
    err = app.DB.Model(token).Where("account_id = ? AND platform = ?", accountId, model.PlatformAccount).First(&token).Error
    return
}

// getPlatformToken 通过fromId和平台号查询对应的token
func (o *oauthService) getPlatformToken(fromId string, platform uint) (token *model.Token, err error) {
    err = app.DB.Model(token).Where("from_id = ? AND platform = ?", fromId, platform).First(&token).Error
    return
}

// NewOAuthService 返回IOAuthService实例
func NewOAuthService() IOAuthService {
    return &oauthService{
        oauth2: map[string]oauth2.IOAuth2{
            oauth2.IOAuth2Github: &oauth2.OAuthGithub{
                ClientSecret: app.Config.OAuthKeys[oauth2.IOAuth2Github].ClientSecret,
                ClientId:     app.Config.OAuthKeys[oauth2.IOAuth2Github].ClientId,
            },
        },
    }
}
