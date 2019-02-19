package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  0,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	Subprotocols:      nil,
	Error:             nil,
	CheckOrigin:       checkOrigin,
	EnableCompression: false,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		peer, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		count := 0
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					count++
					msg := []byte(fmt.Sprintf("hello world: %v", count))
					// write message to peer
					if err = peer.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Fatal(err)
					}
				}
			}
		}()

		for {
			// read message from browser
			_, msg, err := peer.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			data := string(msg)
			fmt.Printf("received: %s\n", data)
			fmt.Printf("from: %s\n", peer.RemoteAddr())
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

/* // BROWSER!
var peer = new WebSocket('ws://localhost:8080')
peer.onopen = function(o){ console.log('open', o) }
peer.onclose = function(c){ console.log('close', c) }
peer.onmessage = function(m){ console.log(m.data) }
peer.onerror = function(e){ console.log('error', e) }
*/
