package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/database"
)

var MongoClient mongo.Client

func main() {
	app := fiber.New()

	godotenv.Load()

	database.ConnectToDatabase()

	publicRouter := app.Group("/")

	Router(publicRouter)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Listen(":8080")
}
