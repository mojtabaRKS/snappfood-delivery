package rest

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	r := gin.New()
	r.RedirectTrailingSlash = false

	return &Server{
		engine: r,
	}
}

const headerTimeout = 10 * time.Second

func (s *Server) Serve(ctx context.Context, address string) error {
	srv := &http.Server{
		Addr:              address,
		Handler:           s.engine,
		ReadHeaderTimeout: headerTimeout,
	}

	log.Printf("rest server starting at: %s", address)
	srvError := make(chan error)
	go func() {
		srvError <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Print("rest server is shutting down")
		return srv.Shutdown(ctx)
	case err := <-srvError:
		return err
	}
}
