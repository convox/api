package api

import (
	"net/http"

	"github.com/convox/logger"
	"github.com/gorilla/mux"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, c Context) *Error

type Router struct {
	*mux.Router
	log *Logger
}

func NewRouter() *Router {
	return &Router{
		Router: mux.NewRouter(),
		log:    NewLogger(),
	}
}

func (r *Router) HandleApi(method, path string, fn HandlerFunc) {
	log := r.log.Namespace("method=%q path=%q", method, path)

	r.HandleFunc(path, apiHandler(fn, log)).Methods(method)
}

func (r *Router) HandleText(method, path, text string) {
	r.HandleFunc(path, textHandler(text)).Methods(method)
}

func apiHandler(fn HandlerFunc, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r, log)

		defer r.Body.Close()

		if err := fn(w, r, c); err != nil {
			if err.Server() {
				err.Record()
			}

			http.Error(w, err.Error(), err.Code())
		}
	}
}

func textHandler(s string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}
}
