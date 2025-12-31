package database

import "tiny-invoicing/models"

// Store is a database adapter that implements handler interfaces.
type Store struct{}

// CreateInvoice calls the existing package-level CreateInvoice function.
func (s *Store) CreateInvoice(invoice *models.Invoice) (int64, error) {
	return CreateInvoice(invoice)
}
