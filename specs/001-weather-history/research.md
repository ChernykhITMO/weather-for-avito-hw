# Research: Weather Query History

## Decision 1: Record history in the existing in-memory repository

**Decision**: Extend the current in-memory weather repository to also keep a
bounded history of successful weather lookups.

**Rationale**: The repository already owns in-memory state and locking. Adding
history there keeps persistence concerns out of `usecase`, preserves the current
architecture, and avoids unnecessary constructor churn or new dependencies.

**Alternatives considered**:
- Separate history repository instance: rejected because it adds wiring and new
  abstraction for a tiny in-memory concern.
- Store history in the handler: rejected because it would violate the
  transport/usecase boundary.

## Decision 2: Record history only after `GetByCity` succeeds

**Decision**: `usecase` will append a history record only after a weather lookup
successfully returns a weather value.

**Rationale**: The feature is explicitly about successful weather requests.
Recording the event in `usecase` keeps business rules out of HTTP handlers and
ensures failed lookups are naturally excluded.

**Alternatives considered**:
- Record before repository lookup: rejected because failures would need later
  rollback or filtering.
- Record in repository `GetByCity`: rejected because repository should not infer
  higher-level product behavior from reads alone.

## Decision 3: Keep only the last 10 entries in newest-first order

**Decision**: The repository will insert new history entries at the front and
truncate the collection to at most 10 records.

**Rationale**: The caller requires newest-to-oldest ordering and a maximum of
10 items. Front insertion plus truncation keeps retrieval simple and bounded.

**Alternatives considered**:
- Append then reverse on read: rejected because it adds extra transformation
  logic on every history request.
- Store unlimited entries: rejected because it violates the feature constraint.

## Decision 4: Return `GET /history` as HTTP 200 with a JSON array

**Decision**: The new endpoint returns HTTP 200 with `[]` when empty and an
array of history records when entries exist.

**Rationale**: This matches predictable API behavior for collection reads and
avoids a special-case error contract for "no history yet".

**Alternatives considered**:
- Return 404 when empty: rejected because empty history is a valid state, not a
  missing resource.
- Wrap the array in an envelope object: rejected because the existing service
  returns plain JSON payloads and no envelope is needed for this feature.

## Decision 5: Cover behavior with repository, usecase, and transport tests

**Decision**: Add tests for bounded ordering in the repository, successful-only
recording in the usecase, and HTTP contract coverage for `GET /history`.

**Rationale**: This provides direct verification at the layer where each rule
is enforced while keeping tests focused and aligned with the constitution.

**Alternatives considered**:
- Only end-to-end HTTP tests: rejected because repository ordering and truncation
  rules would be harder to verify precisely.
- Only usecase tests: rejected because the new endpoint contract also needs
  transport verification.
