package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	forum "forum/app"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	scriptContent, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("error reading the schema")
	}

	_, err = db.Exec(string(scriptContent))
	if err != nil {
		log.Fatalf("Failed to execute SQL script: %v", err)
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "static/posts.html")
	// })

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/posts.html")
	})

	http.HandleFunc("/posts", forum.APIHandler(db))

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})

	http.HandleFunc("/register/submit", forum.RegisterHandler(db))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})

	http.HandleFunc("/login/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			forum.LoginHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/newPost", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := forum.ValidateCookie(db, w, r)
			if err != nil {
				fmt.Println("shit went wrong her1111111")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
				fmt.Println("shit went wrong here 2222222222222222222")
			http.ServeFile(w, r, "static/newPost.html")
		} else if r.Method == http.MethodPost {
				fmt.Println("shit went wrong here 333333333333333")
			forum.PostNewPostHandler(db)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			forum.LogOutHandler(db)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./static/js"))))

	fmt.Println("Server is running on http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
