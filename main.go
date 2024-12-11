package main

import (
	"database/sql"
	"fmt"
	forum "forum/app"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./forum.db") ; if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
    log.Fatalf("Database connection error: %v", err)
	}

	 scriptContent , err := os.ReadFile("schema.sql"); if err != nil {
		log.Fatal("error reading the schema")
	}

	_, err = db.Exec(string(scriptContent)) ; if err != nil {
		log.Fatalf("Failed to execute SQL script: %v", err)
	}


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/posts.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})

	http.HandleFunc("/register/submit", forum.RegisterHandler(db))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	http.HandleFunc("/login/submit", forum.LoginHandler(db))

	http.HandleFunc("/newPost", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./static/newPost.html")
	})

	http.HandleFunc("/newPost/submit", forum.NewPostHandler(db))


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	

	fmt.Println("Server is running on http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
