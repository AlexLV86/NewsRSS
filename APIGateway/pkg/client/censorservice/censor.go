package censor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	Resource string
}
type CensorService interface {
	Censor(ctx context.Context, comment string) error
}

const (
	ctxRequestID = "request_id"
)

// New - создает новый объект для доступа к сервису комментариев
func New(urlservice string) *client {
	return &client{Resource: urlservice}

}

// Censor - запрос к сервису цензурирования для проверки комментария
func (c *client) Censor(ctx context.Context, comment string) error {
	// запрос в сервис комментариев, получение всех комментариев по id новости
	url := c.Resource
	b, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("ошибка маршализации комментария. error: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("ошибка при составлении запроса к сервису цензурирования. error: %v", err)
	}
	q := req.URL.Query()
	// добавляем request_id к запросу через параметры get
	q.Add(ctxRequestID, ctx.Value(ctxRequestID).(string))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при запросе к сервису цензурирования. error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("код ответа от сервиса цензурирования %d. error: %v", resp.StatusCode, err)
	}
	return nil
}
