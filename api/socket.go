package api


import (
	"log"
	"github.com/googollee/go-socket.io"
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
		// trade,err:=NewTradeApi()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// trade.TradePush("btcusdt","huobi",so)
		// trade.TradePush("ethusdt","huobi",so)
		// so.Join("chat")
		// so.On("chat_message", func(msg string) {
		// 	so.Emit("chat_message", msg)
		// 	log.Println("emit:",msg)
		// 	trade.GetNowInfo("btcusdt","huobi",10)
		// 	so.BroadcastTo("chat", "chat_message", msg)
		// })
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
		
	})

	return server
}