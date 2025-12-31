package database

import (
	"database/sql"
	"fmt"
	"os"

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

// Invoice represents an invoice.
type Invoice struct {
	ID         int           `json:"id"`
	CustomerID int           `json:"customer_id"`
	IssueDate  string        `json:"issue_date"`
	DueDate    string        `json:"due_date"`
	Paid       bool          `json:"paid"`
	Total      float64       `json:"total"`
	Items      []InvoiceItem `json:"items,omitempty"`
}

// InvoiceItem represents an item on an invoice.
type InvoiceItem struct {
	ID          int     `json:"id"`
	InvoiceID   int     `json:"invoice_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Total       float64 `json:"total"`
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
func CreateInvoice(invoice *Invoice) (int64, error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec("INSERT INTO invoices (customer_id, issue_date, due_date, paid, total) VALUES (?, ?, ?, ?, ?)",
		invoice.CustomerID, invoice.IssueDate, invoice.DueDate, invoice.Paid, invoice.Total)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, item := range invoice.Items {
		_, err := tx.Exec("INSERT INTO invoice_items (invoice_id, description, quantity, unit_price, total) VALUES (?, ?, ?, ?, ?)",
			invoiceID, item.Description, item.Quantity, item.UnitPrice, item.Total)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return invoiceID, tx.Commit()
}

// GetInvoices retrieves a paginated list of invoices.
func GetInvoices(limit, offset int) ([]Invoice, error) {
	rows, err := DB.Query("SELECT id, customer_id, issue_date, due_date, paid, total FROM invoices ORDER BY issue_date DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.CustomerID, &invoice.IssueDate, &invoice.DueDate, &invoice.Paid, &invoice.Total); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

// GetInvoiceByID retrieves a single invoice by its ID, including its items.
func GetInvoiceByID(id int) (*Invoice, error) {
	var invoice Invoice
	err := DB.QueryRow("SELECT id, customer_id, issue_date, due_date, paid, total FROM invoices WHERE id = ?", id).Scan(
		&invoice.ID, &invoice.CustomerID, &invoice.IssueDate, &invoice.DueDate, &invoice.Paid, &invoice.Total)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT id, invoice_id, description, quantity, unit_price, total FROM invoice_items WHERE invoice_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var item InvoiceItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Total); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	invoice.Items = items

	return &invoice, nil
}

// UpdateInvoiceStatus updates the paid status of an invoice.
func UpdateInvoiceStatus(id int, paid bool) error {
	_, err := DB.Exec("UPDATE invoices SET paid = ? WHERE id = ?", paid, id)
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
