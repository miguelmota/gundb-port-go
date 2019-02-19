package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/miguelmota/gundb-port-go/dup"
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
			//fmt.Printf("received: %s\n", msg)
			//fmt.Printf("from: %s\n", peer.RemoteAddr())

			putValue, ok := js["put"]

			if ok {
				ham.Mix(putValue.(map[string]interface{}), graph)
				fmt.Println("----------------")
				graphJs, err := json.Marshal(graph)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(graphJs))
			}

			for _, peer := range peers {
				if err = peer.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Fatal(err)
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// BROWSER! Use html/index.html
