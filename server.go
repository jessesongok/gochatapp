package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func main()  {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	//configure websockets

	http.HandleFunc("/ws", handleConnections)

	//go routine
	go handleMessages()

	log.Println("Http Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err !=nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err :=upgrader.Upgrade(w, r, nil)
	if err !=nil{
		log.Fatal(err)
	}
	defer ws.Close()

	//Register New Clients
	clients[ws] = true

//Infinite loop

	for{
		var msg Message
		//Read in an new message
		err:= ws.ReadJSON(&msg)

		if err !=nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <-msg
	}

}

func handleMessages()  {
	for{
		//Grab the next message from the broadcast channel
		msg := <-broadcast

		//send it out to every client that is currently connected
		for client :=range clients {
			err := client.WriteJSON(msg)
			if err !=nil{
				log.Printf("error: %v", err)
				client.Close()
				delete(clients,client)
			}
		}
	}
}

