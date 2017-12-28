package api


import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"fmt"
	"time"
)




func Load(e *echo.Echo) *echo.Group{
	server:=NewSocketServer()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins: []string{"*"},
		AllowOrigins: []string{"http://localhost", "http://localhost:8080", "http://localhost:1323"},
	}))
	e.GET("/socket.io/",echo.WrapHandler(server))
	fmt.Println("load admin router")
	api:=e.Group("/api")
	api.GET("/test",func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	var err error
	calM,err=NewCalMachine("ethusdt","huobi")
	if err!=nil{
		fmt.Println(err)
	}
	stop:=int32(time.Now().Unix())
	start:=stop-60
	calRes,err:=calM.CalData(start,stop)
	if err!=nil{
		fmt.Println(err)
	}
	for _,v:=range calRes{
		fmt.Println("buyamount:",v.BuyAmount,"sellamount:",v.SellAmount,"price:",v.Price)
	}
	return api
}