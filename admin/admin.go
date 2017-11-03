package admin

import(
	"github.com/labstack/echo"
	"net/http"
	"fmt"
)
//Load admin router
func Load(e *echo.Echo) *echo.Group{
	fmt.Println("load admin router")
	rAdmin:=e.Group("/admin")
	rAdmin.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "admin page")
	})
	return rAdmin
}