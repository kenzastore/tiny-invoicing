# Tiny Invoicing 🚀

Tiny Invoicing is a lightweight, high-performance invoicing application built with **Go** and **MySQL**. It features a modern, responsive dashboard UI powered by **Bootstrap 5** and robust backend logic designed for reliability.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg?logo=go)
![Database](https://img.shields.io/badge/MySQL-8.0+-4479A1.svg?logo=mysql)

## ✨ Key Features

*   **Modern Dashboard UI:** A clean, responsive interface using Bootstrap 5 and Inter font for managing invoices.
*   **Robust Invoice Generation:** Create invoices with multiple line items, automatic total calculation, and validation.
*   **Auto-Healing Database Logic:** The system automatically handles missing customer data dependencies to prevent Foreign Key errors during demos.
*   **Secure Authentication:** Basic Authentication implementation for secure access.
*   **RESTful API:** Clean JSON API backend that can be consumed by any frontend client.
*   **RFC3339 Time Standardization:** Accurate date/time handling between frontend and backend.

## 🛠️ Tech Stack

*   **Backend:** Go (Golang) `net/http` standard library.
*   **Database:** MySQL (using `go-sql-driver`).
*   **Frontend:** HTML5, CSS3, Bootstrap 5, JavaScript (Fetch API).
*   **Architecture:** Clean architecture separating Models, Handlers, and Database layers.

## ⚙️ Installation & Setup

### Prerequisites
*   Go 1.20 or higher
*   MySQL Server

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/tiny-invoicing.git
cd tiny-invoicing
```

### 2. Database Setup
Create a MySQL database named `tiny_invoicing` (or your preferred name). You can import the initial schema, although the application includes some auto-migration features.

```sql
CREATE DATABASE tiny_invoicing;
```

*(Optional) Import `schema.sql` if you want a fresh start manually.*

### 3. Configure Environment
Set the `DB_DSN` environment variable to point to your MySQL instance.
**Format:** `user:password@tcp(localhost:3306)/dbname?parseTime=true`

**Linux/Mac:**
```bash
export DB_DSN="root:password@tcp(127.0.0.1:3306)/tiny_invoicing?parseTime=true"
```

**Windows (PowerShell):**
```powershell
$env:DB_DSN="root:password@tcp(127.0.0.1:3306)/tiny_invoicing?parseTime=true"
```

### 4. Run the Application
```bash
go run .
```
The server will start on **port 8080**.

## 📖 Usage

1.  Open your browser and navigate to `http://localhost:8080`.
2.  **First Time Login:**
    *   Click **"Setup Admin Account"**.
    *   Click **"Register Admin"** (default: `admin` / `password`).
    *   Go back to Login and sign in.
3.  **Dashboard:**
    *   View your invoice history on the left.
    *   Create new invoices on the right.
    *   **Note:** The system automatically creates a "Demo Client" (ID 1) if no customers exist, ensuring your first invoice always succeeds.

## 🔌 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/invoices` | Retrieve a list of all invoices |
| `GET` | `/api/invoices/{id}` | Get details of a specific invoice |
| `POST` | `/api/invoices` | Create a new invoice |
| `PUT` | `/api/invoices/{id}` | Update invoice status |
| `POST` | `/api/admin/create-user` | Register a new admin user |

## 📂 Project Structure

```
tiny-invoicing/
├── conductor/       # Project management & docs (Conductor)
├── database/        # Database connection & logic
├── handlers/        # HTTP Request handlers
├── models/          # Go structs for DB entities
├── static/          # Frontend assets (HTML/JS/CSS)
├── main.go          # Entry point
├── schema.sql       # Database schema
└── go.mod           # Go dependencies
```

## 📝 License

This project is licensed under the MIT License.
