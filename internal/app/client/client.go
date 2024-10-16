package client

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func connectToWebSocket(wg *sync.WaitGroup, httpAddr string, connectionNum int) {
	defer wg.Done()

	u := url.URL{Scheme: "ws", Host: httpAddr, Path: "/goapp/ws"}
	log.Printf("Starting WS connection to %s (ID: #%d)", u.String(), connectionNum)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Dial error (connection #%d): %v", connectionNum, err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read error for [conn #%d]: %v", connectionNum, err)
				return
			}
			log.Printf("[conn #%d]: %s", connectionNum, message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Printf("Write error for [conn #%d]: %v", connectionNum, err)
				return
			}
		}
	}
}

// starts the client with the given number of connections, by using goroutines
func StartClient(numConnections int, httpAddr string, wg *sync.WaitGroup) {
	wg.Add(numConnections)

	for i := 0; i < numConnections; i++ {
		go connectToWebSocket(wg, httpAddr, i)
	}
}
