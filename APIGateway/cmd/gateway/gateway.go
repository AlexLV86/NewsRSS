package main

import (
	"NewsRSS/APIGateway/pkg/apigateway"
	censor "NewsRSS/APIGateway/pkg/client/censorservice"
	comments "NewsRSS/APIGateway/pkg/client/commentservice"
	news "NewsRSS/APIGateway/pkg/client/newsservice"
	"NewsRSS/APIGateway/pkg/config"

	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
)

// Сервер GoNews.
type server struct {
	//db  storage.Interface
	api *apigateway.API
}

func main() {
	c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	signal.Notify(c, []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}...)

	// Создаём объект сервера.
	var srv server
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalln("Config file read error ", err)
	}

	commentService := comments.New(cfg.CommentService.URL)
	newsService := news.New(cfg.NewsService.URL)
	censorService := censor.New(cfg.CensorService.URL)
	// Создаём объект API и регистрируем обработчики.
	srv.api = apigateway.New(commentService, newsService, censorService)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	s := &http.Server{Addr: ":8080", Handler: srv.api.Router()}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[*] HTTP server is started on localhost:8080")
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln("HTTP server error - ", err)
		}
	}()
	// Wait for an interrupt
	sig := <-c
	log.Println("[*] HTTP server has been stopped. Reason: ", sig)
	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	wg.Wait()
}
