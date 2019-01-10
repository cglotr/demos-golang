package main

import (
	"log"
	"net/http"

	"github.com/cglotr/websocket"
)

var broadcastChan = make(chan string)
var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/chat", chat)
	go broadcast()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func broadcast() {
	for {
		msg := <-broadcastChan

		for conn := range clients {
			err := conn.WriteMessage(1, []byte(msg))
			if err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()
	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			delete(clients, conn)
			break
		}

		broadcastChan <- string(msg[:len(msg)])
	}
}
