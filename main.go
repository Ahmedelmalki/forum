package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	forum "forum/app"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./testingFilter.db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	scriptFile, err := os.Open("schema.sql")
	if err != nil {
		log.Fatalf("Failed to open SQL script file: %v", err)
	}
	defer scriptFile.Close()

	scriptContent, err := io.ReadAll(scriptFile)
	if err != nil {
		log.Fatalf("Failed to read SQL script file: %v", err)
	}
	_, err = db.Exec(string(scriptContent))
	if err != nil {
		log.Fatalf("Failed to execute SQL script: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"{
			forum.ErrorHandler(w, "page not found", 404)
			return
		}
		http.ServeFile(w, r, "static/templates/posts.html")
	})

	http.HandleFunc("/posts", forum.APIHandler(db))

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/register.html")
	})

	http.HandleFunc("/register/submit", forum.RegisterHandler(db))

	// likes
	http.HandleFunc("/like", forum.HandleLikes(db))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/login.html")
	})

	http.HandleFunc("/login/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			forum.LoginHandler(db)(w, r)
		} else {
			forum.ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/newPost", forum.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Validate the user's cookie
			_, err := forum.ValidateCookie(db, w, r)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			// Serve the HTML template for new posts
			http.ServeFile(w, r, "static/templates/newPost.html")
		} else if r.Method == http.MethodPost {
			// Handle new post creation
			forum.PostNewPostHandler(db)(w, r)
		} else {
			// Return a 405 Method Not Allowed error for unsupported methods
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Comments handling
	http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			forum.CreateComment(db)(w, r)
		case http.MethodGet:
			forum.GetComments(db)(w, r)
		default:
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			forum.LogOutHandler(db)(w, r)
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/category", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			forum.Postcategory(db)(w, r)
			return
		}else if r.Method == "POST" {
			forum.CategoryHandler(db)(w, r)
			return
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//filtering by liked posts
	http.HandleFunc("/LikedPosts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			forum.PostLikedPosts(db)(w, r)
			fmt.Println("calling PostLikedPosts")
			return
		}else if r.Method == "POST" {
			forum.GetLikedPosts(db)(w, r)
			return
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//filtering by created posts
	http.HandleFunc("/CtreatedBy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			forum.PostPostsCreatedBy(db)(w, r)
			return
		}else if r.Method == "POST" {
			forum.GetPostsCreatedBy(db)(w, r)
			return
		} else{
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/static/style/", http.StripPrefix("/static/style/", http.FileServer(http.Dir("./static/style"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./static/js"))))

	fmt.Println("Server is running on http://localhost:9898")
	log.Fatal(http.ListenAndServe(":9898", nil))
}
