package httpsrv

import (
	"fmt"
	"goapp/internal/pkg/watcher"
)

func (s *Server) addWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	s.watchers[w.GetWatcherId()] = w
}

func (s *Server) removeWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	// Print statistics before removing watcher.
	for i := range s.sessionStats {
		if s.sessionStats[i].id == w.GetWatcherId() {
			s.sessionStats[i].print()
		}
	}
	// Remove watcher.
	delete(s.watchers, w.GetWatcherId())
}

func (s *Server) notifyWatchers(str string) {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	// Send message to all watchers and increment stats.
	for id := range s.watchers {
		/*
			Problem #1:
				Starting a new goroutine for each watcher to send the message and inc the stats.
				Now each watcher will process the message concurrently.
				Important Note: The watchersLock mutex ensures thread safe access to the map 'watchers'.
		*/
		go func(id string) {
			fmt.Printf("Sending string msg %s, to watcher id: %s\n", str, id)
			s.watchers[id].Send(str)
			s.incStats(id)
		}(id)
	}
}
