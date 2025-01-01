package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type likes struct {
	User_Id      int    `json:"UserId"`
	Post_Id      int    `json:"PostId"`
	LikeCOunt    int    `json:"LikeCOunt"`
	DislikeCOunt int    `json:"DislikeCOunt"`
	CommentId    int    `json:"CommentId"`
	Type         string `json:"Type"`
}

func HandleLikes(db *sql.DB) http.HandlerFunc {
	var like likes

	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		target := "post"
		if like.CommentId != -1 {
			target = "comment"
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
				// fmt.Println("User id  now:",like.User_Id)
				if err != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

				like.LikeCOunt, err = countLikesForPost(db, like.Post_Id, like.Type, target)
				if err != nil {
					http.Error(w, "Error counting likes", http.StatusInternalServerError)
					return
				}
				checkQuery := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM likes WHERE %s_id = ? AND user_id = ?)`, target)

				var exists bool
				err = db.QueryRow(checkQuery, like.Post_Id, like.User_Id).Scan(&exists)
				if err != nil {
					http.Error(w, "Error checking like existence", http.StatusInternalServerError)
					return
				}

				if exists {
					//  "SELECT TypeOfLike FROM likes WHERE post_id = ? AND user_id = ?"
					LiketypeQuery := fmt.Sprintf(`SELECT TypeOfLike FROM likes WHERE %s_id = ? AND user_id = ?`, target)
					var typea string
					db.QueryRow(LiketypeQuery, like.Post_Id, like.User_Id).Scan(&typea)
					if typea == like.Type {
						query := fmt.Sprintf(`DELETE FROM likes WHERE %s_id = ? AND user_id = ?`, target)
						// query := "DELETE FROM likes WHERE post_id = ? AND user_id = ?"
						_, err = db.Exec(query, like.Post_Id, like.User_Id)
						if err != nil {
							http.Error(w, "Error deleting like", http.StatusInternalServerError)
							return
						}
					} else {
						Updatequery := fmt.Sprintf(`UPDATE likes SET TypeOfLike = ? WHERE %s_id = ? AND user_id = ?`, target)
						_, err = db.Exec(Updatequery, like.Type, like.Post_Id, like.User_Id)
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
				if like.Post_Id != 0 {
					like.User_Id, err = ValidateCookie(db, w, r)
					if err != nil {
						http.Error(w, "No Active Session", http.StatusInternalServerError)
						return
					}
					like.LikeCOunt, err = countLikesForPost(db, like.Post_Id, "like", target)
					if err != nil {
						http.Error(w, "Error Counting like", http.StatusInternalServerError)
						return
					}
					like.DislikeCOunt, err = countLikesForPost(db, like.Post_Id, "dislike", target)
					if err != nil {
						http.Error(w, "Error Counting dislike", http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(&like)
					fmt.Println("heeeeeeere")
					fmt.Println(target)
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

func countLikesForPost(db *sql.DB, postID int, liketype string, target string) (int, error) {
	adding := ""
	if target == "post" {
		adding = "AND comment_id = -1"
	}
	query := fmt.Sprintf(`SELECT COUNT(*) FROM likes WHERE %s_id = ? AND TypeOfLike = ? %s`, target, adding)
	var likeCount int
	err := db.QueryRow(query, postID, liketype).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}
