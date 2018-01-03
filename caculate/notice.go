package caculate

import (
	"sync"
	"errors"
)



var CalQueue *NoticeQueue



type NoticeQueue struct{
	queue map[string]chan CalInfo
	sync.Mutex
}
func StartNoticeQueue(){
	CalQueue=new(NoticeQueue)
	CalQueue.queue=make(map[string]chan CalInfo)
}
func realSend(c chan CalInfo,data CalInfo){
	c<-data
}
func NoticeSend(a CalInfo){
	CalQueue.Lock()
	defer CalQueue.Unlock()
	for _,val:=range CalQueue.queue{
		go realSend(val,a)
	}
}
func NoticeJoin(id string,qq chan CalInfo)error{
	CalQueue.Lock()
	defer CalQueue.Unlock()
	if _,ok:=CalQueue.queue[id];ok{
		return errors.New("the id exited")
	}else{
		CalQueue.queue[id]=qq
		return nil
	}
}
func NoticeQuit(id string){
	CalQueue.Lock()
	defer CalQueue.Unlock()
	delete(CalQueue.queue,id)
}