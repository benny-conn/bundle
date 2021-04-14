package internal

type Authorization struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
