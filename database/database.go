package database

import (
	"database/sql"
	"fmt"
	"os"
	"tiny-invoicing/models" // Add models import

	_ "github.com/go-sql-driver/mysql"
)

// DB is the database connection.
var DB *sql.DB

// Customer represents a customer.
type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// User represents a user.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
}

// InitDB initializes the database connection.
func InitDB() error {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return fmt.Errorf("DB_DSN environment variable not set")
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	return DB.Ping()
}

// CreateInvoice creates a new invoice and its items in a transaction.
func CreateInvoice(invoice *models.Invoice) (int64, error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	// Use Status from models.Invoice
	result, err := tx.Exec("INSERT INTO invoices (client_id, issue_date, due_date, status, total) VALUES (?, ?, ?, ?, ?)",
		invoice.ClientID, invoice.IssueDate, invoice.DueDate, invoice.Status, invoice.Total)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, item := range invoice.LineItems {
		_, err := tx.Exec("INSERT INTO line_items (invoice_id, description, quantity, unit_price) VALUES (?, ?, ?, ?)",
			invoiceID, item.Description, item.Quantity, item.UnitPrice)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return invoiceID, tx.Commit()
}

// GetInvoices retrieves a paginated list of invoices.
func GetInvoices(limit, offset int) ([]models.Invoice, error) {
	rows, err := DB.Query("SELECT id, client_id, issue_date, due_date, status, total FROM invoices ORDER BY issue_date DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		if err := rows.Scan(&invoice.ID, &invoice.ClientID, &invoice.IssueDate, &invoice.DueDate, &invoice.Status, &invoice.Total); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

// GetInvoiceByID retrieves a single invoice by its ID, including its items.
func GetInvoiceByID(id int) (*models.Invoice, error) {
	var invoice models.Invoice
	err := DB.QueryRow("SELECT id, client_id, issue_date, due_date, status, total FROM invoices WHERE id = ?", id).Scan(
		&invoice.ID, &invoice.ClientID, &invoice.IssueDate, &invoice.DueDate, &invoice.Status, &invoice.Total)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT id, invoice_id, description, quantity, unit_price FROM line_items WHERE invoice_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.LineItem
	for rows.Next() {
		var item models.LineItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	invoice.LineItems = items

	return &invoice, nil
}

// UpdateInvoiceStatusString updates the status of an invoice.
func UpdateInvoiceStatusString(id int, status string) error {
	_, err := DB.Exec("UPDATE invoices SET status = ? WHERE id = ?", status, id)
	return err
}

// CreateUser creates a new user.
func CreateUser(user *User) (int64, error) {
	result, err := DB.Exec("INSERT INTO users (username, password_hash, is_admin) VALUES (?, ?, ?)",
		user.Username, user.PasswordHash, user.IsAdmin)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT id, username, password_hash, is_admin FROM users WHERE username = ?", username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
