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
		categories := r.URL.Query()["categories"] // Extract categories[] parameter from the query string
		fmt.Println(categories)

		if len(categories) == 0 {
			http.Error(w, "No categories provided", http.StatusBadRequest)
			return
		}
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
			for rows.Next() {
				var post Post
				var category string

				err := rows.Scan(&post.ID, &post.UserName, &post.Title, &post.Content, &post.CreatedAt, &category)
				if err != nil {
					http.Error(w, "Failed to parse posts: "+err.Error(), http.StatusInternalServerError)
					return
				}

				post.Categories = strings.Split(category, ",")
				posts = append(posts, post)
			}
			if err := rows.Err(); err != nil {
				http.Error(w, "Error reading posts: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		user_id := isLoged(db, r)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			http.Error(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func Postcategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/templates/category.html")
		if err != nil {
			http.Error(w, "internal server error", 500)
		}
		template.Execute(w, nil)
	}
}

func GetLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `SELECT 
			    p.id, 
			    p.username, 
			    p.title, 
			    p.content, 
			    p.created_at,
			    COALESCE(GROUP_CONCAT(c.categories), '') AS categories,
			    COUNT(l.id) AS like_count
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
			    l.id IS NOT NULL
			GROUP BY 
			    p.id
						`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Failed to query liked posts: "+err.Error(), http.StatusInternalServerError)
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
				http.Error(w, "Failed to parse liked posts: "+err.Error(), http.StatusInternalServerError)
				return
			}

			post.Categories = strings.Split(category, ",")
			post.Likes = likeCount
			posts = append(posts, post)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error reading liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//
		user_id := isLoged(db, r)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]any{posts, user_id}); err != nil {
			http.Error(w, "Failed to encode posts: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func PostLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/templates/likedPosts.html")
		if err != nil {
			http.Error(w, "internal server error", 500)
		}
		template.Execute(w, nil)
	}
}
