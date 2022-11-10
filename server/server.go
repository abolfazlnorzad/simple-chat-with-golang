package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func main() {
	h := createHub()
	//server := http.Server{
	//	Addr: ":7777",
	//	Handler: websocket.Handler(func(conn *websocket.Conn) {
	//		fmt.Println("dddd")
	//		wsHandler(conn, h)
	//	}),
	//}

	http.Handle("/", websocket.Handler(func(conn *websocket.Conn) {
		wsHandler(conn, h)
	}))

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		return
	}

	//err := server.ListenAndServe()
	//if err != nil {
	//	fmt.Println("moshkel dar serve server")
	//	return
	//}
}

func wsHandler(conn *websocket.Conn, h *hub) {
	go h.run()
	h.AddClientChan <- conn
	for {
		var m Message
		err := websocket.JSON.Receive(conn, &m)
		if err != nil {
			fmt.Println(err)
			h.RemoveClientChan <- conn
			continue
		}
		h.Broadcast <- m
	}
}

type Message struct {
	Text string `json:"text"`
}

type hub struct {
	Clients          map[string]*websocket.Conn
	AddClientChan    chan *websocket.Conn
	RemoveClientChan chan *websocket.Conn
	Broadcast        chan Message
}

func createHub() *hub {
	return &hub{
		Clients:          make(map[string]*websocket.Conn),
		AddClientChan:    make(chan *websocket.Conn),
		RemoveClientChan: make(chan *websocket.Conn),
		Broadcast:        make(chan Message),
	}
}

func (h *hub) run() {
	for {
		select {
		case con := <-h.AddClientChan:
			h.addClient(con)
		case con := <-h.RemoveClientChan:
			h.removeClient(con)
		case m := <-h.Broadcast:
			h.broadcast(m)
		}
	}
}

func (h *hub) addClient(conn *websocket.Conn) {
	h.Clients[conn.RemoteAddr().String()] = conn
}

func (h *hub) removeClient(conn *websocket.Conn) {
	delete(h.Clients, conn.RemoteAddr().String())
}

func (h *hub) broadcast(m Message) {
	for _, conn := range h.Clients {
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
