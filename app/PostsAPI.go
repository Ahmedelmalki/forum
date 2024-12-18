package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

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
		user_id := isLoged(db, r)
		posts, err := FetchPosts(db)
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]any{posts, user_id})
	}
}
