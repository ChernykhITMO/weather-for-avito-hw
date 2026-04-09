# Tasks: Weather Query History

**Input**: Design documents from `/specs/001-weather-history/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Automated tests are REQUIRED for every new behavior. Include explicit
test tasks for each user story and for any shared behavior that changes.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- This project uses `internal/domain`, `internal/repository/memory`,
  `internal/usecase`, `internal/transport/http`, and `cmd/weather-service`.
- Tests live next to the Go files they validate.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Capture current behavior and confirm the exact files that will change

- [X] T001 Inspect the current weather request flow in `internal/transport/http/handler.go`, `internal/usecase/weather_service.go`, and `internal/repository/memory/weather_repository.go`
- [X] T002 Review existing tests and update targets in `internal/transport/http/handler_test.go` and `internal/usecase/weather_service_test.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared building blocks required before either user story can be completed

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Extend the history domain model and repository contract in `internal/domain/weather.go`
- [X] T004 Implement bounded in-memory history storage with newest-first ordering and 10-entry truncation in `internal/repository/memory/weather_repository.go`
- [X] T005 Add repository-level tests for history storage ordering and 10-entry limit in `internal/repository/memory/weather_repository_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Save Successful Weather Lookups (Priority: P1) 🎯 MVP

**Goal**: Record a history entry after each successful weather lookup without changing existing weather response contracts

**Independent Test**: Seed weather data, perform a successful `GET /weather`, then verify exactly one history record exists with city, temperature, condition, and request time; verify a failed lookup does not append history.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T006 [P] [US1] Add usecase test proving a successful weather lookup is saved to history in `internal/usecase/weather_service_test.go`
- [X] T007 [P] [US1] Add usecase test proving a failed weather lookup does not create a history record in `internal/usecase/weather_service_test.go`

### Implementation for User Story 1

- [X] T008 [US1] Add history append and retrieval behavior to `internal/usecase/weather_service.go`
- [X] T009 [US1] Ensure successful `GET /weather` requests trigger history persistence via `internal/transport/http/handler.go` without changing existing payloads

**Checkpoint**: At this point, successful weather lookups should create history entries and failed lookups should not

---

## Phase 4: User Story 2 - View Weather Lookup History (Priority: P2)

**Goal**: Expose `GET /history` returning the latest saved successful lookups as JSON from newest to oldest

**Independent Test**: Create multiple successful history entries, call `GET /history`, and verify HTTP 200, JSON array response, newest-first order, and a maximum of 10 records.

### Tests for User Story 2 ⚠️

- [X] T010 [P] [US2] Add transport test for `GET /history` success and empty-list response in `internal/transport/http/handler_test.go`
- [X] T011 [P] [US2] Add transport test for `GET /history` newest-first ordering in `internal/transport/http/handler_test.go`

### Implementation for User Story 2

- [X] T012 [US2] Add usecase method for reading history from the repository in `internal/usecase/weather_service.go`
- [X] T013 [US2] Add `GET /history` routing and JSON response handling in `internal/transport/http/handler.go`
- [X] T014 [US2] Wire history response fields to the HTTP contract, including `requested_at`, in `internal/transport/http/handler.go`

**Checkpoint**: At this point, `GET /history` should return up to 10 saved entries in reverse chronological order and `[]` when empty

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final verification and documentation alignment across both user stories

- [ ] T015 Update endpoint documentation for `GET /history` and history behavior in `README.md`
- [ ] T016 Run and, if needed, fix targeted automated tests for `internal/repository/memory`, `internal/usecase`, and `internal/transport/http`
- [ ] T017 Run `go test ./...` from the repository root and confirm the quickstart flow in `specs/001-weather-history/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational completion
- **User Story 2 (Phase 4)**: Depends on Foundational completion and reuses history model/repository support from User Story 1
- **Polish (Phase 5)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Starts after foundational tasks and delivers the MVP by recording only successful weather lookups
- **User Story 2 (P2)**: Depends on the history data produced by US1 and exposes it via HTTP

### Within Each User Story

- Tests MUST be written and fail before implementation
- Domain/repository support comes before usecase orchestration
- Usecase changes come before handler changes
- Handler changes stay limited to HTTP mapping and response formatting
- Story complete before moving to the next priority

### Parallel Opportunities

- T006 and T007 can run in parallel after T005
- T010 and T011 can run in parallel after T009
- T015 can run in parallel with T016 after T014

---

## Parallel Example: User Story 1

```bash
Task: "Add usecase test proving a successful weather lookup is saved to history in internal/usecase/weather_service_test.go"
Task: "Add usecase test proving a failed weather lookup does not create a history record in internal/usecase/weather_service_test.go"
```

## Parallel Example: User Story 2

```bash
Task: "Add transport test for GET /history success and empty-list response in internal/transport/http/handler_test.go"
Task: "Add transport test for GET /history newest-first ordering in internal/transport/http/handler_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. Validate that only successful weather lookups create history entries

### Incremental Delivery

1. Build domain and repository support for bounded history
2. Add history recording after successful `GET /weather`
3. Add `GET /history` HTTP read path
4. Validate ordering, empty-state behavior, and test suite

### Parallel Team Strategy

1. One developer handles domain/repository groundwork through T005
2. One developer can prepare US1 tests while another refines history read-path tests after foundation is ready
3. Final integration happens in handler and verification tasks

---

## Notes

- All tasks follow the required checklist format with IDs, optional labels, and file paths
- MVP scope is User Story 1: recording history for successful weather lookups
- The task list intentionally avoids unrelated refactoring and keeps `cmd/weather-service/main.go` untouched unless implementation proves a minimal constructor change is necessary
