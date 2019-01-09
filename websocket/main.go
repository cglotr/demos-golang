package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// TAG ...
var TAG = "MAIN"

var currID = 0
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

	currID++
	myID := currID

	writeMsgString(conn, "Your ID is "+strconv.Itoa(myID)+".")

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(TAG, "ReadMessage", err)
			break
		}

		if w := writeMsg(conn, mt, msg); !w {
			break
		}
	}
}

func writeMsg(conn *websocket.Conn, msgType int, msg []byte) bool {
	err := conn.WriteMessage(msgType, msg)
	if err != nil {
		log.Println(TAG, "WriteMessage", err)
		return false
	}
	return true
}

func writeMsgString(conn *websocket.Conn, msg string) bool {
	return writeMsg(conn, 1, []byte(msg))
}
