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
    CreateOAuthPasswordAccount(param *params.PutOAuthPassword) (*response.JWTToken, error)
    // OAuth2App 使用第三方OAuth2登录
    OAuth2App(appName, code string) (*response.JWTToken, error)
    // BindEmail 可以修改和创建新账号
    BindEmail(user *model.User, email string) error
    // UnBindEmail 解除绑定邮箱帐号
    UnBindEmail(user *model.User) error
    // BindOAuth2App 绑定第三方OAuth2账号
    BindOAuth2App(user *model.User, appName, code string)
    // UnBindOAuth2App 解除绑定第三方OAuth2账号
    UnBindOAuth2App(user *model.User, appName string) error
    // SetPassword 更新密码 先要判断是否绑定邮箱
    SetPassword(user *model.User, password string) error
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
func (o *oauthService) CreateOAuthPasswordAccount(param *params.PutOAuthPassword) (res *response.JWTToken, err error) {
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
    var userinfo *oauth2.UserInfo
    userinfo, err = o.oauth2AppCode2UserInfo(appName, code)
    if utils.ErrNotEmpty(err) {
        return
    }
    var token *model.Token
    platformCode := model.Name2PlatformCode(appName)
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
                token.Platform = uint8(platformCode)
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

// oauth2AppCode2UserInfo 通过传入对应的oauth2应用名和code换取用户信息
func (o *oauthService) oauth2AppCode2UserInfo(appName, code string) (userinfo *oauth2.UserInfo, err error) {
    var appOAuth oauth2.IOAuth2
    var ok bool
    if appOAuth, ok = o.oauth2[appName]; !ok || !app.Config.IsOAuthKeyExists(appName) {
        err = errs.BadRequest("不支持此类app授权登录", nil)
        return
    }
    var accessToken *oauth2.AccessToken
    if accessToken, err = appOAuth.Code2AccessToken(code); utils.ErrNotEmpty(err) {
        err = errs.BadRequest("获取OAuth2访问令牌失败", err)
        return
    }
    if userinfo, err = appOAuth.AccessToken2UserInfo(accessToken.Token); utils.ErrNotEmpty(err) {
        err = errs.BadRequest("获取OAuth2用户信息失败", err)
        return
    }
    return
}

// BindEmail 绑定邮箱
func (o *oauthService) BindEmail(user *model.User, email string) (err error) {
    if user.Email == email {
        return
    }
    _, err = o.getAccountByUsername(email)
    if utils.ErrNotEmpty(err) {
        //邮箱不存在
        if utils.IsRecordNotFound(err) {
            if user.IsBindEmail() {
                //绑定邮箱更新
                err = o.updateUserEmailAccount(user, email)
            } else {
                //从来没有绑定邮箱创建账号和token
                err = o.createUserEmailAccount(user, email)
            }
        }
        if utils.ErrNotEmpty(err) {
            err = errs.DefaultServerError(err)
        }
        return
    }
    //没有返回错误表示邮箱存在
    err = errs.Conflict("邮箱已存在")
    return
}

// createUserEmailAccount 创建用户邮箱账户
func (o *oauthService) createUserEmailAccount(user *model.User, email string) (err error) {
    account := model.NewActivatedAccount(email, app.Config.Crypt.SaltPassword(utils.RandomStr(2, 10)))
    user.Email = email
    //1、开启事务
    err = app.DB.Transaction(func(tx *gorm.DB) (err error) {
        //2、创建account
        err = tx.Create(account).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        //2、更新用户
        err = tx.Save(user).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        //3、创建token
        token := model.NewAccountToken(o.makeRandomToken(account.Username), account.Username, account.ID, user.ID)
        err = tx.Create(token).Error
        return
    })
    return
}

// updateUserEmailAccount 更新用户邮箱账户
func (o *oauthService) updateUserEmailAccount(user *model.User, email string) (err error) {
    account, err := o.getAccountByUsername(user.Email)
    if utils.ErrNotEmpty(err) {
        return
    }
    token, err := o.getPlatformToken(user.Email, model.PlatformAccount)
    if utils.ErrNotEmpty(err) {
        return
    }
    account.Username = email
    token.FromId = email
    //重新生成token
    token.Token = o.makeRandomToken(email)
    user.Email = email
    err = app.DB.Transaction(func(tx *gorm.DB) (err error) {
        err = tx.Save(account).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        err = tx.Save(token).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        err = tx.Save(user).Error
        return
    })
    return
}

// UnBindEmail 解除绑定,如果当前用户只绑定了一个邮箱无法解除绑定
func (o *oauthService) UnBindEmail(user *model.User) (err error) {
    if !user.IsBindEmail() {
        err = errs.BadRequest("未绑定邮箱,无法解绑", nil)
        return
    }
    if !user.IsBindGithub() {
        err = errs.BadRequest("当前只有一种登录方式,无法解绑", nil)
        return
    }
    err = app.DB.Transaction(func(tx *gorm.DB) (err error) {
        //1、删除帐号表
        err = tx.Unscoped().Delete(&model.Account{}, "username = ?", user.Email).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        //2、删除token表
        err = tx.Unscoped().Delete(&model.Token{}, "user_id = ? AND from_id = ?", user.ID, user.Email).Error
        if utils.ErrNotEmpty(err) {
            return
        }
        //3、更新用户表
        return tx.Model(user).Updates(map[string]interface{}{"email": nil}).Error
    })
    if utils.ErrNotEmpty(err) {
        err = errs.DefaultServerError(err)
    } else {
        user.Email = ""
    }
    return err
}

// BindOAuth2App 绑定第三方账号
func (o *oauthService) BindOAuth2App(user *model.User, appName, code string) (err error) {
    var userinfo *oauth2.UserInfo
    switch appName {
    case oauth2.IOAuth2Github:
        if user.IsBindGithub() {
            err = errs.BadRequest("已绑定github请先解除绑定", nil)
            break
        }
        userinfo, err = o.oauth2AppCode2UserInfo(appName, code)
        if user.GithubName == userinfo.Username {
            break
        }
    default:
        err = errs.BadRequest("OAuth2 类型错误", nil)
        break
    }
    if utils.ErrNotEmpty(err) {
        return
    }
    var token *model.Token
    platformCode := model.Name2PlatformCode(appName)
    token, err = o.getPlatformToken(userinfo.FirstId, platformCode)
    if utils.ErrNotEmpty(err) {
        if utils.IsRecordNotFound(err) {
            //1、创建token
            token = model.NewAccountToken(o.makeRandomToken(userinfo.FirstId), userinfo.FirstId, 0, user.ID)
            err = app.DB.Transaction(func(tx *gorm.DB) error {
                token.Platform = uint8(platformCode)
                err = tx.Create(token).Error
                if utils.ErrNotEmpty(err) {
                    return err
                }
                user.UpdateFieldsFromOAuth2Userinfo(userinfo)
                //2、更新用户
                return tx.Save(user).Error
            })
        }
        if utils.ErrNotEmpty(err) {
            err = errs.DefaultServerError(err)
        }
    } else {
        err = errs.Conflict("绑定的账户已被其他人使用,请更换")
    }
    return
}

// UnBindOAuth2App 解除绑定第三方账号,如果当前用户只有一个绑定不允许解绑
func (o *oauthService) UnBindOAuth2App(user *model.User, appName string) (err error) {
    if !user.IsBindEmail() {
        err = errs.BadRequest("当前只有一种登录方式,无法解绑", nil)
        return
    }
    switch appName {
    case oauth2.IOAuth2Github:
        if !user.IsBindGithub() {
            err = errs.BadRequest("没有绑定github无法解除绑定", nil)
            break
        }
        err = app.DB.Transaction(func(tx *gorm.DB) (err error) {
            //1、删除token表
            err = tx.Unscoped().Delete(&model.Token{}, "user_id = ? AND platform = ?", user.ID, model.PlatformGithub).Error
            if utils.ErrNotEmpty(err) {
                return
            }
            //2、更新用户表
            err = tx.Model(user).Updates(map[string]interface{}{"github_name": nil}).Error
            if utils.ErrNotEmpty(err) {
                user.GithubName = ""
            }
            return
        })
        if utils.ErrNotEmpty(err) {
            err = errs.DefaultServerError(err)
        }
        break
    default:
        err = errs.BadRequest("OAuth2 类型错误", nil)
    }
    return
}

// SetPassword 设置密码,先判断是否绑定了账号
func (o *oauthService) SetPassword(user *model.User, password string) (err error) {
    if !user.IsBindEmail() {
        err = errs.BadRequest("您没有绑定邮箱无法设置密码", nil)
        return
    }
    newPassword := app.Config.Crypt.SaltPassword(password)
    err = app.DB.Model(&model.Account{}).
        Where("username = ?", user.Email).
        Update("password", newPassword).
        Error
    return
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
