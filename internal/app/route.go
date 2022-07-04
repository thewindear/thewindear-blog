package app

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/thewindear/thewindear-blog/internal/handle"
)

// 初始化路由和启动服务
func newRoute(conf *conf) (err error) {
    app := fiber.New()
    setRoute(app)
    if err = app.Listen(conf.Server.ListenHost()); err != nil {
        err = fmt.Errorf("启动服务失败: %v", err)
    }
    return
}

// setRoute 设置路径
func setRoute(app *fiber.App) {
    api := app.Group("/api")
    v1 := api.Group("/v1", func(ctx *fiber.Ctx) error {
        ctx.Set("Version", "v1")
        return ctx.Next()
    })
    // 普通帐号密码登录
    v1.Post("/login/oauth/password", handle.LoginOauthPassword)
    // 创建帐号密码账户
    v1.Put("/login/oauth/password", handle.CreateOauthPasswordAccount)
    // 第三方授权登录
    v1.Get("/login/oauth/:app/authorize", handle.LoginOauthAuthorize)
}
