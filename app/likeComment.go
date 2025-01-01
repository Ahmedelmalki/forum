package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Structure pour gérer les likes/dislikes des commentaires
type CommentReaction struct {
	UserId       int    `json:"userId"`
	CommentId    int    `json:"commentId"`
	Type         string `json:"type"`
	LikeCount    int    `json:"likeCount"`
	DislikeCount int    `json:"dislikeCount"`
}

// Gestion des likes/dislikes des commentaires
func HandleCommentReactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reaction CommentReaction
		var err error

		switch r.Method {
		// 🔄 **Gestion des likes et dislikes (POST)**
		case http.MethodPost:
			err = json.NewDecoder(r.Body).Decode(&reaction)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			
			reaction.UserId, err = ValidateCookie(db, w, r)
				// fmt.Println("User id  now:",like.User_Id)
				if err != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}

			reaction.LikeCount, err = countLikesForComment(db, reaction.CommentId, reaction.Type)
				if err != nil {
					http.Error(w, "Error counting likes", http.StatusInternalServerError)
					return
				}

				checkQuery := "SELECT EXISTS(SELECT 1 FROM comment_reactions WHERE comment_id = ? AND user_id = ?)"
				var exists bool
				err = db.QueryRow(checkQuery, reaction.CommentId, reaction.UserId).Scan(&exists)
				if err != nil {
					http.Error(w, "Error checking like existence", http.StatusInternalServerError)
					return
				}


				if exists {
					LiketypeQuery := "SELECT TypeOfLike FROM comment_reactions WHERE comment_id = ? AND user_id = ?"
					var typea string
					db.QueryRow(LiketypeQuery, reaction.CommentId, reaction.UserId).Scan(&typea)
					if typea == reaction.Type {
						query := "DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?"
						_, err = db.Exec(query, reaction.CommentId, reaction.UserId)
						if err != nil {
							http.Error(w, "Error deleting like", http.StatusInternalServerError)
							return
						}
					} else {
						Updatequery := "UPDATE comment_reactions SET TypeOfLike = ? WHERE comment_id = ? AND user_id = ?"
						_, err = db.Exec(Updatequery, reaction.Type, reaction.CommentId, reaction.UserId)
						if err != nil {
							http.Error(w, "Error UPDATNG comment_reactions", http.StatusInternalServerError)
							return
						}
					}
				} else {
					query := "INSERT INTO comment_reactions (user_id, comment_id , TypeOfLike) VALUES (?, ?, ?)"
					_, err = db.Exec(query, reaction.UserId, reaction.CommentId, reaction.Type)
					if err != nil {
						http.Error(w, "Error adding like", http.StatusInternalServerError)
						return
					}
				}


		// 📊 **Récupération des likes et dislikes (GET)**
	case http.MethodGet:
		{
			if reaction.CommentId != 0 {
				reaction.UserId, err = ValidateCookie(db, w, r)
				if err != nil {
					http.Error(w, "No Active Session", http.StatusInternalServerError)
					return
				}
				reaction.LikeCount, err = countLikesForComment(db, reaction.CommentId, "like")
				if err != nil {
					http.Error(w, "Error Counting like", http.StatusInternalServerError)
					return
				}
				reaction.DislikeCount, err = countLikesForComment(db, reaction.CommentId, "dislike")
				if err != nil {
					http.Error(w, "Error Counting dislike", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(&reaction)
				fmt.Println("heeeeeeere")
				fmt.Println(reaction)
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

func countLikesForComment(db *sql.DB, commentID int, liketype string) (int, error) {
	query := "SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND TypeOfLike = ? "
	var likeCount int
	err := db.QueryRow(query, commentID, liketype).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}