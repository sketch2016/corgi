package ConfigEngine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	defaultMaxThread    int = 12
	defaultMaxHtmlCache int = 128
)

const (
	analyseIdle = iota
	analyseServer
	analyseHTMLMap
	analyseHTMLRes
	analyseUser

	analyseSQL
)

const (
	contentRoot      string = "root"
	contentMaxThread string = "maxthread"
	contentCache     string = "cache"
	contentPort      string = "port"
	contentMaxFileUp string = "maxfileupload"
)

const (
	configServerTag  string = "!server"
	configHTMLdirTag string = "!htmlmap"
	configHTMLresTag string = "!htmlres"
	configHTMLuser   string = "!user"
)

var status = analyseIdle

//Config corgi config
type Config struct {
	//server
	ServerPort    int
	MaxFileUpload int64

	//static config
	HTMLMapRoot string
	HTMLMap     map[string]string
	HTMLCache   int

	HTMLResRoot string

	//user
	UserMaxThread int
}

var gConfig = Config{
	ServerPort:    80,
	MaxFileUpload: 1024 * 1024,
	HTMLMapRoot:   "",
	HTMLMap:       make(map[string]string),
	HTMLCache:     defaultMaxHtmlCache,
	HTMLResRoot:   "",
	UserMaxThread: defaultMaxThread,
}

var path = "config.pro"

//LoadConfig load from config.pro
func LoadConfig() error {

	//os.Create("3.txt")
	fmt.Println("StartAnalyze")

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("err is ", err)
		return err
	}

	defer file.Close()

	br := bufio.NewReader(file)
	for {
		bytes, _, c := br.ReadLine()
		fmt.Println("bytes is ", bytes)
		if c == io.EOF {
			fmt.Println("read eof")
			break
		}
		//fmt.Println(string(a))
		//setValue
		val := string(bytes)
		val = strings.Replace(val, " ", "", -1)
		val = filterComment(val)

		fmt.Println("val is ", val)
		if len(val) == 0 {
			continue
		}

		switch val {
		case configHTMLdirTag:
			status = analyseHTMLMap
			continue

		case configHTMLresTag:
			status = analyseHTMLRes
			continue

		case configHTMLuser:
			status = analyseUser

		case configServerTag:
			status = analyseServer

		default:
			//Do nothing
		}

		setValue(val)
	}

	return nil
}

//GetConfig return config
func GetConfig() Config {
	return gConfig
}

func setValue(val string) {
	switch status {
	case analyseHTMLMap:
		setValHTMLMap(val)
	case analyseHTMLRes:
		setValHTMLRes(val)
	case analyseUser:
		setValUser(val)
	case analyseSQL:
		setSQLValue(val)
	case analyseServer:
		setServerValue(val)
	}
}

func setServerValue(val string) {
	v := strings.Split(val, "=")
	if len(v) < 2 {
		return
	}

	switch v[0] {
	case contentPort:
		port, error := strconv.Atoi(v[1])
		if error == nil {
			gConfig.ServerPort = port
		}

	case contentMaxFileUp:
		size, error := strconv.ParseInt(v[1], 10, 64)
		if error == nil {
			gConfig.MaxFileUpload = size
		}
	}
}

func setValUser(val string) {
	v := strings.Split(val, "=")
	if len(v) < 2 {
		return
	}

	switch v[0] {
	case contentMaxThread:
		maxthread, error := strconv.Atoi(v[1])
		if error == nil {
			gConfig.UserMaxThread = maxthread
		}
	}
}

func setValHTMLRes(val string) {
	v := strings.Split(val, "=")
	if len(v) < 2 {
		return
	}

	switch v[0] {
	case contentRoot:
		gConfig.HTMLResRoot = v[1]
	}
}

func setValHTMLMap(val string) {

	v := strings.Split(val, "=")
	if len(v) < 2 {
		return
	}

	fmt.Println("v[0] is ", v[0], "v[1] is ", v[1])
	switch v[0] {
	case contentCache:
		cache, error := strconv.Atoi(v[1])
		if error == nil {
			gConfig.HTMLCache = cache
		}
		return

	case contentRoot:
		gConfig.HTMLMapRoot = v[1]
		return

	default:
		//donothing
	}

	gConfig.HTMLMap[v[0]] = v[1]
}

//setSQLValue set sql config information
func setSQLValue(val string) {

}

func filterComment(val string) string {
	index := strings.Index(val, "#")
	if index > 0 {
		return val[0:index]
	}
	return val
}
