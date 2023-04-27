package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/constants"
	"main.go/database"
	"main.go/utils"
)

func ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {

		throw := func(err error) error {
			fmt.Println("error", err)
			return c.Status(403).JSON(&fiber.Map{
				"message": "Invalid access token",
			})
		}

		loggedInUsers := database.GetCollection(constants.ActiveUserCollection)
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		headers := new(AuthHeaderSchema)

		if err := c.ReqHeaderParser(headers); err != nil {
			return throw(err)
		}

		if err := utils.Validate(*headers); err != nil {
			return c.Status(403).JSON(&fiber.Map{
				"message": err,
			})
		}

		token := headers.AccessToken

		if doc := loggedInUsers.FindOne(ctx, bson.M{"token": token}); doc.Err() == mongo.ErrNoDocuments {
			return throw(doc.Err())
		}

		claims := jwt.MapClaims{}

		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return throw(err)
		}

		for key, value := range claims {
			c.Locals(key, value)
		}

		return c.Next()

	}
}
