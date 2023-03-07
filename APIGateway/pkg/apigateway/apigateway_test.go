package apigateway

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	censor "NewsRSS/APIGateway/pkg/client/censorservice"
	comments "NewsRSS/APIGateway/pkg/client/commentservice"
	news "NewsRSS/APIGateway/pkg/client/newsservice"
)

func TestAPI_newsPage(t *testing.T) {
	commentService := comments.New("http://127.0.0.1:10001/comments")
	newsService := news.New("http://127.0.0.1:10002/news")
	censorService := censor.New("http://127.0.0.1:10003/censor")
	api := New(commentService, newsService, censorService)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/latest?page=1", nil)
	//req.Context()
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
	// Раскодируем JSON.
	var data int
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(data)
}

func TestAPI_newsDetail(t *testing.T) {
	commentService := comments.New("http://127.0.0.1:10001/comments")
	newsService := news.New("http://127.0.0.1:10002/news")
	censorService := censor.New("http://127.0.0.1:10003/censor")
	api := New(commentService, newsService, censorService)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Fatalf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	// Раскодируем JSON.
	var com []comments.Comment
	err = json.Unmarshal(body, &com)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(com)
}

func TestAPI_newsFilter(t *testing.T) {
	commentService := comments.New("http://127.0.0.1:10001/comments")
	newsService := news.New("http://127.0.0.1:10002/news")
	censorService := censor.New("http://127.0.0.1:10003/censor")
	api := New(commentService, newsService, censorService)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/filter?direction=desc&count=11&offset=0", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
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
	// Раскодируем JSON.
	var data interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(data)
}
