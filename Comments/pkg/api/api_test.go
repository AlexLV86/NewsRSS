package api

import (
	"NewsRSS/Comments/pkg/postgres"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a *API

func TestMain(m *testing.M) {
	path := os.Getenv("dbcomm")
	if path == "" {
		m.Run()
	}
	db, err := postgres.New(path)
	if err != nil {
		log.Fatal(err)
	}
	a = New(db)
	os.Exit(m.Run())
}

func TestAPI_comments(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/comments/news/2", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	a.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Fatalf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// // Раскодируем JSON в массив заказов.
	var data []postgres.Comment
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(data)
	// // Проверяем, что получили 2 новости
	// const want = 2
	// got := len(data)
	// if got != want {
	// 	t.Fatalf("получено %d , ожидалось %d", got, want)
	// }
}

func TestAPI_addComment(t *testing.T) {
	c := postgres.Comment{ParentID: 10, Content: "Коммент к 10 второй", NewsID: 2}
	var body io.Reader
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal("Ошибка преобразования структуры ", err)
	}
	body = bytes.NewBuffer(b)
	req := httptest.NewRequest(http.MethodPost, "/comments", body)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	a.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Fatalf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
}
