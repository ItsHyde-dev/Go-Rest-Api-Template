package users

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"main.go/constants"
	"main.go/database"
	"main.go/utils"
)

func GetAllUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCollection := database.GetCollection(constants.UserCollection)

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		data, err := userCollection.Find(ctx, bson.D{})
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		var users []AllUsersSchema

		if err = data.All(ctx, &users); err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		return c.Status(200).JSON(&fiber.Map{
			"message": "Successfully fetched user",
			"data":    users,
		})
	}
}

func UpdateUserDetails() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCollection := database.GetCollection(constants.UserCollection)

		body := new(UpdateUserDetailsSchema)
		if err := c.BodyParser(body); err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		if err := utils.Validate(*body); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		email := c.Locals("email")

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		_, err := userCollection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": *body})

		if err != nil {
			fmt.Println("error", err)
			return c.Status(400).JSON(&fiber.Map{
				"message": "failed to update user details",
			})
		}

		return c.SendStatus(200)
	}
}
