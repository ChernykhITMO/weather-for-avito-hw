# Data Model: Weather Query History

## Weather

Existing weather value returned by the service.

**Fields**
- `city`: non-empty city name as returned by current storage
- `temperature`: numeric temperature value
- `condition`: non-empty condition description

**Validation**
- `city` must remain non-empty after trimming whitespace
- `condition` is trimmed before persistence by the service

## HistoryRecord

Represents one successful weather lookup captured for later retrieval.

**Fields**
- `city`: city from the successful weather response
- `temperature`: temperature from the successful weather response
- `condition`: condition from the successful weather response
- `requested_at`: timestamp of when the lookup succeeded

**Validation**
- All fields mirror a successful weather response
- `requested_at` must be set when the record is created

## HistoryCollection

Bounded ordered list of `HistoryRecord` values exposed by `GET /history`.

**Rules**
- Ordered newest to oldest
- Maximum size is 10
- Contains only records produced by successful weather lookups
- Empty collection is valid

## Relationships

- One successful `GetByCity` operation produces zero or one `HistoryRecord`
  appended to the collection.
- `HistoryRecord` copies values from `Weather`; it does not replace the primary
  weather storage.

## State Transitions

1. `GET /weather?city=...` succeeds.
2. `usecase` constructs a `HistoryRecord` from the returned `Weather` plus the
   current time.
3. Repository inserts the record at the head of the history collection.
4. Repository truncates the collection if it grows beyond 10 records.
5. `GET /history` returns the current bounded ordered collection.
