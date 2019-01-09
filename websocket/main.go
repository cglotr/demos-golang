package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// TAG ...
var TAG = "MAIN"

var ch = make(chan []byte)
var clients = 0
var currID = 0
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func broadcast(conn *websocket.Conn, myID int) {
	for {
		m := <-ch
		msg := string(m[:len(m)])
		senderID, err := strconv.Atoi(strings.Split(msg, ":")[0])

		if err != nil {
			continue
		}

		if myID != senderID {
			writeMsgString(conn, msg)
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	clients++

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(TAG, "Upgrade", err)
		return
	}
	defer conn.Close()
	defer func() { clients-- }()

	currID++
	myID := currID

	writeMsgString(conn, "Your ID is "+strconv.Itoa(myID)+".")
	go broadcast(conn, myID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(TAG, "ReadMessage", err)
			break
		}

		m := []byte{}
		m = append(m, []byte(strconv.Itoa(myID)+": ")...)
		m = append(m, msg...)

		for i := 0; i < clients; i++ {
			ch <- m
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
