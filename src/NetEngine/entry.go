package NetEngine

import (
	"fmt"
	"net/http"
)

//HTTPMainEntry entry
func HTTPMainEntry(w http.ResponseWriter, r *http.Request) {
	//we need check whether it is a template get request
	fmt.Println("HTTPMainEntry,path is ", r.URL)
	if templateHTTPReq(w, r) {
		return
	} else if listenerHTTPReq(w, r) {
		return
	}
}

//Init init
func Init() {
	initTemplate()
	initWebSocketserver()

	http.HandleFunc("/", HTTPMainEntry)
}
