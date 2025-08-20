# dining_room

### Сборка
```bash
go build -o dish cmd/main.go
```

### Запуск
```bash
# Запуск с настройками по умолчанию
./dish

# Или запуск из исходного кода
go run cmd/main.go

# Запуск с кастомными настройками
ADDR=:8080 ./dish
```


## Для теста
```bash
curl -X POST http://localhost:8080/dish \
  -H "Content-Type: application/json" \
  -d '{
    "name": "пюре",
    "price": "100",
    "description": "картоха"
  }'


curl -X POST http://localhost:8080/dish \
  -H "Content-Type: application/json" \
  -d '{
    "name": "котлета",
    "price": "150",
    "description": "говядина"
  }'


  curl -X GET http://localhost:8080/dish 
```