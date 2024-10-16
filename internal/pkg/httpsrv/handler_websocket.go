package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"goapp/internal/pkg/watcher"

	"github.com/gorilla/websocket"
)

// Handler for the WebSocket
func (s *Server) handlerWebSocket(w http.ResponseWriter, r *http.Request) {

	// Not implemented yet.
	// Validate CSRF token.
	// if !validateCSRFToken(r) {
	// 	http.Error(w, "Invalid CSRF token", http.StatusForbidden)
	// 	return
	// }

	// Create and start a watcher.
	var watch = watcher.New()
	if err := watch.Start(); err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to start watcher: %w", err))
		return
	}
	defer watch.Stop()

	s.addWatcher(watch)
	defer s.removeWatcher(watch)

	// Start WS.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
			// Security Suggestion: We should never blindly trust 'any Origin' by return 'true' here.
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}
	// defer func() { _ = c.Close() }()
	defer c.Close()

	log.Printf("websocket started for watcher %s\n", watch.GetWatcherId())
	defer func() {
		log.Printf("websocket stopped for watcher %s\n", watch.GetWatcherId())
	}()

	/*
		Problem #2 & New Feature #3:
			By using goroutines to read and write messages, we are able to handle multiple clients concurrently.
			This will allow us to read and write messages from the client without blocking the main thread.

			By adding this WaitGroup as well, we are able to wait for all goroutines to finish before closing the connection.
			Through this approach, we make sure that all concurrent operations are properly managed and that the function does
			not exit sooner than it should, as it would might leave resources uncleaned.
	*/
	var wg sync.WaitGroup

	// Read done.
	readDoneCh := make(chan struct{})

	// All done.
	doneCh := make(chan struct{})

	defer close(doneCh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(readDoneCh)
		for {
			select {
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
						log.Printf("failed to read message: %v\n", err)
					}
					return
				}
				var m watcher.CounterReset
				if err := json.Unmarshal(message, &m); err != nil {
					log.Printf("failed to unmarshal message: %v\n", err)
					continue
				}
				watch.ResetCounter()
			case <-doneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()

	wg.Add(1)
	// Starting a goroutine to write as well.
	go func() {
		defer wg.Done()
		for {
			select {
			case cv := <-watch.Recv():
				data, err := json.Marshal(cv)
				if err != nil {
					log.Printf("failed to marshal message: %v\n", err)
					continue
				}
				err = c.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("failed to write message: %v\n", err)
					}
					return
				}
			case <-readDoneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()

	wg.Wait()

}
