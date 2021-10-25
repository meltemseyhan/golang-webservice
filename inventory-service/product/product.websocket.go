package product

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func productSocketHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer func(c *websocket.Conn) {
		fmt.Printf("closing web socket \n")
		c.Close()
	}(conn)

	// This part is for testing how to read from socket.
	// Through developer tool's console write javascript:
	// let ws = new WebSocket("ws://localhost:5000/websocket")
	// ws.send("Hello Meltem")
	// ws.close()
	go func(c *websocket.Conn) {
		for {
			_, mess, err := c.ReadMessage()
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("received message: %s \n", mess)
		}
	}(conn)

	for {
		products, err := getTopTenProducts()
		if err != nil {
			log.Println(err)
			break
		}
		err = conn.WriteJSON(products)
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(10 * time.Second)
	}
}
