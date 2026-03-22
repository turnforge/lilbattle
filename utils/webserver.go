package utils

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/panyam/servicekit/middleware"
)

type WebAppServer struct {
	Address       string
	GrpcAddress   string
	AllowLocalDev bool
}

func (s *WebAppServer) StartWithHandler(ctx context.Context, handler http.Handler, srvErr chan error, stopChan chan bool) error {
	if s.AllowLocalDev {
		PrintStartupMessage(s.Address)
	} else {
		log.Println("Starting http web server on: ", s.Address)
	}
	var g middleware.Guard
	g.Use(middleware.RequestLogger())
	if s.AllowLocalDev {
		// nil OriginChecker = allow all origins (dev mode)
		g.Use(middleware.CORS(nil))
	}
	handler = g.Wrap(handler)
	server := &http.Server{
		Addr:        s.Address,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Handler:     handler,
	}

	go func() {
		<-stopChan
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalln(err)
			panic(err)
		}
	}()
	srvErr <- server.ListenAndServe()
	return nil
}

