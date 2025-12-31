package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"tiny-invoicing/auth"
	"tiny-invoicing/database"
	"tiny-invoicing/models" // Add models import
	"tiny-invoicing/response"
)

// CreateInvoice creates a new invoice.
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Basic validation for models.Invoice
	if invoice.ClientID == 0 || invoice.IssueDate.IsZero() || invoice.DueDate.IsZero() || len(invoice.LineItems) == 0 {
		response.Error(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	invoice.CalculateTotal()

	invoiceID, err := database.CreateInvoice(&invoice)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	invoice.ID = int(invoiceID)
	response.JSON(w, http.StatusCreated, invoice)
}

// GetInvoices lists all invoices.
func GetInvoices(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	invoices, err := database.GetInvoices(limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve invoices")
		return
	}

	response.JSON(w, http.StatusOK, invoices)
}


// GetInvoice retrieves a single invoice.
func GetInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/api/invoices/"):])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := database.GetInvoiceByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusNotFound, "Invoice not found")
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to retrieve invoice")
		}
		return
	}

	response.JSON(w, http.StatusOK, *invoice)
}

// UpdateInvoice updates an invoice.
func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/api/invoices/"):])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := database.UpdateInvoiceStatusString(id, payload.Status); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update invoice")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Invoice updated successfully"})
}

// CreateAdminUser creates a new admin user.
// TODO: Remove this endpoint or secure it properly in a production environment.
func CreateAdminUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if creds.Username == "" || creds.Password == "" {
		response.Error(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	hashedPassword, err := auth.HashPassword(creds.Password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := &database.User{
		Username:     creds.Username,
		PasswordHash: hashedPassword,
		IsAdmin:      true,
	}

	if _, err := database.CreateUser(user); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create admin user")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "Admin user created successfully"})
}
