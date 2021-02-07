package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

//存所有的连接
var mapWsconn = make(map[*websocket.Conn]*User)

type User struct {
	UserNmae string
	Pic      string
}

var MapPic = map[int]string{
	0: "images/53f442834079a.jpg",
	1: "images/53f44283a4347.jpg",
	2: "images/hetu.jpg",
	3: "images/11928.gif",
	4: "images/1111020107630.jpg",
	5: "images/1111020331193.jpg",
	6: "images/1111004231713.jpg",
	7: "images/1111003753433.jpg",
	8: "images/1111004953355.jpg",
}

//处理请求
func Handler(ws *websocket.Conn) {
	if _, ok := mapWsconn[ws]; !ok {
		mapWsconn[ws] = new(User)
	}
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}
		ReArr := strings.Split(reply, "=")
		if ReArr[0] == "@@:" {
			var Picture = ""
			for key, _ := range MapPic {
				if va, ok := MapPic[key]; ok {
					Picture = va
					delete(MapPic, key)
					break
				}

			}
			if Picture == "" {
				Picture = "images/11928.gif"
			}
			user := &User{
				UserNmae: ReArr[1],
				Pic:      Picture,
			}
			mapWsconn[ws] = user
			//更新头像
			UpPic(ws, Picture)
			//提示谁上线
			UserOnline(ws, ReArr[1])

		} else if ReArr[0] == "qq:" {
			TN := time.Now().Format("2006-01-02 15:04:05")
			writeBack := "reply" + "&%&" + ReArr[1] + "&%&" + mapWsconn[ws].Pic + "&%&" + mapWsconn[ws].UserNmae + "&%&" + TN
			for con, _ := range mapWsconn {
				if err := websocket.Message.Send(con, writeBack); err != nil {
					DeleteUser(con)
				}
			}
		}
		//向所有客户端发送通知
		SendOnlin()
	}
}

//浏览器更新在线人员
func SendOnlin() {
	//剔除非在线人员
	for con, _ := range mapWsconn {
		if err := websocket.Message.Send(con, "test"); err != nil {
			DeleteUser(con)
		}

	}
	var bB bytes.Buffer
	bB.WriteString("Line&%&")
	bB.WriteString(" <li class='fn-clear selected'><em>所有用户</em></li>")
	for _, user := range mapWsconn {
		if user.UserNmae != "" {
			bB.WriteString("<li class='fn-clear' data-id='1'><span>")
			bB.WriteString("<img src=" + user.Pic + " width='30' height='30'  alt=''/>")
			bB.WriteString("</span><em>" + user.UserNmae + "</em>")
			bB.WriteString("<small class='online' title='在线'></small></li>")
		}
	}
	writeBack := bB.String()
	for con, _ := range mapWsconn {
		if err := websocket.Message.Send(con, writeBack); err != nil {
			DeleteUser(con)
		}

	}
}

//更新头像
func UpPic(ws *websocket.Conn, pic string) {
	writeBack := "up&%&" + pic
	if err := websocket.Message.Send(ws, writeBack); err != nil {
		DeleteUser(ws)
	}
}

//网站标志
func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./views/favicon.png")
}

//删除不在线用户
func DeleteUser(ws *websocket.Conn) {
	username := mapWsconn[ws].UserNmae
	delete(mapWsconn, ws)
	//提示谁下线了
	writeBack := "offline&%&" + username + "已下线……"

	for con, _ := range mapWsconn {
		if err := websocket.Message.Send(con, writeBack); err != nil {
			DeleteUser(con)
		}

	}
}

//广播某某上线
func UserOnline(ws *websocket.Conn, username string) {
	writeBack := "Oneline&%&" + username + "已上线…………"
	for con, _ := range mapWsconn {
		if err := websocket.Message.Send(con, writeBack); err != nil {
			DeleteUser(con)
		}

	}
}
