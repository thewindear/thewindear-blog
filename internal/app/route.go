package app

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/requestid"
    jwtware "github.com/gofiber/jwt/v3"
    "github.com/golang-jwt/jwt/v4"
    "github.com/thewindear/thewindear-blog/internal/data/response/errs"
    "github.com/thewindear/thewindear-blog/internal/handle"
)

// 初始化路由和启动服务
func newRoute(conf *conf) (err error) {
    app := fiber.New(fiber.Config{
        AppName: Config.Application.Name,
    })
    app.Use(requestid.New())
    app.Use(logger.New())
    setRoute(app)
    if err = app.Listen(conf.Application.ListenHost()); err != nil {
        err = fmt.Errorf("启动服务失败: %v", err)
    }
    return
}

// setRoute 设置路径
func setRoute(app *fiber.App) {
    v1 := app.Group("/api/v1", func(ctx *fiber.Ctx) error {
        ctx.Set("Version", "v1")
        return ctx.Next()
    })
    // 普通帐号密码登录
    v1.Post("/login/oauth/password", handle.OAuthPassword)
    // 创建帐号密码账户
    v1.Put("/login/oauth/password", handle.CreateOAuthPassword)
    // 第三方授权登录
    v1.Get("/login/oauth/:app/authorize", handle.OAuthAuthorize)

    //var jwtMiddleware = jwtware.New(makeJWTConfig())

}

// makeJWTConfig 生成jwt中间件配置
func makeJWTConfig() jwtware.Config {
    return jwtware.Config{
        ContextKey: "token",
        SigningKey: []byte(Config.Crypt.Token),
        Claims:     jwt.RegisteredClaims{},
        ErrorHandler: func(ctx *fiber.Ctx, err error) error {
            if err.Error() == "Missing or malformed JWT" {
                return errs.Unauthorized("token不能为空")
            }
            return errs.Unauthorized("token无效")
        },
    }
}
