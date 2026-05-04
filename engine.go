package vodka

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// httprouter wrapper
type Engine struct {
	router *httprouter.Router
}

// creates a new router
func New() *Engine {
	return &Engine{
		router: httprouter.New(),
	}
}

// Runs the http server
func (e *Engine) Run(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Listening and serving HTTP on %s\n", addr)

	// Using net/http
	return http.ListenAndServe(addr, e.router)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.GET(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := &Context{
			Writer:  w,
			Request: r,
		}

		handler(c)
	})
}
