package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func CategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := r.URL.Query()["categories"]
		if len(categories) == 0 {
			ErrorHandler(w, "No categories provided", http.StatusBadRequest)
			return
		}
		query := `
					SELECT 
						p.id, 
						p.username, 
						p.title, 
						p.content, 
						p.created_at, 
						COALESCE(GROUP_CONCAT(c.categories), '') AS categories,
						COUNT(CASE WHEN l.TypeOfLike = 'like' AND comment_id == -1 THEN 1 ELSE NULL END) AS likes,
						COUNT(CASE WHEN l.TypeOfLike = 'dislike' AND comment_id == -1 THEN 1 ELSE NULL END) AS dislikes
					FROM 
						posts AS p
					LEFT JOIN 
						categories AS c ON c.post_id = p.id
					LEFT JOIN 
						likes AS l ON l.post_id = p.id
					WHERE 
						c.categories IN (?` + strings.Repeat(", ?", len(categories)-1) + `)
					GROUP BY 
						p.id;
		`

		args := make([]interface{}, len(categories))
		for i, category := range categories {
			args[i] = category
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			ErrorHandler(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
			fmt.Printf("Error executing query: %v\n", err)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			var categoryString string
			err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &post.CreatedAt, &categoryString, &post.Likes, &post.Dislikes)
			if err != nil {
				ErrorHandler(w, "Failed to parse posts: "+err.Error(), http.StatusInternalServerError)
				fmt.Printf("Error scanning rows: %v\n", err)
				return
			}

			if categoryString != "" {
				post.Categories = strings.Split(categoryString, ",")
			} else {
				post.Categories = []string{}
			}

			posts = append(posts, post)
		}

		if err := rows.Err(); err != nil {
			ErrorHandler(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
			fmt.Printf("Error iterating rows: %v\n", err)
			return
		}

		userID := isLoged(db, r)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]any{posts, userID}); err != nil {
			ErrorHandler(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
			fmt.Printf("Error encoding response: %v\n", err)
		}
	}
}

func Postcategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/templates/category.html")
		if err != nil {
			ErrorHandler(w, "Internal server error: failed to load template", http.StatusInternalServerError)
			fmt.Printf("Error loading template: %v\n", err)
			return
		}
		template.Execute(w, nil)
	}
}

func GetLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := isLoged(db, r)
		if user_id == 0 {
			ErrorHandler(w, "user Unauthorized to make this request", http.StatusUnauthorized)
			return
		}
		query := `	SELECT 
					    p.id AS post_id, 
					    p.username, 
					    p.title, 
					    p.content, 
					    p.created_at,
					    COALESCE(GROUP_CONCAT(c.categories), '') AS categories,
					    (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND TypeOfLike = 'like' AND comment_id == -1 ) AS like_count 
					FROM 
					    posts AS p
					INNER JOIN 
					    likes AS l
					ON 
					    l.post_id = p.id AND l.TypeOfLike = 'like' AND l.user_id = ?
					LEFT JOIN 
					    categories AS c
					ON 
					    c.post_id = p.id
					GROUP BY 
					    p.id, p.username, p.title, p.content, p.created_at
						`
		rows, err := db.Query(query, user_id)
		if err != nil {
			ErrorHandler(w, "Failed to query liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var posts []Post
		for rows.Next() {
			var post Post
			var category string
			var likeCount int

			err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &post.CreatedAt, &category, &likeCount)
			if err != nil {
				ErrorHandler(w, "Failed to parse liked posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			post.Categories = strings.Split(category, ",")
			post.Likes = likeCount
			posts = append(posts, post)
		}
		if err := rows.Err(); err != nil {
			ErrorHandler(w, "Error reading liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			ErrorHandler(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func PostLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/templates/likedPosts.html")
		if err != nil {
			ErrorHandler(w, "internal server error", 500)
		}
		template.Execute(w, nil)
	}
}

func GetPostsCreatedBy(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := isLoged(db, r)
		if user_id == 0 {
			ErrorHandler(w, "user Unauthorized to make this request", http.StatusUnauthorized)
			return
		}
		query := `SELECT 
					    p.id, 
					    p.username, 
					    p.title, 
					    p.content, 
					    p.created_at,
					    COALESCE(GROUP_CONCAT(c.categories), '') AS categories,
					    (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND TypeOfLike = 'like' AND comment_id == -1) AS like_count
					FROM 
					    posts AS p
					LEFT JOIN 
					    likes AS l
					ON 
					    l.post_id = p.id AND l.TypeOfLike = 'like'
					LEFT JOIN 
					    categories AS c
					ON 
					    c.post_id = p.id
					WHERE 
					    p.username = (SELECT username FROM users WHERE id = ?)
					GROUP BY 
					    p.id;
						`
		rows, err := db.Query(query, user_id)
		if err != nil {
			ErrorHandler(w, "Failed to query liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var posts []Post
		for rows.Next() {
			var post Post
			var category string
			var likeCount int

			err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &post.CreatedAt, &category, &likeCount)
			if err != nil {
				ErrorHandler(w, "Failed to parse liked posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			post.Categories = strings.Split(category, ",")
			post.Likes = likeCount
			posts = append(posts, post)
		}
		if err := rows.Err(); err != nil {
			ErrorHandler(w, "Error reading liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			ErrorHandler(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func PostPostsCreatedBy(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/templates/createdBy.html")
		if err != nil {
			ErrorHandler(w, "internal server error", 500)
		}
		template.Execute(w, nil)
	}
}
