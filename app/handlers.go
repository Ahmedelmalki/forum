package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
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

		// handling one session at a time
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
		fmt.Printf("%s logged in successfully!\n", email)
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

func PostNewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}
		cookie, err := r.Cookie("forum_session")
		if err != nil {
			http.Error(w, "Unauthorized to create a post", http.StatusUnauthorized)
			return
		}
		var idForUsername int
		err = db.QueryRow(`SELECT user_id FROM sessions WHERE session= ?;`, cookie.Value).Scan(&idForUsername)
		fmt.Println(idForUsername)
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
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		if len(title) > 50 || len(content) > 1000 {
			http.Error(w, "don't miss with the html plz", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin() // starting a transaction
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

func CategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := r.URL.Query()["categories"] // Extract categories[] parameter from the query string
		fmt.Println(categories)

		if len(categories) == 0 {
			http.Error(w, "No categories provided", http.StatusBadRequest)
			return
		}
		query := `SELECT 
					    p.id, 
					    p.username, 
					    p.title, 
					    p.content, 
					    p.created_at,
						c.categories
					FROM 
					    posts AS p
					INNER JOIN 
					    categories AS c 
					ON 
					    c.post_id = p.id
					WHERE 
					    c.categories = ?;
						`
		var posts []Post
		for _, c := range categories {
			rows, err := db.Query(query, c)
			if err != nil {
				http.Error(w, "Failed to query posts: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()
			for rows.Next() {
				var post Post
				var category string

				err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &post.CreatedAt, &category)
				if err != nil {
					http.Error(w, "Failed to parse posts: "+err.Error(), http.StatusInternalServerError)
					return
				}

				post.Categories = strings.Split(category, ",")
				posts = append(posts, post)
			}
			if err := rows.Err(); err != nil {
				http.Error(w, "Error reading posts: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		user_id := isLoged(db, r)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			http.Error(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func Postcategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/templates/category.html")
		if err != nil {
			http.Error(w, "internal", 500)
		}
		template.Execute(w, nil)
	}
}
