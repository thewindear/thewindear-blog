package handle

import (
    "github.com/gofiber/fiber/v2"
)

// OAuthPassword 通过帐号密码登录
// @Post /api/login/oauth/password
func OAuthPassword(c *fiber.Ctx) error {

    return nil
}

// CreateOAuthPassword 创建帐号密码
// @Put /api/login/oauth/password
func CreateOAuthPassword(c *fiber.Ctx) error {

    return nil
}

// OAuthAuthorize 通过第三方授权登录
// @Get /api/login/oauth/:app/authorize
func OAuthAuthorize(c *fiber.Ctx) error {

    return nil
}
