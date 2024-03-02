package user

type LoginResponse struct {
	Hp    string `json:"hp"`
	Nama  string `json:"nama"`
	Token string `json:"token"`
}
