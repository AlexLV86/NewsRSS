package main

import (
	"NewsRSS/Comments/pkg/api"
	"NewsRSS/Comments/pkg/postgres"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"
)

// Сервер GoNews.
type server struct {
	db  *postgres.Storage
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server
	path := os.Getenv("dbcomm")
	if path == "" {
		os.Exit(1)
	}
	var err error
	srv.db, err = postgres.New(path)
	if err != nil {
		log.Fatal(err)
	}
	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}...)
	//http.ListenAndServe(":10001", srv.api.Router())
	s := &http.Server{Addr: ":10001", Handler: srv.api.Router()}
	go func() {
		log.Println("[*] HTTP server is started on ", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln("HTTP server error - ", err)
		}
	}()
	// Wait for an interrupt
	sig := <-quit
	log.Println("[*] HTTP server has been stopped. Reason: ", sig)
	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
