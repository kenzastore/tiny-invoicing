package models

import "time"

// Invoice represents an invoice in the system.
type Invoice struct {
	ID        int        `json:"id"`
	ClientID  int        `json:"client_id"`
	IssueDate time.Time  `json:"issue_date"`
	DueDate   time.Time  `json:"due_date"`
	Total     float64    `json:"total"`
	Status    string     `json:"status"` // e.g., "draft", "sent", "paid", "void"
	LineItems []LineItem `json:"line_items"`
}

// LineItem represents a single line item on an invoice.
type LineItem struct {
	ID          int     `json:"id"`
	InvoiceID   int     `json:"invoice_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}
