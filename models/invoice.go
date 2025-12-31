package models

import "time"

// Invoice represents an invoice in the system.
type Invoice struct {
	ID         int        `json:"id"`
	CustomerID int        `json:"customer_id"`
	IssueDate  time.Time  `json:"issue_date"`
	DueDate    time.Time  `json:"due_date"`
	Total      float64    `json:"total"`
	Status     string     `json:"status"` // Kita tetap simpan ini di struct untuk UI, tapi di DB akan dipetakan
	LineItems  []LineItem `json:"line_items"`
}

// LineItem represents a single line item on an invoice.
type LineItem struct {
	ID          int     `json:"id"`
	InvoiceID   int     `json:"invoice_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Total       float64 `json:"total"`
}

// CalculateTotal calculates the total amount of the invoice based on its line items.
func (i *Invoice) CalculateTotal() {
	var grandTotal float64
	for j := range i.LineItems {
		i.LineItems[j].Total = float64(i.LineItems[j].Quantity) * i.LineItems[j].UnitPrice
		grandTotal += i.LineItems[j].Total
	}
	i.Total = grandTotal
}
