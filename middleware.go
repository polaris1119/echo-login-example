package main

import "github.com/labstack/echo/v4"

// AutoLogin 如果上次记住了，则自动登录
func AutoLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("username")
		if err == nil && cookie.Value != "" {
			// 实际项目这里可以通过 username 读库获取用户信息
			user := &User{Username: cookie.Value}

			// 放入 context 中
			ctx.Set("user", user)
		}

		return next(ctx)
	}
}
