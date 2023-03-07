package api

import (
	"News/pkg/storage"
	"News/pkg/storage/memdb"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testData(db storage.Interface) {
	p := []storage.Post{{Title: "Test title", Content: "Test content", PubTime: 123214},
		{Title: "Test title 2", Content: "Test content 2", PubTime: 123214}}
	db.AddPosts(p)
}
func TestAPI_newsDetail(t *testing.T) {
	// Создаём чистый объект API для теста.
	db := memdb.New()
	testData(db) // заполняем БД
	api := New(db)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
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
	var data storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что получили 2 новости
	const want = 1
	got := data.ID
	if got != want {
		t.Fatalf("получено %d , ожидалось %d", got, want)
	}
}

func TestAPI_newsPage(t *testing.T) {
	// Создаём чистый объект API для теста.
	db := memdb.New()
	testData(db) // заполняем БД
	api := New(db)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news?page=1", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор.
	api.r.ServeHTTP(rr, req)
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
	var data *storage.PostsPagination
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что получили 2 новости
	const want = 2
	got := len(data.Posts)
	if got != want {
		t.Fatalf("получено %d , ожидалось %d", got, want)
	}
}
