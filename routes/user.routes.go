package routes

import "github.com/gofiber/fiber/v2"

func UserRoutes(app fiber.Router) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"message": "Hello, World!",
			"status":  "success",
			"route":   "user route",
		})
	})

	app.Get("/user", func(c *fiber.Ctx) error {
		c.Body()
	})
}
