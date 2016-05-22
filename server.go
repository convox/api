package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type Server struct {
	n   *negroni.Negroni
	log *Logger
}

func NewServer() Server {
	server := Server{
		n:   negroni.New(),
		log: NewLogger(),
	}

	server.Use(server.log)

	return server
}

func (s *Server) Listen(addr string) {
	s.log.Logf("listen=%q", addr)

	if err := http.ListenAndServe(addr, s.n); err != nil {
		s.log.Error(err)
	}
}

func (s *Server) Use(fn negroni.Handler) {
	s.n.Use(fn)
}

func (s *Server) UseHandler(fn http.Handler) {
	s.n.UseHandler(fn)
}
