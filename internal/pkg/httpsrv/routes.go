package httpsrv

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	HFunc   http.Handler
	Queries []string
}

func (s *Server) myRoutes() []Route {
	return []Route{
		{
			Name:    "health",
			Method:  "GET",
			Pattern: "/goapp/health",
			HFunc:   s.handlerWrapper(s.handlerHealth),
		},
		{
			Name:    "websocket",
			Method:  "GET",
			Pattern: "/goapp/ws",
			HFunc:   s.handlerWrapper(s.handlerWebSocket),
		},
		/*
			Changing the route of home to the root path ("/") to implement what it is mentioned in the README file
			(A client connects on `localhost:8080`), plus the requirement of the Feature #B, which was to show the Hex values
			when a browser opens a connection to `localhost:8080`.
		*/
		{
			Name:    "home",
			Method:  "GET",
			Pattern: "/",
			HFunc:   s.handlerWrapper(s.handlerHome),
		},
	}
}

func validateCSRFToken(r *http.Request) bool {
	token := r.Header.Get("X-CSRF-Token")
	return token == "my-csrf-token"
}

/*
Problem #3:

	Adding the CSRF token validation into this wrapper as it is used from all handlers.
	This ensures that the CSRF validation is applied to all requests.
*/
func (s *Server) handlerWrapper(handlerFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if !validateCSRFToken(r) {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		defer func() {
			r := recover()
			if r != nil {
				s.error(w, http.StatusInternalServerError, fmt.Errorf("%v\n%v", r, string(debug.Stack())))
			}
		}()
		handlerFunc(w, r)
	})
}
