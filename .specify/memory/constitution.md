<!--
Sync Impact Report
- Version change: template -> 1.0.0
- Modified principles:
  - Principle 1 -> I. Simplicity and Readability First
  - Principle 2 -> II. Minimal Impact Changes
  - Principle 3 -> III. HTTP Transport Boundaries
  - Principle 4 -> IV. Mandatory Automated Tests
  - Principle 5 -> V. Predictable API Contracts
- Added sections:
  - Implementation Boundaries
  - Delivery Workflow
- Removed sections:
  - None
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md
  - ✅ .specify/templates/spec-template.md
  - ✅ .specify/templates/tasks-template.md
  - ⚠ pending .specify/templates/commands/*.md (directory absent in repository)
  - ✅ README.md
- Follow-up TODOs:
  - None
-->
# Weather Service Constitution

## Core Principles

### I. Simplicity and Readability First
Every change MUST prefer the simplest design that satisfies the requirement and
keeps the code easy to read. New abstractions, dependencies, and patterns MUST
be introduced only when the current code cannot express the behavior cleanly
without them. Rationale: this service is intentionally small, so readability is
the primary defense against defects and maintenance drift.

### II. Minimal Impact Changes
New behavior MUST change the smallest practical surface area of the existing
codebase. Wide refactors, speculative cleanup, and unrelated restructuring MUST
NOT be bundled with feature delivery unless they are required to keep the
system correct. Rationale: small, contained changes reduce regression risk and
keep review straightforward.

### III. HTTP Transport Boundaries
HTTP handlers MUST be responsible only for receiving HTTP input, validating and
mapping transport data, and returning HTTP responses. Business rules,
orchestration, and domain decisions MUST live in `usecase`. Handlers MUST NOT
contain business logic beyond transport-specific concerns. Rationale: clear
boundaries keep HTTP concerns replaceable and business behavior testable.

### IV. Mandatory Automated Tests
Every new behavior MUST be covered by automated tests that exercise the added or
changed behavior at the appropriate level. A change is not complete until the
relevant tests exist and pass. Tests SHOULD target `usecase` behavior first and
add transport-level coverage when HTTP mapping, validation, or response shaping
changes. Rationale: mandatory tests are the quality gate for safe iteration on
this service.

### V. Predictable API Contracts
API responses MUST be consistent in shape, status semantics, and error handling
for comparable outcomes. New endpoints or response branches MUST follow
established naming and serialization patterns unless a documented incompatibility
forces a change. Rationale: predictable contracts simplify clients, testing, and
future maintenance.

## Implementation Boundaries

- Avoid unnecessary dependencies. New libraries MUST be justified by clear value
  that exceeds the maintenance cost of keeping the dependency.
- Avoid broad refactoring during feature work. Required structural changes MUST
  be minimal and directly tied to the delivered behavior.
- Architecture changes MUST preserve the `transport/http` -> `usecase` ->
  `domain`/`repository` direction of responsibility.
- When a simpler in-repository solution exists, it SHOULD be preferred over
  introducing infrastructure or framework complexity.

## Delivery Workflow

- Commits MUST be small, logically coherent, and scoped to a single intent.
- Plans and tasks MUST make testing explicit for each new behavior rather than
  treating tests as optional polish.
- Reviews MUST check constitution compliance, especially handler/usecase
  separation, regression risk from change scope, and API contract consistency.
- If a task requires violating a principle, the plan MUST document the reason,
  the rejected simpler alternative, and the narrowest acceptable exception.

## Governance

This constitution overrides conflicting local practices for this repository. Any
change to code, plans, tasks, or specifications MUST be reviewed against these
principles. Amendments require: (1) the proposed constitutional change to be
documented, (2) dependent templates and guidance files to be updated in the
same change, and (3) the version to be bumped according to semantic intent.

Versioning policy:
- MAJOR: remove or materially redefine a principle or governance rule.
- MINOR: add a principle or materially expand required process or constraints.
- PATCH: clarify wording without changing expected behavior.

Compliance review expectations:
- Every implementation plan MUST list explicit constitution gates.
- Every specification MUST describe response consistency expectations and test
  coverage for new behavior.
- Every task list MUST include automated test work for changed behavior.

**Version**: 1.0.0 | **Ratified**: 2026-04-09 | **Last Amended**: 2026-04-09
