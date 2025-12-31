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

// EnsureDefaultCustomer creates a default customer if none exists.
// Also performs schema patches to ensure decimal columns are large enough.
func EnsureDefaultCustomer() error {
	// 1. Patch Schema: Widen DECIMAL columns to avoid "Out of range" errors
	// Using DECIMAL(20, 2) allows numbers up to 99,999,999,999,999,999.99
	schemaPatches := []string{
		"ALTER TABLE invoices MODIFY total DECIMAL(20, 2)",
		"ALTER TABLE invoice_items MODIFY total DECIMAL(20, 2)",
		"ALTER TABLE invoice_items MODIFY unit_price DECIMAL(20, 2)",
	}

	for _, query := range schemaPatches {
		_, err := DB.Exec(query)
		if err != nil {
			// Ignore error if table doesn't exist yet, usually it's fine
			fmt.Printf("Schema patch warning: %v\n", err)
		}
	}

	// 2. Ensure Default Customer
	// We use standard SQL logic: Try to select, if missing, insert explicitly with ID=1.
	var exists int
	err := DB.QueryRow("SELECT 1 FROM customers WHERE id = 1").Scan(&exists)
	
	if err == sql.ErrNoRows {
		// Force insert ID 1. Using explicit ID overrides auto-increment in MySQL.
		_, err = DB.Exec("INSERT INTO customers (id, name, email, address) VALUES (1, 'Demo Client', 'demo@example.com', '123 Tech Street')")
		if err != nil {
			return fmt.Errorf("failed to create default customer: %v", err)
		}
		fmt.Println("Default customer (ID: 1) created/restored.")
	} else if err != nil {
		return err
	}
	
	return nil
}

// CreateInvoice creates a new invoice and its items in a transaction.
func CreateInvoice(invoice *models.Invoice) (int64, error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	// AUTO-HEAL: Check if customer exists, if not create it to satisfy Foreign Key
	var exists int
	err = tx.QueryRow("SELECT 1 FROM customers WHERE id = ?", invoice.CustomerID).Scan(&exists)
	if err == sql.ErrNoRows {
		// Customer missing! Auto-create it inside the same transaction
		_, err = tx.Exec("INSERT INTO customers (id, name, email, address) VALUES (?, ?, ?, ?)",
			invoice.CustomerID,
			fmt.Sprintf("Auto Client %d", invoice.CustomerID),
			"auto@demo.com",
			"Auto Created Address")
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to auto-create missing customer: %v", err)
		}
	} else if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Map 'status' to 'paid' boolean (draft/sent = false, paid = true)
	isPaid := (invoice.Status == "paid")

	result, err := tx.Exec("INSERT INTO invoices (customer_id, issue_date, due_date, paid, total) VALUES (?, ?, ?, ?, ?)",
		invoice.CustomerID, invoice.IssueDate, invoice.DueDate, isPaid, invoice.Total)
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
		// Calculate item total for DB
		itemTotal := float64(item.Quantity) * item.UnitPrice
		_, err := tx.Exec("INSERT INTO invoice_items (invoice_id, description, quantity, unit_price, total) VALUES (?, ?, ?, ?, ?)",
			invoiceID, item.Description, item.Quantity, item.UnitPrice, itemTotal)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return invoiceID, tx.Commit()
}

// GetInvoices retrieves a paginated list of invoices.
func GetInvoices(limit, offset int) ([]models.Invoice, error) {
	rows, err := DB.Query("SELECT id, customer_id, issue_date, due_date, paid, total FROM invoices ORDER BY issue_date DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		var isPaid bool
		if err := rows.Scan(&invoice.ID, &invoice.CustomerID, &invoice.IssueDate, &invoice.DueDate, &isPaid, &invoice.Total); err != nil {
			return nil, err
		}
		if isPaid {
			invoice.Status = "paid"
		} else {
			invoice.Status = "draft"
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

// GetInvoiceByID retrieves a single invoice by its ID, including its items.
func GetInvoiceByID(id int) (*models.Invoice, error) {
	var invoice models.Invoice
	var isPaid bool
	err := DB.QueryRow("SELECT id, customer_id, issue_date, due_date, paid, total FROM invoices WHERE id = ?", id).Scan(
		&invoice.ID, &invoice.CustomerID, &invoice.IssueDate, &invoice.DueDate, &isPaid, &invoice.Total)
	if err != nil {
		return nil, err
	}
	
	if isPaid {
		invoice.Status = "paid"
	} else {
		invoice.Status = "draft"
	}

	rows, err := DB.Query("SELECT id, invoice_id, description, quantity, unit_price, total FROM invoice_items WHERE invoice_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.LineItem
	for rows.Next() {
		var item models.LineItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Total); err != nil {
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
