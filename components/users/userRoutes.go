package users

import (
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	app.Get("/", GetAllUsers())

	app.Post("/create", Create())

	app.Post("/login", Login())

	app.Post("/logout", Logout())

}
