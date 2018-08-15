package NetEngine

import (
	"ConfigEngine"
	"Utils"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
)

var listenerTag = "LISTENER"

const (
	HTTP_METHOD_GET  string = "GET"
	HTTP_METHOD_POST string = "POST"
	HTTP_METHOD_PUT  string = "PUT"
)

//HTTPListener listener interface
//type HTTPListener func(w http.ResponseWriter, r *http.Request)
type HTTPListener func(param map[string]interface{}) bool

//SetHTTPGetListener set http listener
func SetHTTPGetListener(url string, listener HTTPListener) {
	RegistGetRouter(url, listener)
}

//SetHTTPPostListener set http listener
func SetHTTPPostListener(url string, listener HTTPListener) {
	RegistPostRouter(url, listener)
}

//HTTPListenerEntry http listener
func listenerHTTPReq(w http.ResponseWriter, r *http.Request) bool {
	path := r.URL.Path

	//check whether path include "/"
	if path[0] == '/' && len(path) > 1 {
		path = path[1:]
	}

	//we should check get/put/set
	fmt.Println("req is ", *r)
	switch r.Method {
	case HTTP_METHOD_GET:
		return dealwithGetReq(path, w, r)

	case HTTP_METHOD_POST:
		return dealwithPostReq(path, w, r)

	default:
		return false
	}
}

func dealwithGetReq(url string, w http.ResponseWriter, r *http.Request) (result bool) {
	listener, m, ok := MatchGetURL(url)
	if ok {
		parseResult := parseForm(w, r)
		for key, val := range m {
			parseResult[key] = val
		}

		v := listener(parseResult)
		Utils.LOGD(listenerTag, "result is ", v)
		rr, err := json.Marshal(v)
		Utils.LOGD(listenerTag, "marshel result is ", rr, " error is ", err)
		if err == nil {
			w.Write(rr)
			return true
		}
	}

	return false
}

func dealwithPostReq(url string, w http.ResponseWriter, r *http.Request) (result bool) {
	listener, m, ok := MatchPostURL(url)
	if ok {
		parseResult := parseForm(w, r)
		for key, val := range m {
			parseResult[key] = val
		}

		v := listener(parseResult)
		Utils.LOGD(listenerTag, "result is ", v)
		rr, err := json.Marshal(v)
		Utils.LOGD(listenerTag, "marshel result is ", rr, " error is ", err)
		if err == nil {
			w.Write(rr)
			return true
		}
	}
	return false
}

func parseForm(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	r.ParseMultipartForm(ConfigEngine.GetConfig().MaxFileUpload)

	var param = make(map[string]interface{})

	urlParams := r.URL.Query()
	fmt.Println("urlparams is ", urlParams)

	//we should analyze data

	for urlKey, urlVal := range urlParams {
		param[urlKey] = urlVal
	}

	for formKey, formVal := range r.Form {
		param[formKey] = formVal
	}

	fmt.Println("r.Form is ", r.Form)

	for postFormKey, postFormVal := range r.PostForm {
		param[postFormKey] = postFormVal
	}

	fmt.Println("r.Form is ", r.PostForm)

	//fmt.Println("r.MultipartForm is ", r.MultipartForm)
	if r.MultipartForm != nil {
		for multiValueKey, multiValueVal := range r.MultipartForm.Value {
			param[multiValueKey] = multiValueVal
		}

		for multiPartKey, multiPartVal := range r.MultipartForm.File {
			param[multiPartKey] = multiPartVal
		}
	}

	for key, head := range r.Header {
		param[key] = head
	}

	return param
}

//AcceptPostFile file up load
func AcceptPostFile(filehead *multipart.FileHeader, buffersize int, path string) {

	file, err := filehead.Open()
	if err == nil {
		fileObj, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)

		defer func() {
			fileObj.Close()
		}()

		for {
			buf := make([]byte, 0, buffersize)
			readsize, readerr := file.Read(buf[0 : buffersize-1])

			if readerr != nil && readsize == 0 {
				break
			}

			fileObj.Write(buf[0 : readsize-1])
		}
	}
}
