# Invoice Management Plan

## Track: Implement basic invoice management functionalities, including creation, viewing, and listing of invoices.

---

## Phase 1: Invoice Data Model and API Endpoints

### Objective
Define the data model for invoices and implement the core API endpoints for creation and retrieval.

### Tasks

- [x] Task: Define Invoice and LineItem structs in a new `models/invoice.go` file. [5b9179e]
- [x] Task: Create a new `migrations/0001_create_invoices_table.up.sql` file with the necessary SQL statements. [28ad979]
- [c] Task: Implement a stub for the `CreateInvoice` handler.
- [c] Task: Implement a stub for the `GetInvoice` handler.
- [c] Task: Implement a stub for the `ListInvoices` handler.
- [ ] Task: Refactor `database/database.go` to use `models.Invoice` and `models.LineItem` structs.
- [ ] Task: Refactor `handlers/handlers.go` to use `models.Invoice` and `models.LineItem` structs.
- [ ] Task: Adjust `migrations/0001_create_invoices_table.up.sql` to match `models.Invoice` and `models.LineItem` struct fields (especially 'status' vs 'paid').
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Invoice Data Model and API Endpoints' (Protocol in workflow.md)

## Phase 2: Invoice Creation

### Objective
Implement the logic and API endpoint for creating new invoices, including input validation and line item management.

### Tasks
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Invoice Creation' (Protocol in workflow.md)

## Phase 3: View and List Invoices

### Objective
Implement the logic and API endpoints for viewing a single invoice and listing all invoices, including basic sorting and pagination.

### Tasks
- [ ] Task: Conductor - User Manual Verification 'Phase 3: View and List Invoices' (Protocol in workflow.md)
