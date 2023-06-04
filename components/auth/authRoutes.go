package auth

import "github.com/gofiber/fiber/v2"

func AuthRoutes(app fiber.Router) {

	app.Post("/signup", Signup())

	app.Post("/login", Login())

	app.Post("/logout", ValidateToken(), Logout())
}
