# weather-service

Небольшой сервис погоды на Go с простой layered/clean-структурой:

- `domain` - сущности и интерфейсы
- `usecase` - бизнес-логика
- `repository` - `in-memory` реализация
- `transport/http` - HTTP-ручки

## Development Principles

- Изменения держим простыми и читаемыми.
- Новое поведение добавляем с минимальным влиянием на существующий код.
- `transport/http` отвечает только за HTTP-ввод/вывод, бизнес-логика живет в
  `usecase`.
- Для каждого нового поведения обязательны автоматические тесты.
- API-ответы должны оставаться предсказуемыми и согласованными.

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
