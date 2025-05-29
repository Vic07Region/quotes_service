package server

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quotes_service/internal/handlers"
	"quotes_service/internal/repository/inmemory"
	"quotes_service/pkg/logger"
	"syscall"
	"time"
)

type Server struct {
	storage    *inmemory.InMemory
	handlers   *handlers.Handler
	httpServer *http.Server
	router     *mux.Router
	stopChan   chan struct{}
	logger     *zap.Logger
}

func NewServer(addr string, debug bool) *Server {
	srv := &Server{}

	r := mux.NewRouter()

	srv.logger = logger.NewLogger(logger.Config{})
	srv.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	srv.router = r
	srv.stopChan = make(chan struct{})

	srv.storage = inmemory.NewInMemory(srv.logger)
	srv.handlers = handlers.NewHandler(srv.storage, srv.logger)

	srv.router.HandleFunc("/quotes", srv.handlers.CreateQuote).Methods("POST")
	srv.router.HandleFunc("/quotes", srv.handlers.GetAllQuotesHandler).Methods("GET")
	srv.router.HandleFunc("/quotes/random", srv.handlers.GetRandomQuoteHandler).Methods("GET")
	srv.router.HandleFunc("/quotes/{id}", srv.handlers.DeleteQuoteHandler).Methods("DELETE")

	return srv
}

func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) StartServer() {
	s.httpServer.Handler = s.router

	go func() {
		s.logger.Info("Starting server on " + s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Ждём сигнал остановки
	s.waitForShutdown()
}

func (s *Server) StopServer() {
	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.logger.Sync()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	close(s.stopChan)
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.StopServer()
}
