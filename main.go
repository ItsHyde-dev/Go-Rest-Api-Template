package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/routes"
)

func main() {
	app := fiber.New()

	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://jwksdbuser:jwk$prod^123@10.166.67.5/jwksprod?directConnection=true&authSource=admin&readPreference=secondary"))

	if err != nil {
		fmt.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = db.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Database connected")

	defer db.Disconnect(ctx)
	defer fmt.Println("Database disconnecting...")

	publicRouter := app.Group("/")

	routes.Router(publicRouter)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Listen(":3000")
}
