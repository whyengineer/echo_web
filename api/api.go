package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/whyengineer/api.cryptobc.info/caculate"
	"golang.org/x/sync/syncmap"
)

var api *api_type

type api_type struct {
	Db *gorm.DB
	Rc *redis.Client
	Mq map[string]*syncmap.Map
}

func Load(e *echo.Echo) *echo.Group {
	api = new(api_type)
	var err error
	api.Mq = make(map[string]*syncmap.Map)
	api.Db, err = gorm.Open("mysql", "test:12345678@tcp(123.56.216.29:3306)/coins?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panic("open database failed")
	}
	//start redis
	api.Rc = redis.NewClient(&redis.Options{
		Addr:     "123.56.216.29:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = api.Rc.Ping().Result()
	if err != nil {
		log.Panic("connect redis err:", err)
	}
	
	server := NewSocketServer()
	e.GET("/socket.io/", echo.WrapHandler(server))
	log.Println("load api router")
	api := e.Group("/api")
	api.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	api.GET("/gettrade", GetTrade)
	api.GET("/getsenc", GetSecond)
	return api
}
func GetSecond(c echo.Context) error {
	ts, _ := strconv.Atoi(c.QueryParam("ts"))
	coin := c.QueryParam("coin")
	for i := ts; i > ts-50; i-- {
		key := coin + ":" + strconv.FormatInt(int64(i), 10)
		data, err := api.Rc.Get(key).Result()
		if err != nil {
			continue
		}
		return c.String(http.StatusOK, data)
	}
	return c.NoContent(http.StatusNoContent)
}

// /api/gettrade?coin=eosusdt&plat=huobi&year=2018&month=1&day=25&num=5&type=day
func GetTrade(c echo.Context) error {
	//ts, _ := strconv.Atoi(c.QueryParam("ts"))
	coin := c.QueryParam("coin")
	plat := c.QueryParam("plat")
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	day, _ := strconv.Atoi(c.QueryParam("day"))
	timetype := c.QueryParam("type")
	num, _ := strconv.Atoi(c.QueryParam("num"))
	switch timetype {
	case "min1":
		hour, _ := strconv.Atoi(c.QueryParam("hour"))
		min, _ := strconv.Atoi(c.QueryParam("min"))
		var d []caculate.Min1TradeTable
		tmp := fmt.Sprintf("%04d%02d%02d%02d%02d", year, month, day, hour, min)
		key, _ := strconv.Atoi(tmp)
		if api.Db.Where("time_key<=? AND prop=? AND coin_type=?",
			key, plat, coin).Order("time_key desc").
			Limit(num).Find(&d).RecordNotFound() {
			return c.NoContent(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, &d)
		}
	case "min5":
		hour, _ := strconv.Atoi(c.QueryParam("hour"))
		min, _ := strconv.Atoi(c.QueryParam("min"))
		var d []caculate.Min5TradeTable
		tmp := fmt.Sprintf("%04d%02d%02d%02d%02d", year, month, day, hour, min)
		key, _ := strconv.Atoi(tmp)
		if api.Db.Where("time_key<=? AND prop=? AND coin_type=?",
			key, plat, coin).Order("time_key desc").
			Limit(num).Find(&d).RecordNotFound() {
			return c.NoContent(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, &d)
		}
	case "min30":
		hour, _ := strconv.Atoi(c.QueryParam("hour"))
		min, _ := strconv.Atoi(c.QueryParam("min"))
		var d []caculate.Min30TradeTable
		tmp := fmt.Sprintf("%04d%02d%02d%02d%02d", year, month, day, hour, min)
		key, _ := strconv.Atoi(tmp)
		if api.Db.Where("time_key<=? AND prop=? AND coin_type=?",
			key, plat, coin).Order("time_key desc").
			Limit(num).Find(&d).RecordNotFound() {
			return c.NoContent(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, &d)
		}
	case "hour1":
		hour, _ := strconv.Atoi(c.QueryParam("hour"))
		var d []caculate.Hour1TradeTable
		tmp := fmt.Sprintf("%04d%02d%02d%02d", year, month, day, hour)
		key, _ := strconv.Atoi(tmp)
		if api.Db.Where("time_key<=? AND prop=? AND coin_type=?",
			key, plat, coin).Order("time_key desc").
			Limit(num).Find(&d).RecordNotFound() {
			return c.NoContent(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, &d)
		}
	case "hour4":
		hour, _ := strconv.Atoi(c.QueryParam("hour"))
		var d []caculate.Hour4TradeTable
		tmp := fmt.Sprintf("%04d%02d%02d%02d", year, month, day, hour)
		key, _ := strconv.Atoi(tmp)
		if api.Db.Where("time_key<=? AND prop=? AND coin_type=?",
			key, plat, coin).Order("time_key desc").
			Limit(num).Find(&d).RecordNotFound() {
			return c.NoContent(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, &d)
		}
	case "day":
		var d []caculate.DayTradeTable
		tmp := fmt.Sprintf("%04d%02d%02d", year, month, day)
		key, _ := strconv.Atoi(tmp)
		if api.Db.Where("time_key<=? AND prop=? AND coin_type=?",
			key, plat, coin).Order("time_key desc").
			Limit(num).Find(&d).RecordNotFound() {
			return c.NoContent(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, &d)
		}

	}
	return nil
}
