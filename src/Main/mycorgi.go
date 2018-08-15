package main

import (
	"ConfigEngine"
	"NetEngine"
	"fmt"
	"mime/multipart"
	"net/http"
)

type testStruct struct {
	Val1 int
	Val2 string
	Val3 int
	Val4 []string
}

//control
func handlerhello(param map[string]interface{}) bool {
	fmt.Println("handlerhello!!!")
	for key, val := range param {
		fmt.Println("key is ", key, " val is ", val)
	}

	ll := []string{"a1", "a2", "a3"}
	var b = testStruct{1, "abc", 3, ll}
	fmt.Println(" b is ", b)
	return true
}

func handlerPost(param map[string]interface{}) bool {
	fmt.Println("handlerpost param is ", param)

	paramMap := param
	object, ok := paramMap["file"]
	//var latch sync.WaitGroup

	if ok {
		filearray := object.([]*multipart.FileHeader)
		for _, fileHead := range filearray {
			NetEngine.AcceptPostFile(fileHead, 1024*1024, fileHead.Filename)
		}
	}

	return true
}

func wshandler(msg string) {
	fmt.Println("wshandler is ", msg)
}

func main() {
	//http.HandleFunc("/", handler)
	//http.HandleFunc("/hello.html", handlerhello)
	//HttpEngine.AddMainController()
	ConfigEngine.LoadConfig()

	NetEngine.Init()

	//HttpEngine.SetHTTPListener("abc", handler)
	//HttpEngine.SetHTTPListener("efg", handlerhello)
	//NetEngine.SetHTTPGetListener("abc", handlerhello)
	NetEngine.SetHTTPPostListener("post1", handlerPost)
	NetEngine.SetWebSocketListener("corgi", wshandler)
	NetEngine.SetHTTPGetListener("abc/:name/:age", handlerhello)

	http.ListenAndServe(":8989", nil)
}
