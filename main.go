package main

import (
	"database/sql"
	"fmt"
	forum "forum/app"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Route to serve the home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/home.html")
	})

	// Route to serve the registration page
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})

	// Route to handle registration form submission
	http.HandleFunc("/register/submit", forum.RegisterHandler(db))

	// Handlers
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})
	http.HandleFunc("/login/submit", forum.LoginHandler(db))
	http.HandleFunc("/posts", forum.LostsHandler)

	// Start server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
