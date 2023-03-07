package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"News/pkg/api"
	"News/pkg/rss"
	"News/pkg/storage"
	"News/pkg/storage/postgres"
)

type Config struct {
	ListRss   []string `json:"rss"`
	ReqPeriod int      `json:"request_period"`
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	//db := memdb.New()
	// Реляционная БД PostgreSQL.
	path := os.Getenv("dbpath")
	if path == "" {
		os.Exit(1)
	}
	db, err := postgres.New(path)
	if err != nil {
		log.Fatal(err)
	}

	chErr := make(chan error)
	chPosts := make(chan []storage.Post)
	// запускаем обработчик ошибок
	go handlerError(chErr)
	// запускаем обработчик для добавления новостей в БД
	go handlerDBPosts(chErr, chPosts, db)
	// читаем новости и засыпаем на время из config.json
	go func() {
		for {
			for _, url := range config.ListRss {
				go readRSS(chErr, chPosts, url)
			}
			time.Sleep(time.Minute * time.Duration(config.ReqPeriod))
		}
	}()

	// Создание объекта API, использующего БД в памяти.
	api := api.New(db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}...)
	// http.ListenAndServe(":10002", api.Router())
	s := &http.Server{Addr: ":10002", Handler: api.Router()}
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

// читаем RSS
func readRSS(chErr chan<- error, chPosts chan<- []storage.Post, url string) {
	var articles *rss.RSS
	articles, err := rss.Get(url)
	if err != nil {
		// пишем ошибку в канал
		chErr <- err
		return
	}
	posts := make([]storage.Post, 0, len(articles.Channel.Item))
	var p storage.Post
	// добавляем каждую новость в storage.Post
	for _, v := range articles.Channel.Item {
		p.PubTime, err = rss.DateToUnix(v.PubDate)
		if err != nil {
			// передать в канал ошибку и continue
			chErr <- fmt.Errorf("date convert error %s, wrap error: %w", v.PubDate, err)
			continue
		}
		p.Content = v.Description
		p.Link = v.Link
		p.Title = v.Title
		posts = append(posts, p)
	}
	// передаем все новости в канал
	chPosts <- posts
}

// Обработчик новостей. Добавляем новости в БД
func handlerDBPosts(chErr chan<- error, chPosts <-chan []storage.Post, db storage.Interface) {
	for posts := range chPosts {
		_, err := db.AddPosts(posts)
		if err != nil {
			// передать в канал ошибку и continue
			chErr <- err
			continue
		}
	}
}

// Обработчик ошибок. Записываем все ошибки в файл
func handlerError(chErr <-chan error) {
	for e := range chErr {
		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE, os.ModePerm) //os.O_WRONLY|os.O_CREATE
		if err != nil {
			panic(err)
		}
		str := fmt.Sprintf("%v\n", e)
		f.WriteString(str)
	}
}

// читаем файл с конфигурацией rss каналов
func readConfig() (Config, error) {
	// достать данные из файла для
	c := Config{}
	data, err := os.ReadFile("cmd/gonews/config.json")
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Println(err)
		return Config{}, err
	}
	return c, nil
}
