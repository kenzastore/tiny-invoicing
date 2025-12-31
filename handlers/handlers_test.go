package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tiny-invoicing/database"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHandlers(t *testing.T) {
	// TODO: Implement actual handler tests with a test server and mocked database
	t.Skip("Skipping handler tests until a test server and mocked database setup is available.")
}

func TestCreateInvoice_InvalidInput(t *testing.T) {
	// Setup a test server
	reqBody := []byte(`{
		"client_id": 0,
		"issue_date": "0001-01-01T00:00:00Z",
		"due_date": "0001-01-01T00:00:00Z",
		"total": 150.75,
		"status": "draft",
		"line_items": []
	}`)
	req, err := http.NewRequest("POST", "/invoices", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateInvoice) // Assuming CreateInvoice is exported

	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "{\"error\":\"Missing required fields\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateInvoice_Success(t *testing.T) {
	// Setup sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Swap global DB
	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	reqBody := []byte(`{
		"client_id": 1,
		"issue_date": "2025-12-31T00:00:00Z",
		"due_date": "2026-01-14T00:00:00Z",
		"status": "draft",
		"line_items": [
			{"description": "Item 1", "quantity": 2, "unit_price": 10.0},
			{"description": "Item 2", "quantity": 1, "unit_price": 5.0}
		]
	}`)
	req, err := http.NewRequest("POST", "/invoices", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO invoices").
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), "draft", 25.0). // 2*10 + 1*5 = 25
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO line_items").
		WithArgs(1, "Item 1", 2, 10.0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO line_items").
		WithArgs(1, "Item 2", 1, 5.0).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateInvoice)

	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateInvoice_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	reqBody := []byte(`{
		"client_id": 1,
		"issue_date": "2025-12-31T00:00:00Z",
		"due_date": "2026-01-14T00:00:00Z",
		"status": "draft",
		"line_items": [{"description": "Item 1", "quantity": 1, "unit_price": 10.0}]
	}`)
	req, err := http.NewRequest("POST", "/invoices", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO invoices").
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), "draft", 10.0).
		WillReturnError(fmt.Errorf("db error"))
	mock.ExpectRollback()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateInvoice)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "{\"error\":\"Failed to create invoice\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateInvoice_InvalidJSON(t *testing.T) {
	reqBody := []byte(`{invalid json}`)
	req, err := http.NewRequest("POST", "/invoices", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateInvoice)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "{\"error\":\"Invalid request payload\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetInvoice_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	issueDate := time.Now().Truncate(time.Second)
	dueDate := issueDate.Add(14 * 24 * time.Hour)

	// Expectations for Invoice
	invoiceRows := sqlmock.NewRows([]string{"id", "client_id", "issue_date", "due_date", "status", "total"}).
		AddRow(1, 1, issueDate, dueDate, "draft", 25.0)

	mock.ExpectQuery("SELECT id, client_id, issue_date, due_date, status, total FROM invoices WHERE id = ?").
		WithArgs(1).
		WillReturnRows(invoiceRows)

	// Expectations for Line Items
	itemRows := sqlmock.NewRows([]string{"id", "invoice_id", "description", "quantity", "unit_price"}).
		AddRow(1, 1, "Item 1", 2, 10.0).
		AddRow(2, 1, "Item 2", 1, 5.0)

	mock.ExpectQuery("SELECT id, invoice_id, description, quantity, unit_price FROM line_items WHERE invoice_id = ?").
		WithArgs(1).
		WillReturnRows(itemRows)

	req, err := http.NewRequest("GET", "/api/invoices/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetInvoice)

	handler.ServeHTTP(rr, req)

	// Assertions
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetInvoice_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	mock.ExpectQuery("SELECT id, client_id, issue_date, due_date, status, total FROM invoices WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/api/invoices/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetInvoice)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetInvoices_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	issueDate := time.Now().Truncate(time.Second)
	dueDate := issueDate.Add(14 * 24 * time.Hour)

	rows := sqlmock.NewRows([]string{"id", "client_id", "issue_date", "due_date", "status", "total"}).
		AddRow(1, 1, issueDate, dueDate, "draft", 25.0).
		AddRow(2, 2, issueDate, dueDate, "paid", 100.0)

	mock.ExpectQuery("SELECT id, client_id, issue_date, due_date, status, total FROM invoices ORDER BY issue_date DESC LIMIT \\? OFFSET \\?").
		WithArgs(20, 0).
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/invoices", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetInvoices)

	handler.ServeHTTP(rr, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetInvoices_Pagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	oldDB := database.DB
	database.DB = db
	defer func() { database.DB = oldDB }()

	rows := sqlmock.NewRows([]string{"id", "client_id", "issue_date", "due_date", "status", "total"}).
		AddRow(1, 1, time.Now(), time.Now(), "draft", 25.0)

	mock.ExpectQuery("SELECT id, client_id, issue_date, due_date, status, total FROM invoices ORDER BY issue_date DESC LIMIT \\? OFFSET \\?").
		WithArgs(10, 5).
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/invoices?limit=10&offset=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetInvoices)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
