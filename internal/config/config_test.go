package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("STORAGE_TYPE")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("BASE_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("ожидался порт 8080, получен %s", cfg.ServerPort)
	}

	if cfg.StorageType != "memory" {
		t.Errorf("ожидался тип memory, получен %s", cfg.StorageType)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("STORAGE_TYPE", "postgres")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("STORAGE_TYPE")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if cfg.ServerPort != "9090" {
		t.Errorf("ожидался порт 9090, получен %s", cfg.ServerPort)
	}

	if cfg.StorageType != "postgres" {
		t.Errorf("ожидался тип postgres, получен %s", cfg.StorageType)
	}
}

func TestLoad_InvalidStorageType(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "invalid")
	defer os.Unsetenv("STORAGE_TYPE")

	_, err := Load()
	if err == nil {
		t.Error("ожидалась ошибка для невалидного типа хранилища")
	}
}
