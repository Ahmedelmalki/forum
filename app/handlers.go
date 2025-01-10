package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var credentials RegisterCredenials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			ErrorHandler(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		username := credentials.UserName
		email := credentials.Email
		password := credentials.Password

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			ErrorHandler(w, "error hashing the password", http.StatusInternalServerError)
		}

		if username == "" || email == "" || password == "" {
			ErrorHandler(w, "All fields are required", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
		res, err := db.Exec(query, username, email, hashed)
		if err != nil {
			ErrorHandler(w, "Error adding user to database", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		user_id, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
			return
		}
		cookie := CookieMaker(w)
		err = InsretCookie(db, int(user_id), cookie, time.Now().Add(time.Hour*24))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s registered successfully\n", email)
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var credentials LoginCredentials

		// Decode JSON body
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			ErrorHandler(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		email := credentials.Email
		password := credentials.Password

		if email == "" || password == "" {
			ErrorHandler(w, "Email and password are required", http.StatusBadRequest)
			return
		}
		var storedPassword string
		var user_id int
		query := "SELECT password, id FROM users WHERE email = ?"
		err = db.QueryRow(query, email).Scan(&storedPassword, &user_id)
		if err != nil {
			if err == sql.ErrNoRows {
				ErrorHandler(w, "User not found", http.StatusUnauthorized)
			} else {
				ErrorHandler(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		err2 := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err2 != nil {
			ErrorHandler(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// handling one session at a time
		deleteQuery := "DELETE FROM sessions WHERE user_id = ?"
		_, err = db.Exec(deleteQuery, user_id)
		if err != nil {
			ErrorHandler(w, "Error cleaning old sessions", http.StatusInternalServerError)
			return
		}
		
		cookie := CookieMaker(w)
		err = InsretCookie(db, user_id, cookie, time.Now().Add(time.Hour*24))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s logged in successfully!\n", email)
	}
}

func AuthenticationMiddleware(next http.HandlerFunc, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow access to the homepage without authentication
		if r.URL.Path == "/" {
			next(w, r)
			return
		}

		// Check if the user has a valid session
		cookie, err := r.Cookie("forum_session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Verify session in the database
		var userID int
		err = db.QueryRow("SELECT user_id FROM sessions WHERE session = ?", cookie.Value).Scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Proceed to the handler if authenticated
		next(w, r)
	}
}

func GetNewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/newPost.html")
	}
}

func LogOutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("forum_session")
		if err != nil {
			ErrorHandler(w, "No active session found", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		query := `DELETE FROM sessions WHERE session = ?`
		res, err := db.Exec(query, sessionID)
		if err != nil {
			fmt.Println("error executing the query")
		}
		rows, _ := res.RowsAffected()
		fmt.Println("rows :", rows)

		http.SetCookie(w, &http.Cookie{
			Name:   "forum_session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func PostNewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			ErrorHandler(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}
		cookie, err := r.Cookie("forum_session")
		if err != nil {
			ErrorHandler(w, "Unauthorized to create a post", http.StatusUnauthorized)
			return
		}
		var idForUsername int
		err = db.QueryRow(`SELECT user_id FROM sessions WHERE session= ?;`, cookie.Value).Scan(&idForUsername)
		if err != nil {
			fmt.Println(err)
			return
		}
		var userName string
		err = db.QueryRow(`SELECT username FROM users WHERE id= ?`, idForUsername).Scan(&userName)
		if err != nil {
			fmt.Println(err)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["categories[]"]
		if title == "" || len(categories) == 0 || content == "" {
			ErrorHandler(w, "All fields are required", http.StatusBadRequest)
			return
		}
		if len(title) > 50 || len(content) > 1000 {
			ErrorHandler(w, "title leght or content lenght exceeded", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback()
		// inserting into posts
		result, err := tx.Exec("INSERT INTO posts (username, title, content) VALUES (?, ?, ?)", userName, title, content)
		if err != nil {
			return
		}
		postID, err := result.LastInsertId()
		if err != nil {
			return
		}
		for _, category := range categories {
			_, err = tx.Exec("INSERT INTO categories (post_id, categories) VALUES (?, ?)", postID, category)
			if err != nil {
				return
			}
		}
		tx.Commit()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
