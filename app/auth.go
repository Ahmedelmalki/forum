package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func InsretCookie(db *sql.DB, user_id int, cookie string) error {
	query := `INSERT INTO sessions (user_id, session) VALUES (?, ?)`
	_, err := db.Exec(query, user_id, cookie)
	if err != nil {
		return err
	}

	return nil
}

func CookieMaker(w http.ResponseWriter) string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	cookie := &http.Cookie{
		Name:  "forum_session",
		Value: u.String(),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	return u.String()
}

func ValidateCookie(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	cookie, err := r.Cookie("forum_session")
	if err != nil {
		return 0, errors.New("error")
	}
	sessionID := cookie.Value
	query1 := `SELECT user_id FROM sessions WHERE session = ?`
	var user_id int
	err = db.QueryRow(query1, sessionID).Scan(&user_id)
	if err != nil {
		log.Printf("Failed to validate session for GET: %v", err)
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return 0, errors.New("error")
	}
	fmt.Println("test")
	return user_id, nil
}

func LogOutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Retrieve the session cookie
		cookie, err := r.Cookie("forum_session")
		if err != nil {
			// If no session cookie is found, user is already logged out
			http.Error(w, "No active session found", http.StatusUnauthorized)
			return
		}

		// Invalidate the session in the database
		sessionID := cookie.Value
		query := `DELETE FROM sessions WHERE session = ?`
		_, err = db.Exec(query, sessionID)
		if err != nil {
			log.Printf("Error invalidating session: %v", err)
			http.Error(w, "Failed to log out", http.StatusInternalServerError)
			return
		}

		// Clear the cookie from the user's browser
		http.SetCookie(w, &http.Cookie{
			Name:  "forum_session",
			Value: "",
			Path:  "/",
		})

		// Redirect to the login page or a confirmation page
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}
