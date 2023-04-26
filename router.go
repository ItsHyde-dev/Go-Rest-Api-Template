package main

import (
	"github.com/gofiber/fiber/v2"
	"main.go/components/users"
)

func Router(app fiber.Router) {
	users.UserRoutes(app.Group("user"))
}
