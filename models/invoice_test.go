package models

import (
	"testing"
	"time"
)

func TestInvoiceStruct(t *testing.T) {
	issueDate := time.Now()
	dueDate := issueDate.Add(24 * 14 * time.Hour)
	invoice := Invoice{
		ID:         1,
		ClientID:   1,
		IssueDate:  issueDate,
		DueDate:    dueDate,
		Total:      100.00,
		Status:     "draft",
		LineItems: []LineItem{
			{
				ID:          1,
				InvoiceID:   1,
				Description: "Test Item",
				Quantity:    1,
				UnitPrice:   100.00,
			},
		},
	}

	if invoice.ID != 1 {
		t.Errorf("Expected Invoice ID to be 1, but got %d", invoice.ID)
	}
	if invoice.LineItems[0].Description != "Test Item" {
		t.Errorf("Expected Line Item Description to be 'Test Item', but got %s", invoice.LineItems[0].Description)
	}
}

func TestInvoice_CalculateTotal(t *testing.T) {
	invoice := Invoice{
		LineItems: []LineItem{
			{Quantity: 2, UnitPrice: 10.50},
			{Quantity: 1, UnitPrice: 5.00},
			{Quantity: 3, UnitPrice: 20.00},
		},
	}
	invoice.CalculateTotal()
	expected := 86.0
	if invoice.Total != expected {
		t.Errorf("Expected Total to be %f, but got %f", expected, invoice.Total)
	}

	if invoice.LineItems[0].Total != 21.0 {
		t.Errorf("Expected first line item total to be 21.0, but got %f", invoice.LineItems[0].Total)
	}
	if invoice.LineItems[1].Total != 5.0 {
		t.Errorf("Expected second line item total to be 5.0, but got %f", invoice.LineItems[1].Total)
	}
	if invoice.LineItems[2].Total != 60.0 {
		t.Errorf("Expected third line item total to be 60.0, but got %f", invoice.LineItems[2].Total)
	}
}
