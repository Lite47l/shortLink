# URL Shortner

## Запуск проекта

### Запуск с memory 
```bash
$env:STORAGE_TYPE="postgres"
$env:SERVER_PORT="8080"
$env:STORAGE_TYPE="memory"
go run cmd/app/main.go
```

### Запуск с Docker
```bush
docker-compose up --build
```

### Тестирование
```bush
go test ./... -v
```

