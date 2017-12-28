package api

import(
)

var calM *CalMachine
var SecondData []CalInfo
var Min1Data []CalInfo
var Min5Data []CalInfo
var Min10Data []CalInfo


type CalMachine struct{
	CoinTrade *Trade
	CoinType string
	Prop string
	
}
type CalInfo struct{
	
	BuyAmount  float64
	SellAmount float64
	Price float64
	Ts int64
}
func NewCalMachine(cointype string,prop string) (*CalMachine,error){
	a:=new(CalMachine)
	a.CoinType=cointype
	a.Prop=prop
	var err error
	a.CoinTrade,err=NewTradeApi()
	if err !=nil{
		return nil,err
	}
	return a,nil
}
//get every second data
func (c *CalMachine)CalData(start int32,stop int32) ([]CalInfo,error) {
	var calres []CalInfo
	for i:=start ;i < stop;i++{
		rawinfo,err:=c.CoinTrade.GetTsInfo(c.CoinType,c.Prop,int64(i))
		if err!=nil{
			return calres,err 
		}
		if rawinfo==nil{
			continue
		}
		var single_cal CalInfo
		single_cal.Ts=int64(i)
		for _,val:=range rawinfo{
			single_cal.Price+=val.Price
			if val.Dir=="buy"{
				single_cal.BuyAmount+=val.Amount
			}else if val.Dir=="sell"{
				single_cal.SellAmount+=val.Amount
			}
		}
		single_cal.Price/=float64(len(rawinfo))
		calres=append(calres,single_cal)
	}
	return calres,nil
}