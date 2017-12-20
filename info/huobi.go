// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package info

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"time"
)

var addr = flag.String("addr", "api.huobi.pro", "huobi websocket")
var comm = make(chan map[string]interface{}, 1)

type Pong struct {
	Ts float64 `json:"pong"`
}
type Sub struct {
	Sub string `json:"sub"`
	Id  string `json:"id"`
}
type HuobiDT struct {
	Ch   string     `json:ch`
	Ts   int64      `json:ts`
	Tick HuobiDTone `json:tick`
}
type HuobiDTone struct {
	Id   int64         `json:id`
	Ts   int64         `json:ts`
	Data []HuobiDTone1 `json:data`
}
type HuobiDTone1 struct {
	Id        int64   `json:id`
	Price     float64 `json:price`
	Amount    float64 `json:amount`
	Direction string  `json:"direction"`
	Ts        int64   `json:ts`
}
type KLine struct {
	symbol string
	period string
}
type MarketDepth struct {
	symbol string
	depth  string
}

func HandlerHuobi(info io.Reader, c *websocket.Conn) {
	// io.Copy(os.Stdout, info)
	a, _ := ioutil.ReadAll(info)
	// log.Println(string(a))
	var data map[string]interface{}
	json.Unmarshal(a, &data)
	if _, err := data["ping"]; err {
		ts := data["ping"].(float64)
		ret := &Pong{
			Ts: ts,
		}
		ret1, err := json.Marshal(ret)
		if err != nil {
			log.Println("json encode err:", err)
			return
		}
		err = c.WriteMessage(websocket.TextMessage, ret1)
		if err != nil {
			log.Println("websocket write err:", err)
			return
		}

	}
	if _, err := data["status"]; err {
		comm <- data
		return
	}
	if _, err := data["ch"]; err {
		trade := HuobiDT{}
		json.Unmarshal(a, &trade)
		re := regexp.MustCompile(`^market\.(.*?)\.trade\.detail$`)
		log.Println(re.FindStringSubmatch(trade.Ch)[1])
		for _, v := range trade.Tick.Data {
			log.Println("ts", v.Ts)
			log.Println("price", v.Price)
			log.Println("amount", v.Amount)
			log.Println("direction", v.Direction)
		}
		return
	}

}
func SubHuobi(topic int, v interface{}, c *websocket.Conn) error {
	sub := new(Sub)
	if topic == 1 {
		sub.Id = "id1"
		t := v.(KLine)
		sub.Sub = fmt.Sprintf("market.%s.kline.%s", t.symbol, t.period)
	} else if topic == 2 {
		sub.Id = "id2"
		t := v.(MarketDepth)
		sub.Sub = fmt.Sprintf("market.%s.depth.%s", t.symbol, t.depth)
	} else if topic == 3 {
		sub.Id = "id3"
		sub.Sub = fmt.Sprintf("market.%s.trade.detail", v.(string))
	} else if topic == 4 {
		sub.Id = "id4"
		sub.Sub = fmt.Sprintf("market.%s.detail", v.(string))
	} else {
		return errors.New("invalid topic")
	}
	ret, err := json.Marshal(sub)
	if err != nil {
		log.Println("json parse err:", err)
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, ret)
	if err != nil {
		log.Println("websocket write err:", err)
		return err
	}
	select {
	case data := <-comm:
		if data["status"].(string) == "ok" {
			err = nil
			log.Printf("id:%s,subbed:%s,ts:%f", data["id"].(string), data["subbed"].(string), data["ts"].(float64))
		} else {
			err_info := data["err-code"].(string) + data["err-msg"].(string)
			log.Println(err_info)
			err = errors.New(err_info)
		}
	case <-time.After(2 * time.Second):
		return errors.New("read time out")
	}
	return err
}
func Huobi() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer c.Close()
		defer close(done)
		if err != nil {
			log.Println("zip:", err)
			return
		}

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			reader := bytes.NewReader(message)
			zr, err := gzip.NewReader(reader)
			if err != nil {
				log.Println("zip decompress:", err)
				return
			}
			go HandlerHuobi(zr, c)
			zr.Close()

		}
	}()
	err = SubHuobi(3, "btcusdt", c)
	if err != nil {
		log.Println("sub err:", err)
	}
	for {
		select {
		case <-done:
			log.Println("socket closed")
			return
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			c.Close()
			return
		}
	}
}
