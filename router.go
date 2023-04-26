package main

import (
	"github.com/gofiber/fiber/v2"
	"main.go/components/auth"
	"main.go/components/users"
)

func Router(app fiber.Router) {
	users.UserRoutes(app.Group("user"))
	auth.AuthRoutes(app.Group("auth"))
}
