package handle

import (
    "github.com/gofiber/fiber/v2"
)

// LoginOauthPassword 通过帐号密码登录
// @Post /api/login/oauth/password
func LoginOauthPassword(c *fiber.Ctx) error {

    return nil
}

// CreateOauthPasswordAccount 创建帐号密码
// @Put /api/login/oauth/password
func CreateOauthPasswordAccount(c *fiber.Ctx) error {

    return nil
}

// LoginOauthAuthorize 通过第三方授权登录
// @Get /api/login/oauth/:app/authorize
func LoginOauthAuthorize(c *fiber.Ctx) error {

    return nil
}
