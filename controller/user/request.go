package user

type LoginRequest struct {
	Hp       string `json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterRequest struct {
	Nama     string `json:"nama" form:"nama" validate:"required"`
	Hp       string `json:"hp" form:"hp" validate:"required,max=13,min=10"`
	Password string `json:"password" form:"password" validate:"required,min=10,uppercase"`
}
