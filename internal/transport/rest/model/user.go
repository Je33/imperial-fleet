package model

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserAuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type UserRefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type UserAuthRes struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}
