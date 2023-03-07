package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_censor(t *testing.T) {
	comment := "Комментарий для проверки валидный"
	b, err := json.Marshal(comment)
	if err != nil {
		t.Fatal("Ошибка преобразования структуры ", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/censor", bytes.NewBuffer(b))
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	a := New()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	a.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Fatalf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	comment = "Комментарий для проверки йцукен невалидный"
	b, err = json.Marshal(comment)
	if err != nil {
		t.Fatal("Ошибка преобразования структуры ", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/censor", bytes.NewBuffer(b))
	// Создаём объект для записи ответа обработчика.
	rr = httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	a.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusBadRequest) {
		t.Fatalf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}
}
