package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func FetchPosts(db *sql.DB, category string) ([]Post, error) {
	baseQuery := `SELECT 
					p.id,
					p.username,
					p.title,
					p.content,
					p.category,
					p.created_at,
					COALESCE(
						(SELECT COUNT(*) FROM likes
						WHERE post_id = p.id AND typeOfLike = 'like'),
						0
						) as likes,
            		COALESCE(
                		(SELECT COUNT(*) FROM likes 
                		WHERE post_id = p.id AND TypeOfLike = 'dislike'),
                		0
            			) as dislikes
			 	FROM posts p
						 `
	var rows *sql.Rows
	var err error

	if category != "" && category != "all" {
		query := baseQuery + " WHERE p.category = ? ORDER BY p.created_at DESC"
		rows, err = db.Query(query, category)
	} else {
		query := baseQuery + " ORDER BY p.created_at DESC"
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UserName,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.CreatedAt,
			&post.Likes,
			&post.Dislikes,
		)
		if err != nil {
			fmt.Printf("error scanning: %v\n", err)
			continue
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return posts, nil
}

// APIHandler serves the posts as JSON
func APIHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		user_id := isLoged(db, r)
		posts, err := FetchPosts(db, category)
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		fmt.Println(json.NewEncoder(os.Stdout).Encode([]any{posts, user_id}))
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
	}
}
