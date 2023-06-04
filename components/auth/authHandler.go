package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"main.go/constants"
	"main.go/database"
	"main.go/utils"
)

func Signup() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData := new(SignupSchema)

		if err := c.BodyParser(userData); err != nil {
			return c.SendStatus(500)
		}

		if err := utils.Validate(*userData); err != nil {
			fmt.Println("error", err)
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		userCollection := database.GetCollection(constants.UserCollection)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		userData.Password = string(hashedPassword)

		_, err = userCollection.InsertOne(context.TODO(), userData)
		if err != nil {
			fmt.Println("error", err)
			return c.SendStatus(500)
		}

		return c.Status(200).JSON(&fiber.Map{
            "message": "successfully signed up",
		})
	}
}

func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {

		req := new(LoginSchema)

		if err := c.BodyParser(req); err != nil {
			fmt.Println("error", err)
			return c.SendStatus(500)
		}

		if err := utils.Validate(*req); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		userCollection := database.GetCollection(constants.UserCollection)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

		data := userCollection.FindOne(ctx, bson.M{
            "email": req.Email,
        })

        if data.Err() != nil {
            fmt.Println(data.Err())
            return c.Status(400).JSON(&fiber.Map{
                "message": "invalid credentials",
            })
        }

		userData := new(LoginSchema)

		if err := data.Decode(userData); err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password))
		if err != nil {
			fmt.Println(err)
			return c.Status(400).JSON(&fiber.Map{
                "message": "invalid credentials",
			})
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": req.Email,
		}).SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		loggedInUsers := database.GetCollection(constants.ActiveUserCollection)

		alreadyLoggedIn := loggedInUsers.FindOne(ctx, bson.M{
            "email": req.Email,
        })

		if alreadyLoggedIn.Err() != mongo.ErrNoDocuments {
			loggedInDoc := new(LoggedIn)
			alreadyLoggedIn.Decode(&loggedInDoc)
			token = loggedInDoc.Token
			goto sendToken
		}

		_, err = loggedInUsers.InsertOne(ctx, bson.M{
			"email": req.Email,
			"token": token,
		})

		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

	sendToken:
		return c.Status(200).JSON(&fiber.Map{
			"message":     "successfully logged in",
			"data": &fiber.Map{
                "accessToken": token,
            },
		})
	}
}

func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {

		body := new(LogoutSchema)

		if err := c.BodyParser(body); err != nil {
			fmt.Println("error", err)
			return c.SendStatus(500)
		}

		if err := utils.Validate(*body); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		loggedInUsers := database.GetCollection(constants.ActiveUserCollection)

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		_, err := loggedInUsers.DeleteMany(ctx, bson.M{"email": body.Email})

		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	}
}
