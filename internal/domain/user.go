package domain

// simple model for user authorization
type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt int
	UpdatedAt int
}

type UserAuthReq struct {
	Email    string
	Password string
}

type UserRegisterReq struct {
	Email      string
	Password   string
	RePassword string
}
