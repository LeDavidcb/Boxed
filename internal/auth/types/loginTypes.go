package types

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	SignedJwt    string `json:"signed-jwt"`
	RefreshToken string `json:"refresh-token"`
}
