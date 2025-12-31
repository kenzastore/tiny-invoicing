# Implementation Plan: Fix handlers tests: remove invalid function mocking; use Interface DI

## Phase 1: Interface and Adapter Implementation [checkpoint: 56ca87c]

### Objective
Establish the contract for dependency injection and provide a concrete implementation that uses the existing database logic in a separate adapter file.

### Tasks
- [x] Task: Define the `InvoiceStore` interface in `handlers/handlers.go`. [d031f63]
- [x] Task: Create `database/store.go` and implement a `Store` struct that satisfies the `InvoiceStore` interface by calling existing logic. [5eab6bb]
- [x] Task: Verify that the project compiles with `go build ./...`. [5eab6bb]
- [x] Task: Conductor - User Manual Verification 'Phase 1: Interface and Adapter Implementation' (Protocol in workflow.md)

## Phase 2: Handler Refactoring for Dependency Injection [checkpoint: a2a8c51]

### Objective
Update the `CreateInvoice` handler to use the `InvoiceStore` interface instead of direct global database access.

### Tasks
- [x] Task: Refactor `handlers.CreateInvoice` to be a method of a new `InvoiceHandler` struct. [333a590]
- [x] Task: Update `main.go` to instantiate the `database.Store` and inject it into the `InvoiceHandler`. [1e15ec9]
- [x] Task: Verify wiring and compilation with `go build ./...`. [1e15ec9]
- [ ] Task: (Optional) Perform a quick manual smoke test if the environment is ready.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Handler Refactoring for Dependency Injection' (Protocol in workflow.md)

## Phase 3: Test Suite Refactoring [checkpoint: 9d69325]

### Objective
Replace brittle tests that use global state manipulation with clean unit tests using the `InvoiceStore` mock, focusing on critical path coverage.

### Tasks
- [x] Task: Create a `MockInvoiceStore` in `handlers/handlers_test.go`. [21e837a]
- [x] Task: **Red Phase:** Update `TestCreateInvoice_Success` to use the mock and verify that it fails as expected. [21e837a]
- [x] Task: **Green Phase:** Update the test setup to properly inject the mock and make the test pass. [21e837a]
- [x] Task: Refactor other `CreateInvoice` tests (invalid input, database errors) to use the mock and remove all `sqlmock` dependencies. [21e837a]
- [x] Task: Verify that `go test ./handlers/...` passes and critical paths (success, invalid input, store error) are covered. [21e837a]
- [x] Task: Conductor - User Manual Verification 'Phase 3: Test Suite Refactoring' (Protocol in workflow.md)
