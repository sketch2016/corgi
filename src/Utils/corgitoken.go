package Utils

import (
	"ConfigEngine"
	"net"
	"strconv"
	"strings"
	"sync"
)

var usernum = 0
var lock sync.Mutex
var ipAddr *string
var port = -1

//GenerateCorgiToken generate a unique token
func GenerateCorgiToken() string {
	//t := time.Now()
	//timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	if ipAddr == nil {
		addrSlice, _ := net.InterfaceAddrs()
		for _, addr := range addrSlice {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if nil != ipnet.IP.To4() {
					ip := ipnet.IP.String()
					ip = strings.Replace(ip, ".", "", -1)
					ipAddr = &ip
					break
				}
			}
		}
	}

	if port == -1 {
		port = ConfigEngine.GetConfig().ServerPort
	}

	lock.Lock()
	usernum++
	lock.Unlock()

	return strconv.Itoa(port) + *ipAddr + strconv.Itoa(usernum)

}
