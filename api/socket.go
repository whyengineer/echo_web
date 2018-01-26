package api

import (
	"log"

	"golang.org/x/sync/syncmap"

	"github.com/googollee/go-socket.io"
	"github.com/whyengineer/api.cryptobc.info/market"
)

// var (
// 	upgrader = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool { return true },
// 	}
// )

func NewSocketServer() *socketio.Server {
	//create a connect

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		log.Println(so.Id())
		done := make(chan struct{})
		so.On("joincoin", func(msg string) {
			//so.Emit("coin", msg)
			log.Println("join the coin:", msg)
			datac := make(chan market.CoinInfo)
			group, ok := api.Mq[msg]
			if ok {
				_, in := group.Load(so.Id())
				if in {
					log.Println("the id has stored")
				} else {
					group.Store(so.Id(), datac)
				}
			} else {
				t := new(syncmap.Map)
				api.Mq[msg] = t
				t.Store(so.Id(), datac)
			}
			go func() {
				select {
				case data := <-datac:
					log.Println(data)
				case <-done:
					log.Println("bye bye")
					return
				}
			}()
		})
		so.On("quitcoin", func(msg string) {
			a, ok := api.Mq[msg]
			if ok {
				a.Delete(so.Id())
			} else {
				log.Println("the coin not in the map")
			}
			log.Println("quit the coin:", msg)
			close(done)
			//it's ok,beacaue no rt wait the channel,if has a rt wait the channel ,then the rt is block always!!!
			done = make(chan struct{})
		})
		so.On("disconnection", func() {
			for _, val := range api.Mq {
				val.Delete(so.Id())
			}
			log.Println("on disconnect")
			close(done)
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)

	})

	return server
}
