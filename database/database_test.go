package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetInvoiceByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	DB = db

	issueDate := time.Now()
	dueDate := issueDate.Add(14 * 24 * time.Hour)

	// Expectations for Invoice
	invoiceRows := sqlmock.NewRows([]string{"id", "client_id", "issue_date", "due_date", "status", "total"}).
		AddRow(1, 1, issueDate, dueDate, "draft", 0.0) // total 0 in DB, should be recalculated

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

	invoice, err := GetInvoiceByID(1)
	if err != nil {
		t.Errorf("GetInvoiceByID returned error: %s", err)
	}

	if invoice.Total != 25.0 {
		t.Errorf("Expected recalculated total to be 25.0, but got %f", invoice.Total)
	}

	if len(invoice.LineItems) != 2 {
		t.Errorf("Expected 2 line items, but got %d", len(invoice.LineItems))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
