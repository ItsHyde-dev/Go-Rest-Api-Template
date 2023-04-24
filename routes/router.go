package routes

import "github.com/gofiber/fiber/v2"

func Router(app fiber.Router) {

	UserRoutes(app.Group("/user"))
}
