package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type PriceTable struct {
	Coin string
}

//123.56.216.29
func DbConnect() {
	db, err := gorm.Open("mysql", "test:12345678@tcp(123.56.216.29:3306)/coins?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("open database failed", err)
		return
	}
	defer db.Close()
}
