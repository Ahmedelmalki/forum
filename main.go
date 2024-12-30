package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	forum "forum/app"

	_ "github.com/mattn/go-sqlite3"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func cleanupVisitors() {
	for {
		time.Sleep(1 * time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 1*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	go cleanupVisitors()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSeen) > 1*time.Minute {
			v.count = 1
			v.lastSeen = time.Now()
		} else {
			v.count++
		}
		mu.Unlock()

		if v.count > 5 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/newPost", rateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			forum.LogOutHandler(db)(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.Handle("/category", CategoryHandler(db))

	http.Handle("/static/style/", http.StripPrefix("/static/style/", http.FileServer(http.Dir("./static/style"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./static/js"))))

	fmt.Println("Server is running on http://localhost:8765")
	log.Fatal(http.ListenAndServe(":8765", nil))
}

func CategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the template file
		tmpl, err := template.ParseFiles("static/templates/category.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create a data structure to pass to the template
		data := map[string]interface{}{
			"Category": r.URL.Query().Get("category"), // Optional: current category
			"RawQuery": r.URL.RawQuery,                // Full raw query string
		}
		fmt.Println(data)
		// Execute the template with the data
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
