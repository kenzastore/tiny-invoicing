# Specification: Fix handlers tests: remove invalid function mocking; use Interface DI

## Overview
This track resolves issues in handler tests caused by brittle and invalid function mocking (e.g., reassigning global database variables). We will transition the `CreateInvoice` handler to use Interface Dependency Injection (DI). This ensures that handlers are decoupled from specific database implementations, making tests reliable, fast, and idiomatic without relying on driver-level mocking.

## Objectives
- Remove direct dependencies on global variables (like `database.DB`) within handlers.
- Define a consumer-owned interface for invoice persistence close to the handlers.
- Implement an adapter in the `database` package that satisfies this interface.
- Refactor handler tests to use a clean mock implementation of the interface.

## Functional Requirements
- **Interface Definition:** Define an `InvoiceStore` interface within the `handlers` package (or `internal/handlers`). 
    - Initially, it only needs to support the `CreateInvoice` operation.
- **Dependency Injection:** Refactor the `CreateInvoice` handler. 
    - Option A: Change `CreateInvoice` to a method on a `Handler` struct that holds the `InvoiceStore`.
    - Option B: Update the handler to accept the interface as a parameter (via a closure or factory).
- **Database Adapter:** Create a small adapter in the `database` package (e.g., `database.Store`) that implements `InvoiceStore` by calling the existing `database.CreateInvoice` logic.
- **Wiring:** Update `main.go` to instantiate the concrete database store and inject it into the handler(s).

## Testing Strategy
- **Mocking:** Create a simple `MockInvoiceStore` struct in the `handlers_test.go` file.
- **Cleanup:** Remove usage of `sqlmock` and reassignment of `database.DB` in the handler tests.
- **Verification:** Ensure tests for `CreateInvoice` verify that the handler correctly processes input and calls the store with the expected data.

## Acceptance Criteria
- `CreateInvoice` handler is refactored to use the injected `InvoiceStore` interface.
- Handler tests pass using the mock interface instead of global variable manipulation.
- `main.go` correctly wires the database adapter to the handler.
- No new "repository" package is introduced; logic remains in `handlers` and `database`.

## Out of Scope
- Refactoring retrieval handlers (Get/List) or other features not mentioned.
- Introducing a complex Repository pattern or new external packages.
