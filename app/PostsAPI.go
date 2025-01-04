package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func FetchPosts(db *sql.DB, category string) ([]Post, error) {
	baseQuery := `
					SELECT 
					    p.id,
					    p.username,
					    p.title,
					    p.content,
					    COALESCE(GROUP_CONCAT(c.categories, ','), '') AS categories,
					    p.created_at,
					    COALESCE(
					        (SELECT COUNT(*) FROM likes 
					         WHERE post_id = p.id AND typeOfLike = 'like' AND comment_id = -1), 
					        0
					    ) AS likes,
					    COALESCE(
					        (SELECT COUNT(*) FROM likes 
					         WHERE post_id = p.id AND typeOfLike = 'dislike' AND comment_id = -1), 
					        0
					    ) AS dislikes
					FROM posts p
					LEFT JOIN categories c ON c.post_id = p.id
					`
	var rows *sql.Rows
	var err error

	if category != "" && category != "all" {
		query := baseQuery + `
		WHERE p.id IN (
			SELECT post_id 
			FROM categories 
			WHERE categories = ?
		)
		GROUP BY p.id
		ORDER BY p.created_at DESC
		`
		rows, err = db.Query(query, category)
	} else {
		query := baseQuery + `
		GROUP BY p.id
		ORDER BY p.created_at DESC
		`
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var categoryString string
		err := rows.Scan(
			&post.ID,
			&post.UserName,
			&post.Title,
			&post.Content,
			&categoryString,
			&post.CreatedAt,
			&post.Likes,
			&post.Dislikes,
		)
		if err != nil {
			fmt.Printf("error scanning: %v\n", err)
			continue
		}
		if categoryString != "" {
			post.Categories = splitStringByComma(categoryString)
		} else {
			post.Categories = []string{}
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return posts, nil
}

func splitStringByComma(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, ",")
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
		/* this line sets an HTTP response header to control how the response
		is cached by clients (browsers) and intermediate caches (proxies).*/
		w.Header().Set("Cache-Control", "no-cache")
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
	}
}

func CategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := r.URL.Query()["categories[]"] // Extract categories[] parameter from the query string

		// If no categories are selected, return an error or handle it as "all"
		if len(categories) == 0 {
			http.Error(w, "No categories provided", http.StatusBadRequest)
			return
		}

		// Construct a dynamic SQL query for SQLite
		query := `SELECT 
    p.id, 
    p.username, 
    p.title, 
    p.content, 
    p.created_at,
	c.categories
FROM 
    posts AS p
INNER JOIN 
    categories AS c 
ON 
    c.post_id = p.id
WHERE 
    c.categories = ?;
						`
		var posts []Post
		for _, c := range categories {
			rows, err := db.Query(query, c)
			if err != nil {
				http.Error(w, "Failed to query posts: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// Parse the rows into a slice of Post structs
			for rows.Next() {
				var post Post
				var category string

				// Adjust the scan fields to match your database schema
				err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &category, &post.CreatedAt)
				if err != nil {
					http.Error(w, "Failed to parse posts: "+err.Error(), http.StatusInternalServerError)
					return
				}

				// Split the category string into an array if it's stored as a comma-separated value
				post.Categories = strings.Split(category, ",")
				posts = append(posts, post)
			}
			if err := rows.Err(); err != nil {
				http.Error(w, "Error reading posts: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			http.Error(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
