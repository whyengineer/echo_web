package api

import(
	"github.com/go-redis/redis"
	"strconv"
	"encoding/json"
	"time"
	"github.com/googollee/go-socket.io"
)
type CoinInfo struct {
	CoinType string
	Amount   float64
	Price    float64
	Dir      string
	Ts       int64
	Prop     string
}

type Trade struct{
	rclient *redis.Client
}

func (t* Trade)TradePush(cointype string,prop string,so socketio.Socket){
	go func(){
		for{
			data,err:=t.GetNowInfo(cointype,prop,3)
			if err!=nil{
				return
			}
			b,err:=json.Marshal(data)
			if err!=nil{
				return
			}
			err=so.Emit(cointype,string(b))
			if err!=nil{
				return
			}
			time.Sleep(3*time.Second)
		}
	}()
}
//get period trade info
func (t *Trade)GetNowInfo(cointype string,prop string,period int32)([]CoinInfo,error){
	nowts:=int32(time.Now().Unix())
	var w []CoinInfo
	for i:=nowts;i>(nowts-period);i--{
		tmp,err:=t.GetTsInfo(cointype,prop,int64(i))
		if err!=nil{
			return w,err
		}
		w=append(w,tmp...)
	}
	return w,nil
}
func (t *Trade)GetTsInfo(cointype string,prop string,ts int64)([]CoinInfo,error){
	
	var value []CoinInfo
	key:=cointype+":"+prop+":"+strconv.FormatInt(ts,10)
	val, err := t.rclient.Get(key).Result()
	if err == redis.Nil{
		return nil,nil
	}else if err!=nil{
		return nil,err
	}else{
		//log.Println(val)
		json.Unmarshal([]byte(val), &value)
		return value,nil
	}
}
func NewTradeApi()(*Trade,error){
	a:=new(Trade)
	//start redis
	a.rclient = redis.NewClient(&redis.Options{
		Addr:     "123.56.216.29:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := a.rclient.Ping().Result()
	if err!=nil{
		return nil,err
	}
	return a,nil
}