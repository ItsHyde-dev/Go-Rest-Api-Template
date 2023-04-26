package users

type AllUsersSchema struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type UpdateUserDetailsSchema struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
