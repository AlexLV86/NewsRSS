package news

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type client struct {
	Resource string //
}
type NewsService interface {
	PostsDetail(ctx context.Context, id int) ([]byte, error)            // получение детальной информации о новости
	Posts(ctx context.Context, page int, filter string) ([]byte, error) // получение заданного кол-ва публикаций и фильтр новостей по заголовку
}

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Pagination, пагинация страниц для новостей.
type Pagination struct {
	PageNum    int // номер страницы
	PageItems  int // кол-во элементов на странице
	TotalPages int // кол-во страниц
}

type PostsPagination struct {
	Posts []Post
	Pages Pagination
}

const (
	ctxRequestID = "request_id"
)

// New - создает новый объект для доступа к сервису комментариев
func New(urlservice string) *client {
	return &client{Resource: urlservice}

}

// PostsDetail - запрос к сервису новостей на получение детальной информации о новости
func (c *client) PostsDetail(ctx context.Context, id int) ([]byte, error) {
	url := fmt.Sprintf("%s/%d", c.Resource, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при формировании запроса к сервису новостей. error: %v", err)
	}
	q := req.URL.Query()
	// добавляем request_id к запросу через параметры get
	q.Add(ctxRequestID, ctx.Value(ctxRequestID).(string))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к сервису новостей. error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный код ответа от сервиса новостей %d. error: %v", resp.StatusCode, err)
	}
	post, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа от сервиса новостей. error: %v", err)
	}
	return post, nil
}

// Posts - получение новостей со страницы page по фильтру заголовка filter
func (c *client) Posts(ctx context.Context, page int, filter string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.Resource, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при формировании запроса к сервису новостей. error: %v", err)
	}
	q := req.URL.Query()
	// добавляем request_id к запросу через параметры get
	q.Add(ctxRequestID, ctx.Value(ctxRequestID).(string))
	q.Add("page", strconv.Itoa(page))
	q.Add("s", filter)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к сервису новостей. error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный код ответа от сервиса новостей %d. error: %v", resp.StatusCode, err)
	}
	posts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа от сервиса новостей. error: %v", err)
	}
	return posts, nil
}
