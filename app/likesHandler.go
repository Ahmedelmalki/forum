package forum

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func HandleLikes(db *sql.DB) http.HandlerFunc {
	var like Likes

	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		target := "post"
		targetId := like.User_Id

		if like.CommentId != -1 {
			target = "comment"
			targetId = like.CommentId

		}
		switch r.Method {
		case http.MethodPost:
			{
				err := json.NewDecoder(r.Body).Decode(&like)
				if err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}
				like.User_Id, err = ValidateCookie(db, w, r)
				if err != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

				like.LikeCOunt, err = countLikesForPost(db, like.Post_Id, like.CommentId, like.Type, target)
				if err != nil {
					http.Error(w, "Error counting likes", http.StatusInternalServerError)
					return
				}
				checkQuery := `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND comment_id = ? AND user_id = ?)`

				var exists bool
				err = db.QueryRow(checkQuery, like.Post_Id, like.CommentId, like.User_Id).Scan(&exists)
				if err != nil {
					http.Error(w, "Error checking like existence", http.StatusInternalServerError)
					return
				}

				if exists {
					//  "SELECT TypeOfLike FROM likes WHERE post_id = ? AND user_id = ?"
					LiketypeQuery := `SELECT TypeOfLike FROM likes WHERE post_id = ? AND comment_id = ? AND user_id = ?`
					var typea string
					db.QueryRow(LiketypeQuery, like.Post_Id, like.CommentId, like.User_Id).Scan(&typea)
					if typea == like.Type {
						query := `DELETE FROM likes WHERE post_id = ? AND comment_id = ? AND user_id = ?`
						// query := "DELETE FROM likes WHERE post_id = ? AND user_id = ?"
						_, err = db.Exec(query, like.Post_Id, like.CommentId, like.User_Id)
						if err != nil {
							http.Error(w, "Error deleting like", http.StatusInternalServerError)
							return
						}
					} else {
						Updatequery := `UPDATE likes SET TypeOfLike = ? WHERE post_id = ? AND comment_id = ? AND user_id = ?`
						_, err = db.Exec(Updatequery, like.Type, like.Post_Id, like.CommentId, like.User_Id)
						if err != nil {
							http.Error(w, "Error UPDATNG likeS", http.StatusInternalServerError)
							return
						}
					}
				} else {
					if target == "post" {
					}
					query := "INSERT INTO likes (user_id, post_id , comment_id , TypeOfLike) VALUES (?, ?, ?, ?)"
					_, err = db.Exec(query, like.User_Id, like.Post_Id, like.CommentId, like.Type)
					if err != nil {
						fmt.Println(err)
						http.Error(w, "Error adding like", http.StatusInternalServerError)
						return
					}
				}

			}
		case http.MethodGet:
			{
				if targetId != 0 {
					like.User_Id, err = ValidateCookie(db, w, r)
					if err != nil {
						http.Error(w, "No Active Session", http.StatusInternalServerError)
						return
					}
					like.LikeCOunt, err = countLikesForPost(db, like.Post_Id, like.CommentId, "like", target)
					if err != nil {
						http.Error(w, "Error Counting like", http.StatusInternalServerError)
						return
					}
					like.DislikeCOunt, err = countLikesForPost(db, like.Post_Id, like.CommentId, "dislike", target)
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

func countLikesForPost(db *sql.DB, postID int, CommentId int, liketype string, target string) (int, error) {
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ? AND comment_id = ? AND TypeOfLike = ? `
	var likeCount int
	err := db.QueryRow(query, postID, CommentId, liketype).Scan(&likeCount)
	fmt.Println("########\n", likeCount, "\n#######")
	if err != nil {
		return 0, errors.New("error counting likes")
	}

	return likeCount, nil
}
