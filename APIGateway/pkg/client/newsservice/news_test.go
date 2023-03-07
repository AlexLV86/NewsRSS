package news

import (
	"context"
	"encoding/json"
	"testing"
)

func Test_client_PostsDetail(t *testing.T) {
	news := New("http://127.0.0.1:10002/news")
	ctx := context.Background()
	// устанавливаем request_id в контекст
	ctx = context.WithValue(ctx, ctxRequestID, "test")
	newsID := 1089
	data, err := news.PostsDetail(ctx, newsID)
	if err != nil {
		t.Fatal(err)
	}
	var p Post
	err = json.Unmarshal(data, &p)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(p)
}

func Test_client_Posts(t *testing.T) {
	news := New("http://127.0.0.1:10002/news")
	ctx := context.Background()
	// устанавливаем request_id в контекст
	ctx = context.WithValue(ctx, ctxRequestID, "test")
	page := 1
	filter := ""
	data, err := news.Posts(ctx, page, filter)
	if err != nil {
		t.Fatal(err)
	}
	var p PostsPagination
	err = json.Unmarshal(data, &p)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(p)
}
