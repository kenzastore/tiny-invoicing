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
