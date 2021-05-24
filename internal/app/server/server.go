package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/harunnryd/tempolalu/internal/app/handler"
	"github.com/harunnryd/tempolalu/internal/pkg/chicustom"
)

type Server interface {
	Router(handler handler.Handler) (w chicustom.Router)
	GetHTTPServer() *http.Server
	GracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool)
}

type server struct {
	Addr         string
	Handler      handler.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewServer(addr string, handler handler.Handler, readTimeout time.Duration, writeTimeout time.Duration, idleTimeout time.Duration) Server {
	return &server{Addr: addr, Handler: handler, ReadTimeout: readTimeout, WriteTimeout: writeTimeout, IdleTimeout: idleTimeout}
}

func (s *server) GetHTTPServer() *http.Server {
	return &http.Server{
		Addr:         s.Addr,
		Handler:      s.Router(s.Handler),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
	}
}

func (s *server) Router(handler handler.Handler) (w chicustom.Router) {
	w = chicustom.NewRouter(chi.NewRouter())
	w.Route("/v1", func(r chi.Router) {
		router := r.(chicustom.Router)
		router.Action(chicustom.NewRest(http.MethodPost, "/transactions", handler.Transaction().CreateTransaction))
		router.Action(chicustom.NewRest(http.MethodDelete, "/transactions/{transaction_id}", handler.Transaction().DeleteTransactionByID))
	})

	return
}

// GracefullShutdown ...
func (s *server) GracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}
