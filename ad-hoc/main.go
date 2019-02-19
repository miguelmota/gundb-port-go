package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/miguelmota/gundb-port-go/dup"
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
	d := dup.NewDup()

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
					msg := make(map[string]interface{})
					id := strconv.Itoa(count)
					msg["#"] = d.Track(id)

					js, err := json.Marshal(msg)
					if err != nil {
						log.Fatal(err)
					}

					// write message to peer
					if err = peer.WriteMessage(websocket.TextMessage, js); err != nil {
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

			var js map[string]interface{}
			err = json.Unmarshal(msg, &js)
			if err != nil {
				log.Fatal(err)
			}

			id := js["#"].(string)

			// comment out this line to test
			if d.Check(id) {
				continue
			}

			d.Track(id)
			fmt.Printf("received: %s\n", js)
			fmt.Printf("from: %s\n", peer.RemoteAddr())
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

/* // BROWSER!
var peer = new WebSocket('ws://localhost:8080')
peer.onopen = function(o){ console.log('open', o) }
peer.onclose = function(c){ console.log('close', c) }
peer.onmessage = function(m){
	var msg = JSON.parse(m.data)
	console.log(msg)
	peer.send(m.data)
};
peer.onerror = function(e){ console.log('error', e) }
*/
