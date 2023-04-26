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
			return c.SendStatus(400)
		}

		if err := utils.Validate(*userData); err != nil {
			fmt.Println("error", err)
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		userCollection := database.GetCollection(constants.Mongo["Users"])

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		userData.Password = string(hashedPassword)

		_, err = userCollection.InsertOne(context.TODO(), userData)
		if err != nil {
			fmt.Println("error", err)
			return c.SendStatus(400)
		}
		return c.SendStatus(200)
	}
}

func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// logic for login

		req := new(LoginSchema)

		if err := c.BodyParser(req); err != nil {
			fmt.Println("error", err)
			return c.SendStatus(400)
		}

		if err := utils.Validate(*req); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		userCollection := database.GetCollection(constants.Mongo["Users"])

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		data := userCollection.FindOne(ctx, bson.D{{"email", req.Email}})

		userData := new(LoginSchema)

		if err := data.Decode(userData); err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password))
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": req.Email,
		}).SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		loggedInUsers := database.GetCollection(constants.Mongo["ActiveUsers"])

		alreadyLoggedIn := loggedInUsers.FindOne(ctx, bson.D{{"email", req.Email}})

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
			return c.SendStatus(400)
		}

	sendToken:
		return c.Status(200).JSON(&fiber.Map{
			"message":     "successfully logged in",
			"accessToken": token,
		})
	}
}

func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {

		body := new(LogoutSchema)

		if err := c.BodyParser(body); err != nil {
			fmt.Println("error", err)
			return c.SendStatus(400)
		}

		if err := utils.Validate(*body); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		loggedInUsers := database.GetCollection(constants.Mongo["ActiveUsers"])

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		_, err := loggedInUsers.DeleteMany(ctx, bson.M{"email": body.Email})

		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		return c.SendStatus(200)
	}
}
