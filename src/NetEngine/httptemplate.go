package NetEngine

import (
	"ConfigEngine"
	"Utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var templateTag = "Template"

//Template null struct
type Template struct{}

var lConfig *ConfigEngine.Config

var contentCache *Utils.LRUCache

func initTemplate() {
	contentCache = Utils.CreateLRUCache(128)
}

//HTTPReq httprequest
func templateHTTPReq(writer http.ResponseWriter, request *http.Request) (ret bool) {
	path := request.URL.Path
	//check whether path include "/"
	if path[0] == '/' && len(path) > 1 {
		path = path[1:]
	}

	fmt.Println("templateHTTPReq path is ", path)

	if lConfig == nil {
		t := ConfigEngine.GetConfig()
		lConfig = &t
	}

	html, ok := ConfigEngine.GetConfig().HTMLMap[path]
	if ok {
		fullpath := ConfigEngine.GetConfig().HTMLMapRoot + "/" + html
		val := templateReadFile(fullpath)
		writer.Write(val)
		return true
	}

	//check whether it is a js or img
	resRoot := ConfigEngine.GetConfig().HTMLResRoot + "/" + path

	val := templateReadFile(resRoot)
	if val != nil {
		//we should check the file type
		requestType := path[strings.LastIndex(path, "."):]

		switch requestType {
		case ".css":
			writer.Header().Set("content-type", "text/css")

		case ".js":
			writer.Header().Set("content-type", "text/javascript")

		case ".svg":
			writer.Header().Set("content-type", "image/svg+xml")

		default:
			//nothing
		}

		writer.Write(val)
		return true
	}

	return false
}

func templateReadFile(path string) (ret []byte) {
	//read cache start
	val := contentCache.GetCache(path)
	if val != nil {
		Utils.LOGD(templateTag, "cache hit :", path)
		return val.([]byte)
	}
	//read cache end

	fin, err := os.Open(path)
	defer fin.Close()
	if err != nil {
		fmt.Println("templateReadFile error:", err)
		return nil
	}
	fd, _ := ioutil.ReadAll(fin)

	//add cache start
	contentCache.AddCache(path, fd)
	//add cache end

	return fd
}
