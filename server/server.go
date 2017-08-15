package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pragkent/slackwork/config"
	"github.com/pragkent/slackwork/sink"
)

type Server struct {
	c   *config.Config
	sc  *sink.Controller
	srv *http.Server
}

func New(addr string, c *config.Config) (*Server, error) {
	s := &Server{
		c: c,
	}

	sc, err := sink.NewController(c)
	if err != nil {
		return nil, err
	}

	s.sc = sc

	r := mux.NewRouter()
	s.registerHandlers(r)

	s.srv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return s, nil
}

func (s *Server) Addr() string {
	return s.srv.Addr
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return s.srv.Shutdown(ctx)
}
