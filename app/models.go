package forum

import "time"

type RegisterCredenials struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	ID         int      `json:"Id"`
	UserName   string   `json:"Username"`
	Title      string   `json:"Title"`
	Content    string   `json:"Content"`
	Categories []string `json:"Categories"`
	CreatedAt  string   `json:"Created_at"`
	Likes      int      `json:"Likes"`
	Dislikes   int      `json:"Dislikes"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserName  string    `json:"username"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"Likes"`
	Dislikes  int       `json:"Dislikes"`
}


