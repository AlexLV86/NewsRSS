package rss

import "testing"

func TestRSS_Get(t *testing.T) {
	var articles *RSS
	//url := "https://habr.com/ru/rss/best/daily/?fl=ru"
	//url := "https://www.finam.ru/analysis/conews/rsspoint/"
	url := "https://realty.interfax.ru/ru/rss/"
	// url := "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml"
	articles, err := Get(url)
	if err != nil {
		t.Fatalf("Ошибка при получении rss: %v", err)
	}
	t.Log(articles.Channel.Item[0].Link)
	t.Log(articles.Channel.Item[1].Link)
	t.Log(articles.Channel.Item[2].Link)
	var got, nowant int64 = 0, 0
	got = int64(len(articles.Channel.Item))
	if got == nowant {
		t.Fatalf("got %d, no want %d", got, nowant)
	}
	pubDate := "Mon, 2 Jan 2006 15:04:05 -0700"
	got, err = DateToUnix(pubDate)
	if err != nil {
		t.Fatalf("Ошибка конвертации даты: %v", err)
	}
	want := int64(1136239445)
	if got != want {
		t.Fatalf("got %d, no want %d", got, want)
	}
}
