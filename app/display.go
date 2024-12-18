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

type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserName  string    `json:"username"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
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
		posts, err := FetchPosts(db) 
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]any{posts, user_id})
	}
}

/*************Start Comment handler functions*****************/
// Function to create a new comment
func CreateComment(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        
		userID, authenticated := ValidateCookie(db, w,r)
        if authenticated != nil {
            http.Error(w, authenticated.Error(), http.StatusUnauthorized)
            return
        }

		var comment Comment
        
        err := json.NewDecoder(r.Body).Decode(&comment)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        stmt, err := db.Prepare("INSERT INTO comments(post_id, user_id, content, created_at) VALUES(?, ?, ?, ?)")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer stmt.Close()

		comment.UserID = userID

		if comment.Content == "" {
			http.Error(w, "Need to add a comment", http.StatusBadRequest)
			return
		}

        _, err = stmt.Exec(comment.PostID, comment.UserID, comment.Content, time.Now())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

// Function to get comments
func GetComments(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        postID := r.URL.Query().Get("post_id")
        rows, err := db.Query(`SELECT com.id, com.user_id, us.username, com.content, com.created_at FROM comments com
            JOIN users us ON com.user_id = us.id
            WHERE com.post_id = ?
            ORDER BY com.created_at ASC
        `, postID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var comments []Comment
        for rows.Next() {
            var comment Comment
            err := rows.Scan(&comment.ID, &comment.UserID,&comment.UserName, &comment.Content, &comment.CreatedAt)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            comments = append(comments, comment)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(comments)
    }
}
/*************End Comment handler functions*****************/
