package main

import (
	"github.com/gofiber/fiber/v2"
	"main.go/components/auth"
	"main.go/components/users"
)

func Router(app fiber.Router) {

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendStatus(200)
    })

    app.Post("/", func(c *fiber.Ctx) error {
        return c.SendStatus(200)
    })

	users.UserRoutes(app.Group("user"))
	auth.AuthRoutes(app.Group("auth"))
}
