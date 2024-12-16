package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Post struct {
	ID        int
	UserName  string
	Title     string
	Content   string
	Category  string
	CreatedAt time.Time
}

func FetchPosts(db *sql.DB) ([]Post, error) {
	query := "SELECT id, username, title, content, category, created_at FROM posts ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
		if err != nil {
			fmt.Printf("error scanning:\n %v \n", err)
			continue
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// APIHandler serves the posts as JSON
func APIHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user_id int
		cookie, err := r.Cookie("forum_session")
		fmt.Println("################--Cookie : ", cookie)
		if err != nil {
			user_id = 0
		} else {
			sessionID := cookie.Value
			query1 := `SELECT user_id FROM sessions WHERE session = ?`
			err = db.QueryRow(query1, sessionID).Scan(&user_id)
			if err != nil {
				user_id = 0
			}
		}
		fmt.Println("Error : ", err)
		posts, err := FetchPosts(db) // Use the FetchPosts function from the earlier example
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Println("##########-id", user_id)
		json.NewEncoder(w).Encode([]any{posts, user_id})
	}
}
