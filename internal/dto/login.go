package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	MemberData Member `json:"member"`
}

type Member struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
