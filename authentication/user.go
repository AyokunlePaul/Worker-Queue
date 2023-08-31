package authentication

type User struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Balance  int64  `json:"balance"`
	Verified bool   `json:"verified"`
}
