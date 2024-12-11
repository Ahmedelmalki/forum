package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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

		fmt.Printf("%s registered successfully\n", email)
		cookieMaker(w)
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
		    http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		    return
		}
		
		fmt.Printf("%s logged in successfully!\n", email)
		cookieMaker(w)
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}


func cookieMaker(w http.ResponseWriter) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	cookie := &http.Cookie{
		Name:  "forum_session",
		Value: u.String(),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func NewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			http.ServeFile(w, r, "./static/newPost.html")
		case http.MethodPost:
			title := r.FormValue("title")
			category := r.FormValue("category")
			content := r.FormValue("content")

			if title == "" || category == "" || content == "" {
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			query := "INSERT INTO posts (title, category, content) VALUES (?, ?, ?)"
			_, err := db.Exec(query, title, category, content)
			if err != nil {
				log.Printf("Error adding post: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}



