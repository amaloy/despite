package ws

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

var thechan chan net.Conn

func echoHandler(ws *websocket.Conn) {
	thechan <- ws
	fmt.Println("echoHandler")
	buff := make([]byte, 1024)
	_, err := ws.Read(buff)
	if err != nil {
		panic(err)
	}
	fmt.Println(buff)
	ws.Write(buff)
	ws.Close()
}

// Listen - Start listening for http and websocket connections
func Listen(addr string, webRoot string) {
	thechan = make(chan net.Conn)
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(webRoot)))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
