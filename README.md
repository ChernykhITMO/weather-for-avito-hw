# weather-service

Небольшой сервис погоды на Go с простой layered/clean-структурой:

- `domain` - сущности и интерфейсы
- `usecase` - бизнес-логика
- `repository` - `in-memory` реализация
- `transport/http` - HTTP-ручки

## Run

```bash
go run ./cmd/weather-service
```

Сервис стартует на `:8080`.

## Endpoints

### POST /weather

Создать или обновить погоду для города.

```bash
curl -X POST http://localhost:8080/weather \
  -H "Content-Type: application/json" \
  -d '{"city":"Moscow","temperature":18.5,"condition":"Cloudy"}'
```

### GET /weather?city=Moscow

Получить погоду по городу.

```bash
curl "http://localhost:8080/weather?city=Moscow"
```

### GET /health

Проверка состояния сервиса.

```bash
curl "http://localhost:8080/health"
```

## Tests

```bash
go test ./...
```
