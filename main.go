package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type Invoice struct {
	ID         int           `json:"id"`
	CustomerID int           `json:"customer_id"`
	IssueDate  string        `json:"issue_date"`
	DueDate    string        `json:"due_date"`
	Status     string        `json:"status"`
	Notes      string        `json:"notes"`
	Items      []InvoiceItem `json:"items,omitempty"`
}

var db *sql.DB

func main() {
	dsn := os.Getenv("INVOICE_DB_DSN")
	if dsn == "" {
		dsn = "cloud:cloud.kenzastore.my.id@tcp(127.0.0.1:3306)/cloud?parseTime=true&charset=utf8mb4&loc=Local"
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to invoice_db")

	// API routes
	http.HandleFunc("/api/invoices", invoicesHandler) // GET, POST

	// static frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // different from notes app
	}
	log.Println("Invoice API on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// GET /api/invoices -> list invoices (no items for now)
func invoicesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listInvoices(w, r)
	case http.MethodPost:
		createInvoice(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func listInvoices(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, customer_id, issue_date, due_date, status, notes
		FROM invoices
		ORDER BY id DESC`,
	)
	if err != nil {
		log.Println("listInvoices:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var invs []Invoice
	for rows.Next() {
		var inv Invoice
		var issue, due time.Time
		if err := rows.Scan(&inv.ID, &inv.CustomerID, &issue, &due, &inv.Status, &inv.Notes); err != nil {
			log.Println("scan:", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		inv.IssueDate = issue.Format("2006-01-02")
		inv.DueDate = due.Format("2006-01-02")
		invs = append(invs, inv)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invs)
}

// POST /api/invoices with JSON:
//
//	{
//	  "customer_id": 1,
//	  "issue_date": "2025-01-01",
//	  "due_date": "2025-01-15",
//	  "status": "sent",
//	  "notes": "Optional",
//	  "items": [{ "description": "...", "quantity": 2, "unit_price": 10.5 }]
//	}
func createInvoice(w http.ResponseWriter, r *http.Request) {
	var body Invoice
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if body.CustomerID == 0 || len(body.Items) == 0 {
		http.Error(w, "customer_id and items required", http.StatusBadRequest)
		return
	}

	issue, err := time.Parse("2006-01-02", body.IssueDate)
	if err != nil {
		http.Error(w, "invalid issue_date", http.StatusBadRequest)
		return
	}
	due, err := time.Parse("2006-01-02", body.DueDate)
	if err != nil {
		http.Error(w, "invalid due_date", http.StatusBadRequest)
		return
	}
	if body.Status == "" {
		body.Status = "draft"
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("tx begin:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO invoices (customer_id, issue_date, due_date, status, notes)
		VALUES (?, ?, ?, ?, ?)`,
		body.CustomerID, issue, due, body.Status, body.Notes,
	)
	if err != nil {
		log.Println("insert invoice:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	invoiceID64, _ := res.LastInsertId()
	invoiceID := int(invoiceID64)

	for _, it := range body.Items {
		_, err := tx.Exec(`
			INSERT INTO invoice_items (invoice_id, description, quantity, unit_price)
			VALUES (?, ?, ?, ?)`,
			invoiceID, it.Description, it.Quantity, it.UnitPrice,
		)
		if err != nil {
			log.Println("insert item:", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("commit:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	body.ID = invoiceID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}
