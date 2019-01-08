package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TAG ...
var TAG = "MAIN"

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(TAG, "Upgrade", err)
		return
	}
	defer conn.Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(TAG, "ReadMessage", err)
			break
		}
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println(TAG, "WriteMessage", err)
			break
		}
	}
}
