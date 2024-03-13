package handler

type LoginRequest struct {
	Hp       string `json:"hp" form:"hp"`
	Password string `json:"password" form:"password" `
}
