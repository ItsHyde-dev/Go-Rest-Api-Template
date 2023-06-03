package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/database"
)

var MongoClient mongo.Client

func main() {
	app := fiber.New()
    app.Use(cors.New())
    app.Use(limiter.New())
    app.Use(logger.New())

	godotenv.Load()

	database.ConnectToDatabase()

	publicRouter := app.Group("/")

	Router(publicRouter)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	port := os.Getenv("PORT")

    fmt.Println("listening on " + port)

	if port == "" {
        port = ":8080"
	} else {
        port = ":" + port
    }

	log.Fatal(app.Listen(port))
}
