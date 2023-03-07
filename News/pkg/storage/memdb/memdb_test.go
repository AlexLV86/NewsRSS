package memdb

import (
	"News/pkg/storage"
	"testing"
)

func TestDB_Posts(t *testing.T) {
	db := New()
	p := []storage.Post{{Title: "Test title", Content: "Test content", PubTime: 123214},
		{Title: "Test titrle 2", Content: "Test content 2", PubTime: 123214}}
	got, err := db.AddPosts(p)
	want := 2
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	page, filter := 1, ""
	posts, err := db.Posts(page, filter)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(posts)
	page, filter = 1, "title"
	posts, err = db.Posts(page, filter)
	if err != nil {
		t.Fatal(err)
	}
	got = len(posts.Posts)
	want = 1
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	t.Log(posts)
	newsID := 1
	post, err := db.PostsDetail(newsID)
	if err != nil {
		t.Fatal(err)
	}
	got = post.ID
	want = 1
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	t.Log(post)
	newsID = 10
	post, err = db.PostsDetail(newsID)
	if err != nil {
		t.Fatal(err)
	}
	got = post.ID
	want = 0
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	t.Log(post)
}
