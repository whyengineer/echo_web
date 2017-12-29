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
	CoinType string
	Prop string
	SecInfo map[int32]CalInfo
	Min1Info map[int32]CalInfo
	Min5Info map[int32]CalInfo
}
func (c *CalMachine)GetMin5Info(ts int32){
	var minCalInfo CalInfo
	stopTs:=ts
	startTs:=ts-300
	minCalInfo.Ts=ts
	var j int32
	for i:=startTs;i < stopTs;i+=60{
		val,ok:=c.Min1Info[i/60*60]
		if ok {
			j++
			minCalInfo.BuyAmount+=val.BuyAmount
			minCalInfo.SellAmount+=val.SellAmount
			minCalInfo.Price+=val.Price
		}
	}
	if j!=0{
		minCalInfo.Price/=float64(j)
	}
	
	c.Min5Info[stopTs/300*300]=minCalInfo
}
func (c *CalMachine)GetMin1Info(ts int32){
	var minCalInfo CalInfo
	
	stopTs:=ts
	startTs:=ts-60
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
	if j!=0{
		minCalInfo.Price/=float64(j)
	}	
	c.Min1Info[stopTs/60*60]=minCalInfo
	//log.Println(minCalInfo.Price,minCalInfo.BuyAmount,minCalInfo.SellAmount)
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

	b:=int32(time.Now().Unix())
	a:=b-3600	
	go func(){
		for i:=a;i< b;i++{
			c.GetSecInfo(i)
		}
		log.Println("init seconde data done")
		for i:=a+60;i< b;i+=60{
			c.GetMin1Info(i)
		}
		log.Println("init 1 minute data done")
		for i:=a+300;i< b;i+=300{
			c.GetMin5Info(i)
		}
		log.Println("init 5 minute data done")
		log.Println("strat delete timeout data")
		secondTick:=time.Tick(1*time.Second)
		for{
			<-secondTick
			delete(c.SecInfo,a)
			delete(c.Min1Info,a/60*60)
			delete(c.Min5Info,a/300*300)
		}
		
		
		
		
	}()

	go func(){
		secondTick:=time.Tick(1*time.Second)
		min1Tick:=time.Tick(1*time.Minute)
		min5Tick:=time.Tick(5*time.Minute)
		for{
			select{
			case nowSecTime:=<-secondTick:
				nowTs:=int32(nowSecTime.Unix())
				c.GetSecInfo(nowTs)
				val,ok:=c.SecInfo[nowTs]
				if ok {
					log.Println("second data:",val.BuyAmount,val.SellAmount,val.Price)
				}
			case nowMin1Time:=<-min1Tick:
				nowTs:=int32(nowMin1Time.Unix())
				c.GetMin1Info(nowTs)
				val,ok:=c.Min1Info[nowTs/60*60]
				if ok {
					log.Println("1 minute data:",val.BuyAmount,val.SellAmount,val.Price)
				}
			case nowMin5Time:=<-min5Tick:
				nowTs:=int32(nowMin5Time.Unix())
				c.GetMin1Info(nowTs)
				val,ok:=c.Min5Info[nowTs/300*300]
				if ok {
					log.Println("5 minute data:",val.BuyAmount,val.SellAmount,val.Price)
				}
			}
			
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