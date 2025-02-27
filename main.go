package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	forum "forum/app"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./finalTest.db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

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
	return db
}

func main() {
	db := initDB()
	defer db.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
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
			_, err := forum.ValidateCookie(db, w, r)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			http.ServeFile(w, r, "static/templates/newPost.html")
		} else if r.Method == http.MethodPost {
			forum.PostNewPostHandler(db)(w, r)
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

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
		} else if r.Method == "POST" {
			forum.CategoryHandler(db)(w, r)
			return
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/LikedPosts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			forum.PostLikedPosts(db)(w, r)
			return
		} else if r.Method == "POST" {
			forum.GetLikedPosts(db)(w, r)
			return
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/CtreatedBy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			forum.PostPostsCreatedBy(db)(w, r)
			return
		} else if r.Method == "POST" {
			forum.GetPostsCreatedBy(db)(w, r)
			return
		} else {
			forum.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/static/style/", http.StripPrefix("/static/style/", http.FileServer(http.Dir("./static/style"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./static/js"))))

	fmt.Println("Server is running on http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
