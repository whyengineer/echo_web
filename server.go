package main

import (
	"os"
	"os/signal"

	"github.com/labstack/echo"
	"github.com/whyengineer/echo_web/admin"
	"github.com/whyengineer/echo_web/api"
	"github.com/whyengineer/echo_web/app"
	// "github.com/gorilla/sessions"
	// "github.com/labstack/echo-contrib/session"
)

func main() {
	e := echo.New()
	//e.Use(middleware.Logger())
	admin.Load(e)
	app.Load(e)
	api.Load(e)
	e.Start("127.0.0.1:1323")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}
