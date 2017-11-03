package app


import(
	"github.com/labstack/echo"
	"path/filepath"
	"fmt"
	//"net/http"

)
// Load application router
func Load(e *echo.Echo) *echo.Group{
	fmt.Println("load app router")
	rApp:=e.Group("")
	path,_:=filepath.Abs(".")
	staticPath:=filepath.Join(path,"static")
	fmt.Println("app path:",staticPath)	
	rApp.Static("/static",staticPath)
	rApp.GET("/", index)
	return rApp
}
func index(c echo.Context)error{
	// return c.File("./front_end/index.html")
	// path,_:=filepath.Abs(".")
	// fmt.Println("index file path",path)
	return c.File("./app/front_end/index.html")
	//return c.String(http.StatusOK, "Hello, World!")
}