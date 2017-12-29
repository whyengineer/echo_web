package caculate

import(
	"testing"
	"log"
	"time"
)
func Test_caculate(t *testing.T){
	cm,err:=NewCalMachine("ethusdt","huobi")
	if err!=nil{
		log.Println(err)
		return
	}
	cm.StartCal()
	for{
		//cm.GetMin1Info(1514455800)
		//cm.GetMin5Info(1514455800)
		val,time.Sleep(5*time.Second)
	}
}