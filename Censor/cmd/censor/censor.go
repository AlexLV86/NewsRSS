package main

import (
	"censor/pkg/api"
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
	api *api.API
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}...)

	// Создаём объект сервера.
	var srv server

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New()
	// Запускаем веб-сервер сервиса цензурирования на порту 10003 на всех интерфейсах.

	//http.ListenAndServe(":10003", srv.api.Router())
	s := &http.Server{Addr: ":10003", Handler: srv.api.Router()}
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
