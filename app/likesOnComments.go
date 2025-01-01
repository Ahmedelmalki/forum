package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type likesOnCmnts struct {
	User_Id      int    `json:"UserId"`
	Comment_Id   int    `json:"CommentId"`
	LikeCount    int    `json:"LikeCount"`
	DisLikeCount int    `json:"DislikeCount"`
	Type         string `json:"Type"`
}

/*
CREATE TABLE if NOT EXISTS likesOnComment(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    TypeOfLike TEXT not NULL,
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
    FOREIGN KEY (comment_id) REFERENCES comments (id)
);
*/

func HandleLikesOnComments(db *sql.DB) http.HandlerFunc {
	fmt.Println("HandleLikesOnComments called")
	var likeOnComment likesOnCmnts

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.Method {
		case http.MethodPost:
			{
				err := json.NewDecoder(r.Body).Decode(&likeOnComment)
				if err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				likeOnComment.User_Id, err = ValidateCookie(db, w, r)
				// fmt.Println("User id  now:",like.User_Id)
				if err != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
				likeOnComment.LikeCount, err = countLikesForComment(db, likeOnComment.Comment_Id, likeOnComment.Type)
				if err != nil {
					http.Error(w, "Error counting likes", http.StatusInternalServerError)
					return
				}
				checkQuery := "SELECT EXISTS(SELECT 1 FROM likesOnComment WHERE comment_id = ? AND user_id = ?)"
				var exists bool
				err = db.QueryRow(checkQuery, likeOnComment.Comment_Id, likeOnComment.User_Id).Scan(&exists)
				if err != nil {
					http.Error(w, "Error checking like existence", http.StatusInternalServerError)
					return
				}

				if exists {
					LiketypeQuery := "SELECT TypeOfLike FROM likesOnComment WHERE comment_id = ? AND user_id = ?"
					var typea string
					db.QueryRow(LiketypeQuery, likeOnComment.Comment_Id, likeOnComment.User_Id).Scan(&typea)
					if typea == likeOnComment.Type {
						query := "DELETE FROM likesOnComment WHERE comment_id = ? AND user_id = ?"
						_, err = db.Exec(query, likeOnComment.Comment_Id, likeOnComment.User_Id)
						if err != nil {
							http.Error(w, "Error deleting like", http.StatusInternalServerError)
							return
						}
					} else {
						Updatequery := "UPDATE likesOnComment SET TypeOfLike = ? WHERE comment_id = ? AND user_id = ?"
						_, err = db.Exec(Updatequery, likeOnComment.Type, likeOnComment.Comment_Id, likeOnComment.User_Id)
						if err != nil {
							http.Error(w, "Error UPDATNG likeS", http.StatusInternalServerError)
							return
						}
					}
				} else {
					query := "INSERT INTO likesOnComment (user_id, comment_id , TypeOfLike) VALUES (?, ?, ?)"
					_, err = db.Exec(query, likeOnComment.User_Id, likeOnComment.Comment_Id, likeOnComment.Type)
					if err != nil {
						http.Error(w, "Error adding like", http.StatusInternalServerError)
						return
					}
				}

			}
		case http.MethodGet:
			{
				if likeOnComment.Comment_Id != 0 {
					likeOnComment.User_Id, err = ValidateCookie(db, w, r)
					if err != nil {
						http.Error(w, "No Active Session", http.StatusInternalServerError)
						return
					}
					likeOnComment.LikeCount, err = countLikesForComment(db, likeOnComment.Comment_Id, "like")
					if err != nil {
						http.Error(w, "Error Counting like", http.StatusInternalServerError)
						return
					}
					likeOnComment.DisLikeCount, err = countLikesForComment(db, likeOnComment.Comment_Id, "dislike")
					if err != nil {
						http.Error(w, "Error Counting dislike", http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(&likeOnComment)
					fmt.Println(likeOnComment)
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
	fmt.Println("countLikesForComment called")
	query := "SELECT COUNT(*) FROM likesOnComment WHERE comment_id = ? AND TypeOfLike = ? "
	var LikeCount int
	err := db.QueryRow(query, commentID, liketype).Scan(&LikeCount)
	if err != nil {
		return 0, err
	}

	return LikeCount, nil
}
