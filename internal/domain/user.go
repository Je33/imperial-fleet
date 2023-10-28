package domain

// simple model for user authorization
type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

// model of authorisation request
type UserAuthReq struct {
	Email    string
	Password string
}

// model of registration request
type UserRegisterReq struct {
	Email      string
	Password   string
	RePassword string
}
