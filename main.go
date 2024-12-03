package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Register the /login route and attach the login handler
	http.HandleFunc("/login", loginHandler(db))

	// Other routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Forum!"))
	})

	// Start server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Serve the login form (HTML)
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, `
				<form action="/login" method="POST">
					<input type="email" name="email" placeholder="Email" required />
					<input type="password" name="password" placeholder="Password" required />
					<button type="submit">Login</button>
				</form>
			`)
			return
		}

		if r.Method == http.MethodPost {
			// Fetch login details from the form
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Fetch the user from the database
			var hashedPassword string
			err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				} else {
					log.Printf("Error querying user: %v\n", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}

			// Compare the provided password with the hashed password
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
			if err != nil {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			// Successful login
			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: "some-unique-token", // Ideally, generate a secure session token
				Path:  "/",
			})
			fmt.Fprintln(w, "Login successful!")
		}
	}
}