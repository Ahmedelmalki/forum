package forum

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type likes struct {
	User_Id      int    `json:"UserId"`
	Post_Id      int    `json:"PostId"`
	LikeCount    int    `json:"LikeCOunt"`
	DisLikeCount int    `json:"DislikeCOunt"`
	Type         string `json:"Type"`
}

func HandleLikes(db *sql.DB) http.HandlerFunc {
	var like likes

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.Method {
		case http.MethodPost:
			{
				err := json.NewDecoder(r.Body).Decode(&like)
				if err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				like.User_Id, err = ValidateCookie(db, w, r)
				// fmt.Println("User id  now:",like.User_Id)
				if err != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
				like.LikeCount, err = countLikesForPost(db, like.Post_Id, like.Type)
				if err != nil {
					http.Error(w, "Error counting likes", http.StatusInternalServerError)
					return
				}
				checkQuery := "SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ?)"
				var exists bool
				err = db.QueryRow(checkQuery, like.Post_Id, like.User_Id).Scan(&exists)
				if err != nil {
					http.Error(w, "Error checking like existence", http.StatusInternalServerError)
					return
				}

				if exists {
					LiketypeQuery := "SELECT TypeOfLike FROM likes WHERE post_id = ? AND user_id = ?"
					var typea string
					db.QueryRow(LiketypeQuery, like.Post_Id, like.User_Id).Scan(&typea)
					if typea == like.Type {
						query := "DELETE FROM likes WHERE post_id = ? AND user_id = ?"
						_, err = db.Exec(query, like.Post_Id, like.User_Id)
						if err != nil {
							http.Error(w, "Error deleting like", http.StatusInternalServerError)
							return
						}
					} else {
						Updatequery := "UPDATE likes SET TypeOfLike = ? WHERE post_id = ? AND user_id = ?"
						_, err = db.Exec(Updatequery, like.Type, like.Post_Id, like.User_Id)
						if err != nil {
							http.Error(w, "Error UPDATNG likeS", http.StatusInternalServerError)
							return
						}
					}
				} else {
					query := "INSERT INTO likes (user_id, post_id , TypeOfLike) VALUES (?, ?, ?)"
					_, err = db.Exec(query, like.User_Id, like.Post_Id, like.Type)
					if err != nil {
						http.Error(w, "Error adding like", http.StatusInternalServerError)
						return
					}
				}

			}
		case http.MethodGet:
			{
				if like.Post_Id != 0 {
					like.User_Id, err = ValidateCookie(db, w, r)
					if err != nil {
						http.Error(w, "No Active Session", http.StatusInternalServerError)
						return
					}
					like.LikeCount, err = countLikesForPost(db, like.Post_Id, "like")
					if err != nil {
						http.Error(w, "Error Counting like", http.StatusInternalServerError)
						return
					}
					like.DisLikeCount, err = countLikesForPost(db, like.Post_Id, "dislike")
					if err != nil {
						http.Error(w, "Error Counting dislike", http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(&like)
				}
			}

		default:
			{
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}
		}
	}
}

func countLikesForPost(db *sql.DB, postID int, liketype string) (int, error) {
	query := "SELECT COUNT(*) FROM likes WHERE post_id = ? AND TypeOfLike = ? "
	var LikeCount int
	err := db.QueryRow(query, postID, liketype).Scan(&LikeCount)
	if err != nil {
		return 0, err
	}

	return LikeCount, nil
}
