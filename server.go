package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/whyengineer/echo_web/admin"
	"github.com/whyengineer/echo_web/app"
	// "github.com/gorilla/sessions"
	// "github.com/labstack/echo-contrib/session"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	admin.Load(e)
	app.Load(e)
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
