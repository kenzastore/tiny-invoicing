# Invoice Management Plan

## Track: Implement basic invoice management functionalities, including creation, viewing, and listing of invoices.

---

## Phase 1: Invoice Data Model and API Endpoints [checkpoint: 2205983]

### Objective
Define the data model for invoices and implement the core API endpoints for creation and retrieval.

### Tasks

- [x] Task: Define Invoice and LineItem structs in a new `models/invoice.go` file. [5b9179e]
- [x] Task: Create a new `migrations/0001_create_invoices_table.up.sql` file with the necessary SQL statements. [28ad979]
- [c] Task: Implement a stub for the `CreateInvoice` handler.
- [c] Task: Implement a stub for the `GetInvoice` handler.
- [c] Task: Implement a stub for the `ListInvoices` handler.
- [x] Task: Refactor `database/database.go` to use `models.Invoice` and `models.LineItem` structs. [99d4745]
- [x] Task: Refactor `handlers/handlers.go` to use `models.Invoice` and `models.LineItem` structs. [b4ca47c]
- [x] Task: Adjust `migrations/0001_create_invoices_table.up.sql` to match `models.Invoice` and `models.LineItem` struct fields (especially 'status' vs 'paid'). [e3dddc2]
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Invoice Data Model and API Endpoints' (Protocol in workflow.md)

## Phase 2: Invoice Creation [checkpoint: f3c25d9]

### Objective
Implement the logic and API endpoint for creating new invoices, including input validation and line item management.

### Tasks
- [x] Task: Implement input validation for `CreateInvoice` handler. [a2f1cc9]
- [x] Task: Implement line item processing in `CreateInvoice` handler (calculate total for each line item and invoice total). [454236c]
- [x] Task: Write tests for `CreateInvoice` handler. [c0f9c12]
- [x] Task: Conductor - User Manual Verification 'Phase 2: Invoice Creation' (Protocol in workflow.md)

## Phase 3: View and List Invoices

### Objective
Implement the logic and API endpoints for viewing a single invoice and listing all invoices, including basic sorting and pagination.

### Tasks
- [~] Task: Implement full logic for `GetInvoice` handler (retrieve from DB with line items).
- [ ] Task: Implement full logic for `ListInvoices` handler (pagination and sorting).
- [ ] Task: Write tests for `GetInvoice` handler.
- [ ] Task: Write tests for `ListInvoices` handler.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: View and List Invoices' (Protocol in workflow.md)
