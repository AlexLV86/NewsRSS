package comments

import (
	"context"
	"encoding/json"
	"testing"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func Test_client_AddComment(t *testing.T) {
	comment := New("http://127.0.0.1:10001/comments")
	ctx := context.Background()
	// устанавливаем request_id в контекст
	ctx = context.WithValue(ctx, ctxRequestID, "test")
	c := Comment{ParentID: 0, Content: "Тестовый комментарий из апигетвей", NewsID: 2}
	err := comment.AddComment(ctx, c)
	if err != nil {
		t.Fatalf("\t%s\t%v", failed, err)
	}
	t.Logf("\t%s\t", success)
}

func Test_client_Comments(t *testing.T) {
	comment := New("http://127.0.0.1:10001/comments")
	ctx := context.Background()
	// устанавливаем request_id в контекст
	ctx = context.WithValue(ctx, ctxRequestID, "test")
	newsID := 1
	data, err := comment.Comments(ctx, newsID)
	if err != nil {
		t.Fatalf("\t%s\t%v", failed, err)
	}
	var c []Comment
	err = json.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	t.Log(c)
}
