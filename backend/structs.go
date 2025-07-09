package backend

type User struct {
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Gender     string `json:"gender"`
	Identifier string `json:"identifier"`
}

type LoginInput struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Note struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
