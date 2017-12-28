package caculate

import(
	"github.com/go-redis/redis"
	"strconv"
	"encoding/json"
	"time"
	"log"
)

type CoinInfo struct {
	CoinType string  `json:CoinType`
	Amount   float64 `json:Amount`
	Price    float64 `json:"Price`
	Dir      string  `json:Dir`
	Ts       int64   `json:Ts`
	Prop     string  `json:Prop`
}

type CalInfo struct{	
	BuyAmount  float64
	SellAmount float64
	Price float64
	Ts int32
}

type CalMachine struct{
	rclient *redis.Client
	lastTs int32
	nowTs int32
	CoinType string
	Prop string
	SecInfo map[int32]CalInfo
	Min1Info map[int32]CalInfo
	Min5Info map[int32]CalInfo
}
func (c *CalMachine)GetMin5Info(ts int32){
	var minCalInfo CalInfo
	startTs:=(ts/300)*300
	stopTs:=startTs+300
	minCalInfo.Ts=startTs
	var j int32
	for i:=startTs;i < stopTs;i+=60{
		val,ok:=c.SecInfo[i]
		if ok {
			j++
			minCalInfo.BuyAmount+=val.BuyAmount
			minCalInfo.SellAmount+=val.SellAmount
			minCalInfo.Price+=val.Price
		}
	}
	minCalInfo.Price/=float64(j)
	c.Min5Info[startTs]=minCalInfo
}
func (c *CalMachine)GetMin1Info(ts int32){
	var minCalInfo CalInfo
	startTs:=(ts/60)*60
	stopTs:=startTs+60
	minCalInfo.Ts=startTs
	var j int32
	for i:=startTs;i < stopTs;i++{
		val,ok:=c.SecInfo[i]
		if ok {
			j++
			minCalInfo.BuyAmount+=val.BuyAmount
			minCalInfo.SellAmount+=val.SellAmount
			minCalInfo.Price+=val.Price
		}
	}
	minCalInfo.Price/=float64(j)
	c.Min1Info[startTs]=minCalInfo
	log.Println(minCalInfo.Price,minCalInfo.BuyAmount,minCalInfo.SellAmount)
}
func (c *CalMachine)GetSecInfo(ts int32){
	var singleV CoinInfo
	var value []CoinInfo
	var single_cal CalInfo
	single_cal.Ts=ts
	key:=c.CoinType+":"+c.Prop+":"+strconv.FormatInt(int64(ts),10)
	val, err := c.rclient.Get(key).Result()
	if err == redis.Nil{
		return 
	}else if err!=nil{
		return
	}else{
		//log.Println(val)
		err:=json.Unmarshal([]byte(val),&singleV)
		if err !=nil{
			err=json.Unmarshal([]byte(val), &value)
			if err!=nil{
				log.Println(err)
			}
		}else{
			single_cal.Price=singleV.Price
			if singleV.Dir=="buy"{
				single_cal.BuyAmount=singleV.Amount
			}else{
				single_cal.SellAmount=singleV.Amount
			}
			c.SecInfo[ts]=single_cal
			return
		}
	}
	for _,val:=range value{
		single_cal.Price+=val.Price
		if val.Dir=="buy"{
			single_cal.BuyAmount+=val.Amount
		}else if val.Dir=="sell"{
			single_cal.SellAmount+=val.Amount
		}
	}
	single_cal.Price/=float64(len(value))
	c.SecInfo[ts]=single_cal
}

func (c* CalMachine)StartCal(){
	//start init the ram data
	c.nowTs=int32(time.Now().Unix())
	c.lastTs=c.nowTs-3600

	a:=c.lastTs
	b:=c.nowTs
	go func(){
		for i:=a;i< b;i++{
			c.GetSecInfo(i)
		}
		log.Println("init seconde data done")
	}()

	go func(){
		tick:=time.Tick(1*time.Second)
		for{
			<-tick
			c.GetSecInfo(c.nowTs)
			_,ok:=c.SecInfo[c.nowTs]
			if ok {
				// log.Println(val.BuyAmount,val.SellAmount,val.Price)
			}
			c.nowTs++
			delete(c.SecInfo,c.lastTs)
			c.lastTs++
		}
	}()
}
func NewCalMachine(cointype string,prop string)(*CalMachine,error){
	c:=new(CalMachine)
	c.CoinType=cointype
	c.Prop=prop
	c.SecInfo=make(map[int32]CalInfo)
	c.Min1Info=make(map[int32]CalInfo)
	c.Min5Info=make(map[int32]CalInfo)
	
	//start redis
	c.rclient = redis.NewClient(&redis.Options{
		Addr:     "123.56.216.29:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := c.rclient.Ping().Result()
	if err!=nil{
		return nil,err
	}
	return c,nil
}