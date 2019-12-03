package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", func(ctx echo.Context) error {
		tpl, err := template.ParseFiles("template/login.html")
		if err != nil {
			ctx.Logger().Error("parse file error:", err)
			return err
		}

		ctx.Logger().Info("this is login page...")

		data := map[string]interface{}{
			"msg": ctx.QueryParam("msg"),
		}

		if user, ok := ctx.Get("user").(*User); ok {
			data["username"] = user.Username
			data["had_login"] = true
		} else {
			sess := getCookieSession(ctx)
			if flashes := sess.Flashes("username"); len(flashes) > 0 {
				data["username"] = flashes[0]
			}
			sess.Save(ctx.Request(), ctx.Response())
		}

		var buf bytes.Buffer
		err = tpl.Execute(&buf, data)
		if err != nil {
			return err
		}

		return ctx.HTML(http.StatusOK, buf.String())
	})

	// 登录
	e.POST("/login", func(ctx echo.Context) error {
		username := ctx.FormValue("username")
		passwd := ctx.FormValue("passwd")
		rememberMe := ctx.FormValue("remember_me")

		if username == "polaris" && passwd == "123567" {
			// 用标准库种 cookie
			cookie := &http.Cookie{
				Name:     "username",
				Value:    username,
				HttpOnly: true,
			}
			if rememberMe == "1" {
				cookie.MaxAge = 7*24*3600	// 7 天
			}
			ctx.SetCookie(cookie)

			return ctx.Redirect(http.StatusSeeOther, "/")
		}

		// 用户名或密码不对，用户名回填，通过 github.com/gorilla/sessions 包实现
		sess := getCookieSession(ctx)
		sess.AddFlash(username, "username")
		err := sess.Save(ctx.Request(), ctx.Response())
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/?msg="+err.Error())
		}

		return ctx.Redirect(http.StatusSeeOther, "/?msg=用户名或密码错误")
	})

	// 退出登录
	e.GET("/logout", func(ctx echo.Context) error {
		cookie := &http.Cookie{
			Name:    "username",
			Value:   "",
			Expires: time.Now().Add(-1e9),
			MaxAge:  -1,
		}
		ctx.SetCookie(cookie)

		return ctx.Redirect(http.StatusSeeOther, "/")
	})
}

func getCookieSession(ctx echo.Context) *sessions.Session {
	sess, _ := cookieStore.Get(ctx.Request(), "request-scope")
	return sess
}
