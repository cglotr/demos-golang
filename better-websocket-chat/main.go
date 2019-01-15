package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/cglotr/websocket"
	"github.com/gorilla/mux"
)

var broadcastChan = make(chan string, 2)
var clients = make(map[*websocket.Conn]string)
var upgrader = websocket.Upgrader{}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/chat", chat)

	go broadcast()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func broadcast() {
	for {
		name := <-broadcastChan
		msg := <-broadcastChan

		for conn := range clients {
			if clients[conn] == name {
				continue
			}

			err := conn.WriteMessage(1, []byte(name+": "+msg))
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
	clients[conn] = ""

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			delete(clients, conn)
			break
		}

		msgStr := string(msg[:len(msg)])

		if m, _ := regexp.MatchString(":name\\s.*", msgStr); m {
			name := strings.Split(msgStr, " ")[1]
			clients[conn] = name
			continue
		}

		if clients[conn] == "" {
			conn.WriteMessage(1, []byte("Tell me your name first."))
			conn.WriteMessage(1, []byte("Use the command :name <your_name>."))
			continue
		}

		broadcastChan <- clients[conn]
		broadcastChan <- msgStr
	}
}
