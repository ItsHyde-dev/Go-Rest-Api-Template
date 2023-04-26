package utils

import "github.com/go-playground/validator"

func Validate[T any](request T) *string {
	err := validator.New().Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			response := "Please provide " + err.Field()
			return &response
		}
	}
	return nil
}
