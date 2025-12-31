package main

import (
	"log"
	"net/http"
	"time"

	"tiny-invoicing/auth"
	"tiny-invoicing/database"
	"tiny-invoicing/handlers"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	// Set up router
	mux := http.NewServeMux()

	invoiceHandler := &handlers.InvoiceHandler{
		Store: &database.Store{},
	}

	// Public route to create an admin user (for demo purposes)
	mux.HandleFunc("/api/admin/create-user", handlers.CreateAdminUser)

	// API routes
	mux.HandleFunc("/api/invoices", auth.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			invoiceHandler.GetInvoices(w, r)
		case http.MethodPost:
			invoiceHandler.CreateInvoice(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/invoices/", auth.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			invoiceHandler.GetInvoice(w, r)
		case http.MethodPut:
			invoiceHandler.UpdateInvoice(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Static file server
	mux.Handle("/", http.FileServer(http.Dir("static")))

	// Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}