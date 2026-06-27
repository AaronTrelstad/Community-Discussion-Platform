package requests

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Eamil    string `json:"email"`
	Password string `json:"password"`
}
