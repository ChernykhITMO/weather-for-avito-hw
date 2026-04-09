# Quickstart: Weather Query History

## Prerequisites

- Start the service:

```bash
go run ./cmd/weather-service
```

## Seed Weather Data

Create weather entries that can later be retrieved successfully:

```bash
curl -X POST http://localhost:8080/weather \
  -H "Content-Type: application/json" \
  -d '{"city":"Moscow","temperature":18.5,"condition":"Cloudy"}'
```

```bash
curl -X POST http://localhost:8080/weather \
  -H "Content-Type: application/json" \
  -d '{"city":"Kazan","temperature":21.3,"condition":"Sunny"}'
```

## Create History Entries

Perform successful weather lookups:

```bash
curl "http://localhost:8080/weather?city=Moscow"
```

```bash
curl "http://localhost:8080/weather?city=Kazan"
```

These successful lookups should create history entries.

## Read History

```bash
curl "http://localhost:8080/history"
```

Expected behavior:
- Returns HTTP 200 with JSON.
- Entries are ordered newest first.
- Returns at most 10 entries.
- Returns `[]` if no successful lookups have occurred yet.

## Verify Failure Does Not Record History

Trigger a failed lookup:

```bash
curl "http://localhost:8080/weather?city=Unknown"
```

Read history again:

```bash
curl "http://localhost:8080/history"
```

Expected behavior:
- Failed lookup returns the existing error contract.
- History remains unchanged after the failed lookup.
