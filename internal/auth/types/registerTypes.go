package types

type UserRegisterRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"Password"`
}
