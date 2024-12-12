package forum

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func InsretCookie(db *sql.DB, user_id int, cookie string) error {
	query := `INSERT INTO sessions (user_id, session) VALUES (?, ?)`
	_, err := db.Exec(query, user_id, cookie) ; if err != nil {
		return err
	}
	
	return nil
}

func CookieMaker(w http.ResponseWriter) string{
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