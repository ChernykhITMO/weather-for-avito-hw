# Feature Specification: Weather Query History

**Feature Branch**: `001-weather-history`  
**Created**: 2026-04-09  
**Status**: Draft  
**Input**: User description: "Добавить фичу истории запросов погоды.

Пользователь может запросить текущую погоду по городу. После успешного получения погоды сервис должен сохранять запись в историю запросов. Пользователь должен иметь возможность получить историю ранее успешных запросов через отдельный HTTP endpoint.

Критерии приемки:
- после успешного запроса погоды запись сохраняется в историю
- запись истории содержит город, температуру, описание и время запроса
- доступен endpoint GET /history
- GET /history"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Save Successful Weather Lookups (Priority: P1)

As a user, I want successful weather lookups to be recorded automatically so
that the service keeps a history of what weather data was actually returned.

**Why this priority**: Without storing successful lookups, the feature has no
history to show and does not deliver user value.

**Independent Test**: Perform a successful weather lookup for a city and verify
that a new history record exists with the returned weather details and a request
timestamp.

**Acceptance Scenarios**:

1. **Given** weather data exists for a city, **When** the user successfully
   requests the current weather for that city, **Then** the service saves a
   history record for that successful lookup.
2. **Given** a successful weather lookup has just completed, **When** the
   history record is created, **Then** it contains the city, temperature,
   condition description, and the time when the lookup succeeded.
3. **Given** a weather lookup fails, **When** the service returns an error,
   **Then** no new history record is created for that failed attempt.

---

### User Story 2 - View Weather Lookup History (Priority: P2)

As a user, I want to retrieve the history of successful weather lookups so that
I can inspect what weather responses were previously returned by the service.

**Why this priority**: Reading history is only useful after the service stores
history, so it depends on the first user story.

**Independent Test**: Create multiple successful history records and request
the history endpoint to verify that all stored successful lookups are returned
with the expected fields in reverse chronological order.

**Acceptance Scenarios**:

1. **Given** the service has one or more successful weather lookup records,
   **When** the user sends `GET /history`, **Then** the service returns a list
   of previously saved successful lookup records.
2. **Given** multiple successful weather lookups were saved at different times,
   **When** the user sends `GET /history`, **Then** the returned list is ordered
   from newest request to oldest request.
3. **Given** no successful weather lookups were saved yet, **When** the user
   sends `GET /history`, **Then** the service returns a successful empty list
   response rather than an error.

### Edge Cases

- What happens when a weather lookup fails because the city is unknown?
  The service returns the existing error response and MUST NOT append a history
  record.
- What happens when the same city is requested multiple times successfully?
  Each successful lookup is recorded as a separate history entry with its own
  request timestamp.
- What happens when history is requested before any successful lookups exist?
  The service returns an empty list with the normal success contract.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST save a history record after each successful
  current-weather request.
- **FR-002**: The system MUST save only successful weather requests in history.
- **FR-003**: Each history record MUST include the city, temperature, condition
  description, and the time when the successful weather response was produced.
- **FR-004**: Users MUST be able to retrieve saved history records through
  `GET /history`.
- **FR-005**: The history response MUST include all previously saved successful
  weather requests.
- **FR-006**: The history response MUST return records ordered from newest
  request to oldest request.
- **FR-007**: When no history exists, `GET /history` MUST return an empty list
  using the normal success response contract.
- **FR-008**: Adding history tracking MUST NOT change the response contract of
  the existing weather lookup endpoint for successful or failed requests.

### API Consistency Requirements *(mandatory for API changes)*

- `GET /history` returns a success response with HTTP 200 and JSON content.
- A successful response returns a list of history records, where each record
  includes city, temperature, condition, and request time fields.
- If no records exist, the endpoint still returns HTTP 200 with an empty list.
- Errors from existing weather lookup behavior keep their current shape and
  status semantics; history storage is internal and does not introduce new
  response branches for successful lookups.

### Test Coverage Requirements *(mandatory)*

- Add automated coverage proving that a successful weather lookup creates one
  history record with the required fields.
- Add automated coverage proving that a failed weather lookup does not create a
  history record.
- Add automated coverage proving that `GET /history` returns saved records in
  reverse chronological order.
- Add automated coverage proving that `GET /history` returns an empty list when
  no successful lookups exist.
- Use business-level tests for history creation behavior and transport-level
  tests for HTTP mapping and response payloads.

### Key Entities *(include if feature involves data)*

- **Weather Lookup History Record**: A stored representation of one successful
  weather lookup, including city, returned temperature, returned condition
  description, and request timestamp.
- **History Collection**: An ordered set of successful lookup records returned
  to the user through the history endpoint.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of successful weather lookups create exactly one history
  record.
- **SC-002**: 100% of failed weather lookups create zero new history records.
- **SC-003**: Users can retrieve previously saved successful lookups through a
  single history request without needing any follow-up action.
- **SC-004**: 100% of returned history entries include city, temperature,
  condition description, and request time.
- **SC-005**: When multiple history entries exist, 100% of history responses
  present them from newest to oldest.

## Assumptions

- The feature applies to the existing HTTP API and does not introduce
  authentication or per-user history separation.
- History is scoped to successful current-weather lookups served by this
  service instance.
- The existing weather lookup behavior and payload shape remain unchanged.
- Request time means the time at which the service successfully completed the
  weather lookup and prepared the response.
