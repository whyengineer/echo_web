package api


import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	cal "github.com/whyengineer/echo_web/caculate"
	"net/http"
	"log"
	"strconv"
)


var btccm *cal.CalMachine
var ethcm *cal.CalMachine
var eoscm *cal.CalMachine

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
	btccm,err=cal.NewCalMachine("btcusdt","huobi")
	if err!=nil{
		log.Println(err)
	}
	btccm.StartCal()
	ethcm,err=cal.NewCalMachine("ethusdt","huobi")
	if err!=nil{
		log.Println(err)
	}
	ethcm.StartCal()
	eoscm,err=cal.NewCalMachine("eosusdt","huobi")
	if err!=nil{
		log.Println(err)
	}
	eoscm.StartCal()
	return api
}

func GetTradeInfo(c echo.Context) error {
	timeType := c.QueryParam("timetype")
	ts,_:=strconv.Atoi(c.QueryParam("ts"))
	cointype:=c.QueryParam("cointype")
	if cointype=="btcusdt"{
		if timeType=="second"{
			t:=btccm.GetSecInfo(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
		if timeType=="minute1"{
			t:=btccm.GetMin1Info(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
		if timeType=="minute5"{
			t:=btccm.GetMin5Info(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
	}
	if cointype=="ethusdt"{
		if timeType=="second"{
			t:=ethcm.GetSecInfo(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
		if timeType=="minute1"{
			t:=ethcm.GetMin1Info(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
		if timeType=="minute5"{
			t:=ethcm.GetMin5Info(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
	}
	if cointype=="eosusdt"{
		if timeType=="second"{
			t:=eoscm.GetSecInfo(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
		if timeType=="minute1"{
			t:=eoscm.GetMin1Info(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
		if timeType=="minute5"{
			t:=eoscm.GetMin5Info(int32(ts))
			return c.JSON(http.StatusOK,&t)
		}
	}
	return nil
}