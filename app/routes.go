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
			fmt.Println("error 1")
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password") // No encryption here

		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
		_, err := db.Exec(query, username, email, password)
		if err != nil {
			http.Error(w, "Error adding user to database", http.StatusInternalServerError)
			fmt.Println("error 2")
			return
		}

		fmt.Println("User registered successfully!")
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
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
		query := "SELECT password FROM users WHERE email = ?"
		err := db.QueryRow(query, email).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
				fmt.Println("here")
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
		fmt.Printf("%s logged in successfully!\n", email)
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}
