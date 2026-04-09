# Contract: History API

## GET /history

Returns up to 10 most recent successful weather lookup records in descending
request time order.

### Request

- Method: `GET`
- Path: `/history`
- Query parameters: none
- Request body: none

### Success Response

- Status: `200 OK`
- Content-Type: `application/json`

```json
[
  {
    "city": "Moscow",
    "temperature": 18.5,
    "condition": "Cloudy",
    "requested_at": "2026-04-09T15:04:05Z"
  },
  {
    "city": "Kazan",
    "temperature": 21.3,
    "condition": "Sunny",
    "requested_at": "2026-04-09T14:58:10Z"
  }
]
```

### Empty Success Response

- Status: `200 OK`
- Content-Type: `application/json`

```json
[]
```

### Behavioral Rules

- Only successful `GET /weather` requests create history records.
- Failed lookups do not change history.
- The newest record appears first.
- The response contains at most 10 records.
- Existing `/weather` and `/health` contracts remain unchanged.
