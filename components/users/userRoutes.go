package users

import (
	"github.com/gofiber/fiber/v2"
	"main.go/components/auth"
)

func UserRoutes(app fiber.Router) {
	app.Get("/", GetAllUsers())

	app.Post("/details/update", auth.ValidateToken(), UpdateUserDetails())

}
