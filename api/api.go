package api


import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"fmt"
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
	return api
}