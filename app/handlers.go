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
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var credentials RegisterCredenials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		username := credentials.UserName
		email := credentials.Email
		password := credentials.Password

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "error hashing the password", http.StatusInternalServerError)
		}

		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
		res, err := db.Exec(query, username, email, hashed)
		if err != nil {
			http.Error(w, "Error adding user to database", http.StatusInternalServerError)
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
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var credentials LoginCredentials

		// Decode JSON body
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		email := credentials.Email
		password := credentials.Password

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}
		var storedPassword string
		var user_id int
		query := "SELECT password, id FROM users WHERE email = ?"
		err = db.QueryRow(query, email).Scan(&storedPassword, &user_id)
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

		//handling one session at a time
        deleteQuery := "DELETE FROM sessions WHERE user_id = ?"
        _, err = db.Exec(deleteQuery, user_id)
        if err != nil {
            http.Error(w, "Error cleaning old sessions", http.StatusInternalServerError)
            return
        }
		
		cookie := CookieMaker(w)
		err = InsretCookie(db, user_id, cookie, time.Now().Add(time.Hour*24))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(user_id, cookie)
		fmt.Printf("%s logged in successfully!\n", email)
	}
}

func GetNewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/newPost.html")
	}
}

func PostNewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}
		cookie, err := r.Cookie("forum_session")
		if err != nil {
			return
		}
		var idForUsername int
		// stupid misstake : i was slecting id instad of user_id which is linked to the users table
		query1 := `SELECT user_id FROM sessions WHERE session= ?;`
		err = db.QueryRow(query1, cookie.Value).Scan(&idForUsername)
		fmt.Println(idForUsername)
		if err != nil {
			fmt.Println(err)
			return
		}
		var userName string
		query2 := `SELECT username FROM users WHERE id= ?`
		err = db.QueryRow(query2, idForUsername).Scan(&userName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(userName)

		title := r.FormValue("title")
		category := r.FormValue("category")
		content := r.FormValue("content")
		if title == "" || category == "" || content == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		if len(title) > 50 || len(content) > 1000 {
			http.Error(w, "don't miss with the html plz", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO posts (username, title, category, content) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(query, userName, title, category, content)
		if err != nil {
			log.Printf("Error adding post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogOutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("forum_session")
		if err != nil {
			http.Error(w, "No active session found", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		fmt.Printf("Method: %s Cookie: %+v\n", r.Method, cookie)
		query := `DELETE FROM sessions WHERE session = ?`
		res, err := db.Exec(query, sessionID)
		if err != nil {
			fmt.Println("error executing the query")
		}
		rows, _ := res.RowsAffected()
		fmt.Println(rows)

		http.SetCookie(w, &http.Cookie{
			Name:   "forum_session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
