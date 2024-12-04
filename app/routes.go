package forum

import (
	"database/sql"
	"fmt"
	"net/http"
)

// Handler to process registration form submission
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password") // No encryption here

		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			http.Error(w, "Error adding user to database", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "User registered successfully!")
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Check if the user exists in the database
		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// For simplicity, we're not hashing passwords here
		if storedPassword != password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}
}

func LostsHandler(w http.ResponseWriter, r *http.Request) {
	 http.ServeFile(w, r, "static/posts.html")
}