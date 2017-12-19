// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package info

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "api.huobi.pro", "huobi websocket")

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
			fmt.Println(zr)
			log.Println("recv:", zr.Header.Comment, zr.Header.Extra, zr.Header.Name, zr.Header.OS)
			zr.Close()
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			// if err != nil {
			// 	log.Println("write:", err)
			// 	return
			// }
			log.Println(t)
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
