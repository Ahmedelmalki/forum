package forum

import (
	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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
		password := r.FormValue("password") 

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "error hashing the password", http.StatusInternalServerError)
		}
		

		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		
		query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
		_, err = db.Exec(query, username, email, hashed)
		if err != nil {
			http.Error(w, "Error adding user to database", http.StatusInternalServerError)
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
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		err2 := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err2 != nil {
		    // If the password doesn't match the hash, respond with an Unauthorized status
		    http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		    return
		}
		
		fmt.Printf("%s logged in successfully!\n", email)
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}


