package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	port   uint
}

func (s *Server) Start() error {
	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprintf(w, "hello, %s", p)
	})

	return s.listenAndServe(r)
}

func (s *Server) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", s.port)

	s.server = &http.Server{
		Addr:    address,
		Handler: handler,
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener")
	}

	log.Info("Listening HTTP server", zap.String("address", address))
	return s.server.Serve(listener)
}

func New(port uint) *Server {
	return &Server{port: port}
}
