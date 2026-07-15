package service

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	shortURL, err := Generate()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if len(shortURL) != 10 {
		t.Errorf("ожидалась длина 10, получена %d", len(shortURL))
	}

	shortURL2, err := Generate()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if shortURL == shortURL2 {
		t.Error("две последовательные генерации вернули одинаковый результат")
	}

	charset := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890_"
	for _, ch := range shortURL {
		found := false
		for _, validCh := range charset {
			if ch == validCh {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("недопустимый символ '%c' в сгенерированной строке", ch)
		}
	}
}

func TestGenerate_Multiple(t *testing.T) {
	generated := make(map[string]bool)

	for i := 0; i < 100; i++ {
		shortURL, err := Generate()
		if err != nil {
			t.Fatalf("неожиданная ошибка на итерации %d: %v", i, err)
		}

		if generated[shortURL] {
			t.Errorf("коллизия на итерации %d: %s", i, shortURL)
		}
		generated[shortURL] = true
	}
}
