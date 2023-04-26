package users

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"main.go/constants"
	"main.go/database"
)

func GetAllUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCollection := database.GetCollection(constants.MONGO_COLLECTIONS["USERS"])

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		data, err := userCollection.Find(ctx, bson.D{})
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		var users []CreateSchema

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

func Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData := new(CreateSchema)

		if err := c.BodyParser(userData); err != nil {
			return c.SendStatus(400)
		}

		if err := validate(*userData); err != nil {
			fmt.Println("error", err)
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		userCollection := database.GetCollection(constants.MONGO_COLLECTIONS["USERS"])

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

		req := new(CreateSchema)

		if err := c.BodyParser(req); err != nil {
			fmt.Println("error", err)
			return c.SendStatus(400)
		}

		if err := validate(*req); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
			})
		}

		userCollection := database.GetCollection(constants.MONGO_COLLECTIONS["USERS"])

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		data := userCollection.FindOne(ctx, bson.D{{"email", req.Email}})

		userData := new(CreateSchema)

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

		loggedInUsers := database.GetCollection(constants.MONGO_COLLECTIONS["LOGGED_IN_USERS"])

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

		// we get the accessToken and we invalidate it

		return c.SendStatus(200)
	}
}

func validate[T any](request T) *string {
	err := validator.New().Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			response := "Please provide " + err.Field()
			return &response
		}
	}
	return nil
}
