package users

type CreateSchema struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoggedIn struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
