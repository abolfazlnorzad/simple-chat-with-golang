package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	conn, err := websocket.Dial("ws://localhost:7777", "", CreateDemoIp())
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close()

	go receive(conn)

	send(conn)

}

func CreateDemoIp() string {
	//192.168.1.1
	var arry [4]int
	for i := 0; i < len(arry); i++ {
		rand.Seed(time.Now().UnixNano())
		arry[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arry[0], arry[1], arry[2], arry[3])
}

func receive(conn *websocket.Conn) {
	for {
		var m Message
		err := websocket.JSON.Receive(conn, &m)
		if err != nil {
			log.Fatalln("Error in Recieve Data :", err)
			continue
		}
		fmt.Println("Message from Server :", m.Text)
	}
}

func send(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		m := Message{
			Text: text,
		}
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			fmt.Println("Error in Send Data, ", err)
			continue
		}
	}

}
