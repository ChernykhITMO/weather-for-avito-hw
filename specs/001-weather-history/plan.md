# Implementation Plan: Weather Query History

**Branch**: `001-weather-history` | **Date**: 2026-04-09 | **Spec**: [/Users/arseniychernykh/study/go_avito_hw_5/weather-for-avito-hw/specs/001-weather-history/spec.md](/Users/arseniychernykh/study/go_avito_hw_5/weather-for-avito-hw/specs/001-weather-history/spec.md)
**Input**: Feature specification from `/specs/001-weather-history/spec.md`

**Note**: This plan follows the current Go service layout under `internal/`
and the explicit user constraints for in-memory history, minimal changes, and
mandatory automated tests.

## Summary

Add in-memory tracking of successful `GET /weather` lookups and expose the
latest 10 history entries through `GET /history`. The implementation keeps the
existing repository/usecase/transport separation by extending the domain
contracts, adding a bounded in-memory history store in the repository layer,
orchestrating successful lookup recording in `usecase`, and returning JSON from
HTTP handlers without changing the existing `POST /weather` or `GET /weather`
response contracts.

## Technical Context

**Language/Version**: Go 1.24.0  
**Primary Dependencies**: Go standard library (`net/http`, `encoding/json`,
`sync`, `testing`)  
**Storage**: In-memory repository state only; latest 10 successful history
records retained in process memory  
**Testing**: `go test ./...` using standard library `testing`,
`net/http/httptest`, and in-memory repository tests  
**Target Platform**: Linux/macOS server process over HTTP  
**Project Type**: Single-binary web service  
**Performance Goals**: Preserve current in-memory request behavior and return
history in constant bounded size (maximum 10 records)  
**Constraints**: Must use current `repository/usecase/transport` architecture;
save only successful weather lookups; return history newest-first; avoid
unrelated refactoring; avoid `main.go` changes unless wiring forces a minimal
constructor update  
**Scale/Scope**: Single service instance, no persistence across restarts,
history endpoint limited to 10 most recent successful requests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Simplicity and readability: PASS. The design adds one new domain entity, a
  small repository contract extension, and one HTTP endpoint without new
  dependencies.
- Minimal impact changes: PASS. The plan keeps the current flow and only
  touches `domain`, `repository/memory`, `usecase`, `transport/http`, and
  constructor wiring if needed.
- Transport boundaries: PASS. History recording is orchestrated in `usecase`
  after successful lookup; handlers remain limited to request parsing and
  response shaping.
- Mandatory automated tests: PASS. The plan includes usecase, repository, and
  transport-level tests for history behavior.
- Predictable API contracts: PASS. Existing weather endpoints keep their current
  payloads; the new endpoint returns consistent JSON with HTTP 200 and an empty
  list when no history exists.

Post-design re-check: PASS. Phase 1 artifacts preserve the same decisions and
show no unjustified violations.

## Project Structure

### Documentation (this feature)

```text
specs/001-weather-history/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── history-api.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── weather-service/
    └── main.go

internal/
├── domain/
│   └── weather.go
├── repository/
│   └── memory/
│       └── weather_repository.go
├── transport/
│   └── http/
│       ├── handler.go
│       └── handler_test.go
└── usecase/
    ├── weather_service.go
    └── weather_service_test.go
```

**Structure Decision**: Keep the existing single-binary service structure. Add
history-related types to `internal/domain`, bounded in-memory storage to
`internal/repository/memory`, orchestration methods in `internal/usecase`, and
`GET /history` handling plus HTTP tests in `internal/transport/http`. Keep
`cmd/weather-service/main.go` unchanged unless constructor wiring must reflect a
new interface method set.

## Complexity Tracking

No constitution violations are expected. The feature intentionally avoids new
packages, external storage, and unrelated refactoring.
