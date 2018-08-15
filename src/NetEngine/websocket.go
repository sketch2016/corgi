package NetEngine

import (
	"Utils"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

var webSocketTag = "WEB SOCKET"

var wsToken = "token"

//WebSocketClient web socket server struct
type WebSocketClient struct {
	Conn        *websocket.Conn
	IP          string
	Port        int
	NetIdentity string
}

//WebSocketClientMgr ws mgr struct
type WebSocketClientMgr struct {
	WebSocketClients map[string]*WebSocketClient
	maplock          sync.RWMutex
}

var wsMgr *WebSocketClientMgr
var initMutex sync.RWMutex

func (p *WebSocketClient) startWsListener(ch chan int, handler WebsocketHandler) {
	for {
		var reply string
		var err = websocket.Message.Receive(p.Conn, &reply)
		if err != nil {
			ch <- 1
			return
		}

		handler(reply)
	}
}

func websocketEntry(conn *websocket.Conn) {
	param := conn.Request().URL.Query()
	tokenarr, ok := param[wsToken]
	if !ok || len(tokenarr) < 1 {
		return
	}
	token := tokenarr[0]
	//check whether client already exits
	if wsMgr.isExist(token) {
		return
	}

	//TODO
	//check whether token is valid

	remoteAddr := conn.Request().RemoteAddr
	v := strings.Split(remoteAddr, ":")
	var client = WebSocketClient{
		Conn: conn,
		IP:   v[0],
	}
	port, error := strconv.Atoi(v[1])
	if error == nil {
		client.Port = port
	}

	//start analyse token
	client.NetIdentity = token
	wsMgr.addClient(&client)
	var waitChan = make(chan int)

	path := conn.Request().URL.Path
	//check whether path include "/"
	if path[0] == '/' && len(path) > 1 {
		path = path[1:]
	}

	vfun, ok := wsHandlerMap[path]
	if ok {
		go client.startWsListener(waitChan, vfun)
		<-waitChan
	}

	Utils.LOGD(webSocketTag, "websocket client disconnect : ", client.NetIdentity)
}

//CreateWebSocketClientMgr create a ws manager
func CreateWebSocketClientMgr() *WebSocketClientMgr {
	initMutex.RLock()
	if wsMgr != nil {
		initMutex.RUnlock()
		return wsMgr
	}
	initMutex.RUnlock()
	initMutex.Lock()
	wsMgr = new(WebSocketClientMgr)
	wsMgr.WebSocketClients = make(map[string]*WebSocketClient)
	initMutex.Unlock()
	return wsMgr
}

func (p *WebSocketClientMgr) sendMsg(ip string, msg string) {
	v, ok := wsMgr.WebSocketClients[ip]
	if ok {
		websocket.Message.Send(v.Conn, msg)
	}
}

func (p *WebSocketClientMgr) sendGroupMsg(msg string) {
	for _, client := range p.WebSocketClients {
		websocket.Message.Send(client.Conn, msg)
	}
}

func (p *WebSocketClientMgr) addClient(client *WebSocketClient) {
	p.maplock.Lock()
	p.WebSocketClients[client.NetIdentity] = client
	p.maplock.Unlock()
}

func (p *WebSocketClientMgr) removeClient(client *WebSocketClient) {
	p.maplock.Lock()
	delete(p.WebSocketClients, client.NetIdentity)
	p.maplock.Unlock()
}

func (p *WebSocketClientMgr) isExist(token string) bool {
	p.maplock.RLock()
	_, ok := p.WebSocketClients[token]
	p.maplock.RUnlock()

	return ok
}

func initWebSocketserver() {
	wsMgr = CreateWebSocketClientMgr()
	wsHandlerMap = make(map[string]WebsocketHandler)
}

//WebsocketHandler handler
type WebsocketHandler func(message string)

var wsHandlerMap map[string]WebsocketHandler
var wsHandlerMapLock sync.RWMutex

//SetWebSocketListener set listner
func SetWebSocketListener(wspath string, handler WebsocketHandler) {
	wsHandlerMapLock.Lock()
	wsHandlerMap[wspath] = handler
	wsHandlerMapLock.Unlock()

	http.Handle("/"+wspath, websocket.Handler(websocketEntry))
}
