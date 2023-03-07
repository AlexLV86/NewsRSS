package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig()
	if err != nil {
		t.Fatal("Ошибка чтения config ", err)
	}
	t.Log(c)
}
