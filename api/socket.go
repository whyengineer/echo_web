package api


import (
	"log"
	"github.com/googollee/go-socket.io"
	cal "github.com/whyengineer/echo_web/caculate"
	"encoding/json"
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
		dataC:=make(chan cal.CalInfo)
		cal.NoticeJoin(so.Id(),dataC)
		done:=make(chan struct{})
		go func(){
			for{
				select{
				case data:=<-dataC:
					//log.Println(so.Id(),data.CoinType,data.Price,data.BuyAmount,data.SellAmount)
					as,_:=json.Marshal(&data)
					so.Emit(data.CoinType,string(as))
				case <-done:
					log.Println("bye bye")
					return
				}	
			}
		}()
		

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
			close(done)
			log.Println("on disconnect")
			cal.NoticeQuit(so.Id())
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
		
	})

	return server
}