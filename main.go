package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"webchat/controllers"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/echo", websocket.Handler(controllers.Handler))
	http.Handle("/css/", http.FileServer(http.Dir("static")))
	http.Handle("/js/", http.FileServer(http.Dir("static")))
	http.Handle("/images/", http.FileServer(http.Dir("static")))
	http.Handle("/font/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/favicon.png", controllers.Favicon)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/index1.html")
	})

	//IP := GetLocalIP()
	IP:="10.5.164.115"
	if IP == "" {
		panic("Failure to start the service")
	}
	//UpdateWsHost(IP)
	fmt.Println("Listening and Serving on " + IP + ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
func UpdateWsHost(IP string) {
	var viewPath = "./views/index.html"
	f, err := os.Open(viewPath)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	hostLine := ""
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		valueList := strings.Split(string(b), "=")
		if len(valueList) != 2 {
			continue
		} else if valueList[0] == "var wsuri " {
			hostLine = "var wsuri =" + valueList[1]
			break
		}
	}
	newHostLine := "var wsuri = \"ws://" + IP + ":8080/echo\";"
	TemFile, _ := ioutil.ReadFile(viewPath)
	TemContent := string(TemFile)
	TemContent = strings.Replace(TemContent, hostLine, newHostLine, -1)
	ioutil.WriteFile("./views/index1.html", []byte(TemContent), os.ModeType)
}

//获得本机ip
func GetLocalIP() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return ip
			}

		}
	}
	return ""
}
