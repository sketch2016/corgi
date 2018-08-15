package NetEngine

import (
	"ConfigEngine"
	"sync"
)

var sessionMap = make(map[string]*Session)
var rwLock sync.RWMutex
var once sync.Once

//Session one group means one user
type Session struct {
	IP          string
	Port        int
	NetIdentity string
	ThreadToken chan int
}

//GetSession get group session
func GetSession(netIdentity string) (r *Session) {

	rwLock.RLock()
	var ret = sessionMap[netIdentity]
	rwLock.RUnlock()

	return ret
}

//CreatSession create group session
func CreatSession(ip string) (r *Session) {
	rwLock.Lock()
	val := new(Session)
	val.IP = ip
	val.Port = 0 //no use
	threadNum := ConfigEngine.GetConfig().UserMaxThread
	val.ThreadToken = make(chan int, threadNum)
	rwLock.Unlock()
	return val
}

//AcquireThreadToken one client can only use limited thread resouce
func (p *Session) AcquireThreadToken() {
	p.ThreadToken <- 1
}

//ReleaseThreadToken release thread resource
func (p *Session) ReleaseThreadToken() {
	<-p.ThreadToken
}
