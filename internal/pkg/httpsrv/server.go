package httpsrv

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"goapp/internal/pkg/config"
	"goapp/internal/pkg/watcher"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
	Problem #3 :
		Using the 'gorilla/csrf' package to protect against CSRF attacks in this program. This package provides CSRF protection middleware
		for Go web applications. This will automatically generate and validate CSRF tokens for our forms, ensuring that requests are legit.
*/

// Server is the HTTP server.
type Server struct {
	strChan      <-chan string                   // String channel.
	server       *http.Server                    // Gorilla HTTP server.
	watchers     map[string]*watcher.Watcher     // Counter watchers (k: counterId).
	watchersLock *sync.RWMutex                   // Counter lock.
	sessionStats []sessionStats                  // Session stats.
	quitChannel  chan struct{}                   // Quit channel.
	running      sync.WaitGroup                  // Running goroutines.
	csrfProtect  func(http.Handler) http.Handler // CSRF protection middleware.
}

func New(strChan <-chan string) *Server {
	s := Server{}
	s.strChan = strChan
	s.server = nil // Set below.
	s.watchers = make(map[string]*watcher.Watcher)
	s.watchersLock = &sync.RWMutex{}
	s.sessionStats = []sessionStats{}
	s.quitChannel = make(chan struct{})
	s.running = sync.WaitGroup{}
	s.csrfProtect = csrf.Protect([]byte("32-byte-long-auth-key"))
	return &s
}

func (s *Server) Start() error {
	// Create router.
	r := mux.NewRouter()

	// Register routes.
	for _, route := range s.myRoutes() {
		/*
			Enhancement:
				Using the "ANY" as an HTTP method in route handling seems a bit dangerous, from a security perspective.
				An attacker can use this to exploit the application by sending requests to the server which are not expected.
		*/
		if route.Method == "GET" {
			r.Handle(route.Pattern, route.HFunc).Methods(route.Method)
			if route.Queries != nil {
				r.Handle(route.Pattern, route.HFunc).Methods(route.Method).Queries(route.Queries...)
			}
		} else {
			log.Printf("Unsupported HTTP method: %s for route: %s", route.Method, route.Pattern)
		}
	}

	csrfRouter := s.csrfProtect(r)

	/*
		Enhancement:
			Reading configuration from the config file and setting the HTTP server address and port accordingly.
	*/
	httpHost, httpPort := config.GetConfig()
	httpAddr := httpHost + httpPort

	// Debug.
	log.Printf("HTTP server listening on %s\n", httpAddr)

	// Create HTTP server.
	s.server = &http.Server{
		Addr:         httpAddr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      handlers.CombinedLoggingHandler(os.Stdout, csrfRouter),
	}

	// Start HTTP server.
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	s.running.Add(1)
	go s.mainLoop()

	return nil
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("error: %v\n", err)
	}

	close(s.quitChannel)
	s.running.Wait()
}

func (s *Server) mainLoop() {
	defer s.running.Done()

	for {
		select {
		case str := <-s.strChan:
			s.notifyWatchers(str)
		case <-s.quitChannel:
			return
		}
	}
}

/*
Problem #1:

	Using a pointer for increasing the elements of the original slice.

Also, improving code's readability and sequence of the Server's methods, by grouping them together here,
instead of leaving the incStats() method in watcher.go.
*/
func (s *Server) incStats(id string) {
	// Find and increment.
	for i := range s.sessionStats {
		if s.sessionStats[i].id == id {
			s.sessionStats[i].inc()
			return
		}
	}
	// Not found, add new.
	s.sessionStats = append(s.sessionStats, sessionStats{id: id, sent: 1})
}
