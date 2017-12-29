package api


import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	cal "github.com/whyengineer/echo_web/caculate"
	"net/http"
	"log"
	"strconv"
)


var cm *cal.CalMachine

func Load(e *echo.Echo) *echo.Group{
	server:=NewSocketServer()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins: []string{"*"},
		AllowOrigins: []string{"http://localhost", "http://localhost:8080", "http://localhost:1323"},
	}))
	e.GET("/socket.io/",echo.WrapHandler(server))
	log.Println("load admin router")
	api:=e.Group("/api")
	api.GET("/test",func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	api.GET("/gettrade",GetTradeInfo)
	var err error
	cm,err=cal.NewCalMachine("btcusdt","huobi")
	if err!=nil{
		log.Println(err)
	}
	cm.StartCal()
	return api
}

func GetTradeInfo(c echo.Context) error {
	timeType := c.QueryParam("timetype")
	ts,_:=strconv.Atoi(c.QueryParam("ts"))
	if timeType=="second"{
		val,ok:=cm.SecInfo[int32(ts)]
		if ok{
			return c.JSON(http.StatusOK,&val)
		}else{
			empty:=&cal.CalInfo{}
			return c.JSON(http.StatusOK,empty)
		}
	}
	if timeType=="minute1"{
		val,ok:=cm.Min1Info[int32(ts)/60*60]
		if ok{
			return c.JSON(http.StatusOK,&val)
		}else{
			empty:=&cal.CalInfo{}
			return c.JSON(http.StatusOK,empty)
		}
	}
	if timeType=="minute5"{
		val,ok:=cm.Min5Info[int32(ts)/300*300]
		if ok{
			return c.JSON(http.StatusOK,&val)
		}else{
			empty:=&cal.CalInfo{}
			return c.JSON(http.StatusOK,empty)
		}
	}
	return nil
}