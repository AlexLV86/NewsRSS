package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CommentService struct {
		URL string `json:"url"`
	} `json:"comment_service"`
	NewsService struct {
		URL string `json:"url"`
	} `json:"news_service"`
	CensorService struct {
		URL string `json:"url"`
	} `json:"censor_service"`
}

// читаем файл с конфигурацией rss каналов
func ReadConfig() (Config, error) {
	// достать данные из файла для
	c := Config{}
	data, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
