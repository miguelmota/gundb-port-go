package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/miguelmota/gundb-port-go/dup"
	"github.com/miguelmota/gundb-port-go/get"
	"github.com/miguelmota/gundb-port-go/ham"
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
	var peers []*websocket.Conn
	d := dup.NewDup()
	graph := make(map[string]interface{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		peer, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		peers = append(peers, peer)

		for {
			// read message from browser
			_, msg, err := peer.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			var js map[string]interface{}
			err = json.Unmarshal(msg, &js)
			if err != nil {
				log.Fatal(err)
			}

			id := js["#"].(string)

			if d.Check(id) {
				continue
			}

			d.Track(id)
			fmt.Printf("received: %s\n", js)
			fmt.Printf("from: %s\n", peer.RemoteAddr())

			putValue, ok := js["put"]
			if ok {
				ham.Mix(putValue.(map[string]interface{}), graph)
			}

			getValue, ok := js["get"]
			if ok {
				ack := get.Get(getValue.(map[string]interface{}), graph)
				if ack != nil {
					j, err := json.Marshal(map[string]interface{}{
						"#":   d.Track(d.Random()),
						"@":   id,
						"put": ack,
					})
					if err != nil {
						log.Fatal(err)
					}

					emit(peers, j)
				}
			}

			emit(peers, msg)

		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func emit(peers []*websocket.Conn, msg []byte) {
	for _, peer := range peers {
		if err := peer.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Fatal(err)
		}
	}
}

// BROWSER! Use html/index.html
