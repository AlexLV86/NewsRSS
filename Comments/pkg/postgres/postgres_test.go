package postgres

import (
	"log"
	"os"
	"testing"
)

var s *Storage

func TestMain(m *testing.M) {
	path := os.Getenv("dbcomm")
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
func TestStorage_Comments(t *testing.T) {

	c := Comment{Content: "Прекрасная тестовый комментарий для постгрес!",
		NewsID: 2}
	err := s.AddComment(c)
	if err != nil {
		t.Fatal(err)
	}

	newsid := 1
	comm, err := s.Comments(newsid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(comm)
}
