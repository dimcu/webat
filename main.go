package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
	"webat/at"
)

var atConn *at.Serial

var wc *websocket.Conn
var _webPort *string
var _atPort *string

func onRead(buf []byte) {
	if wc == nil {
		return
	}

	if err := websocket.Message.Send(wc, string(buf)); err != nil {
		fmt.Println("Can't send")
		return
	}

}

func onConnect() {
	log.Println("serial onConnect")

	if wc == nil {
		return
	}

	if err := websocket.Message.Send(wc, "serial onConnect"); err != nil {
		fmt.Println("Can't send")
		return
	}
}
func onDisconnect() {

	log.Println("serial onDisconnect")

	if wc == nil {
		return
	}

	if err := websocket.Message.Send(wc, "serial onDisconnect"); err != nil {
		fmt.Println("Can't send")
		return
	}

}

func wsHandler(ws *websocket.Conn) {
	wc = ws
	log.Println(ws.RemoteAddr())
	for {
		if wc == nil {
			return
		}
		var reply string
		if err := websocket.Message.Receive(wc, &reply); err != nil {
			log.Println(err)
			return
		}

		log.Println(reply)
		atConn.SendAtCmd(reply + "\r\n")

	}
}

func atInit() {

	var err error
	atConn, err = at.OpenAtSerial(*_atPort, 115200)
	if err != nil {
		return
	}
	atConn.SetCallback(onRead, onConnect, onDisconnect)
	time.Sleep(time.Millisecond * 100)

}

func main() {

	_webPort = flag.String("webPort", "8833", "web端口")
	_atPort = flag.String("atPort", "/dev/ttyUSB2", "at端口")

	flag.Parse()

	atInit()

	fsh := http.FileServer(http.Dir("html"))

	http.Handle("/", http.StripPrefix("/", fsh))

	http.Handle("/ws", websocket.Handler(wsHandler))

	if err := http.ListenAndServe(":"+*_webPort, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
