package comments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	Resource string //
}
type CommentService interface {
	AddComment(ctx context.Context, comment Comment) error
	Comments(ctx context.Context, newsID int) ([]byte, error)
}

// Comment - комментарий.
type Comment struct {
	ID        int
	ParentID  int
	NewsID    int
	Content   string
	CreatedAt int64
}

const (
	ctxRequestID = "request_id"
)

// New - создает новый объект для доступа к сервису комментариев
func New(urlservice string) *client {
	return &client{Resource: urlservice}

}

// AddComment - запрос к сервису комментариев на добавление комментария
func (c *client) AddComment(ctx context.Context, comment Comment) error {
	// запрос в сервис комментариев, получение всех комментариев по id новости
	url := c.Resource
	b, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("ошибка маршализации комментария. error: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("ошибка при составлении запроса к сервису комментариев. error: %v", err)
	}
	q := req.URL.Query()
	// добавляем request_id к запросу через параметры get
	q.Add(ctxRequestID, ctx.Value(ctxRequestID).(string))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при запросе к сервису комментариев. error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("код ответа от сервиса комментариев %d. error: %v", resp.StatusCode, err)
	}
	return nil
}

// Comments - получение всех комментариев по id новости
func (c *client) Comments(ctx context.Context, newsID int) ([]byte, error) {
	// запрос в сервис комментариев, получение всех комментариев по id новости
	url := fmt.Sprintf("%s/news/%d", c.Resource, newsID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при формировании запроса к сервису комментариев. error: %v", err)
	}
	q := req.URL.Query()
	// добавляем request_id к запросу через параметры get
	q.Add(ctxRequestID, ctx.Value(ctxRequestID).(string))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к сервису комментариев. error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный код ответа от сервиса комментариев %d. error: %v", resp.StatusCode, err)
	}
	comments, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа от сервиса комментариев. error: %v", err)
	}
	return comments, nil
}
