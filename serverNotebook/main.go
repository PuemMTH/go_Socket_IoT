package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/inancgumus/screen"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	screen.Clear()

	for {

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var message map[string]interface{}
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println(err)
			continue
		}

		// fmt.Printf("Received JSON: %+v\n", message)

		// fmt.Printf("Received message: %s\n", message["message"])

		// Send a response back to the client
		response := map[string]interface{}{
			"status": "OK",
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			continue
		}

		if err := conn.WriteMessage(messageType, responseJSON); err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Message received: %s\n", p)
		// for key, value := range message {
		// 	fmt.Printf("%s: %s\n", key, value)
		// }

	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
