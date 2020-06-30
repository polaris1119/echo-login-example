package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// 这里的 studyecho ，实际应用，不应该写在代码中，应该写在配置或环境变量中
var cookieStore = sessions.NewCookieStore([]byte("studyecho"))

func init() {
	rand.Seed(time.Now().UnixNano())

	os.Mkdir("log", 0755)
}

func main() {
	// 创建 echo 实例
	e := echo.New()

	// 配置日志
	configLogger(e)

	// 注册静态文件路由
	e.Static("img", "img")
	e.File("/favicon.ico", "img/favicon.ico")

	// 设置中间件
	setMiddleware(e)

	// 注册路由
	RegisterRoutes(e)

	// 启动服务
	e.Logger.Fatal(e.Start(":2020"))
}

func configLogger(e *echo.Echo) {
	// 定义日志级别
	e.Logger.SetLevel(log.INFO)
	// 记录业务日志
	echoLog, err := os.OpenFile("log/echo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	// 同时输出到文件和终端
	e.Logger.SetOutput(io.MultiWriter(os.Stdout, echoLog))
}

func setMiddleware(e *echo.Echo) {
	// access log 输出到文件中
	accessLog, err := os.OpenFile("log/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	// 同时输出到终端和文件
	middleware.DefaultLoggerConfig.Output = accessLog
	e.Use(middleware.Logger())

	// 自定义 middleware
	e.Use(AutoLogin)

	e.Use(middleware.Recover())
}
