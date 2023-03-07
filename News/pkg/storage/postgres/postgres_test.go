package postgres

import (
	"News/pkg/storage"
	"log"
	"os"
	"testing"
)

var s *Storage

func TestMain(m *testing.M) {
	path := os.Getenv("dbpath")
	if path == "" {
		m.Run()
	}
	var err error
	s, err = New(path)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
func TestStorage_Posts(t *testing.T) {
	post := []storage.Post{{Title: "Тестовая статья", Content: "Прекрасная тестовая статья для постгрес!",
		PubTime: 123242, Link: "https://ya3.ru"},
		{Title: "Тестовая статья", Content: "Прекрасная тестовая статья для постгрес!",
			PubTime: 123244524, Link: "https://ya4.ru"}}
	_, err := s.AddPosts(post)
	if err != nil {
		t.Fatal(err)
	}
	data, err := s.Posts(1, "пут")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)

	data, err = s.Posts(1, "")

	if err != nil {
		t.Fatal(err)
	}
	want := 15
	got := len(data.Posts)
	if want != got {
		t.Fatalf("Хотел получить %d получил %d", want, got)
	}

	t.Log(data)
}

func TestStorage_PostsDetail(t *testing.T) {
	newsID := 1293
	p, err := s.PostsDetail(newsID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
	newsID = 1000000
	p, err = s.PostsDetail(newsID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}
