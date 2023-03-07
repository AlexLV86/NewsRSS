package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ItemStruct struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Content struct {
	XMLName xml.Name     `xml:"channel"`
	Title   string       `xml:"title"`
	Item    []ItemStruct `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Content  `xml:"channel"`
}

// получаем публикации из rss канала по url
func Get(url string) (*RSS, error) {
	// подготавливаем запрос на получение rss
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// получаем ответ от rss канала
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// читаем тело ответа в xml формате
	text, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rss RSS
	err = xml.Unmarshal(text, &rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

// конвертер даты из формата
// Mon, 2 Jan 2006 15:04:05 -0700
// в UNIX
func DateToUnix(pubDate string) (int64, error) {
	layoutUTC := "Mon, 2 Jan 2006 15:04:05 -0700"
	layoutGMT := "Mon, 2 Jan 2006 15:04:05 GMT"
	t1p, err := time.Parse(layoutUTC, pubDate)
	if err != nil {
		if t1p, err = time.Parse(layoutGMT, pubDate); err != nil {
			return 0, fmt.Errorf("parse date error %s, ", pubDate)
		}
	}
	return t1p.Unix(), nil
}
