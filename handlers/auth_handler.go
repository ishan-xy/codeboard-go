package handlers

import (
	"backend/config"
	"backend/database"
	"log"
	"time"

	utils "github.com/ItsMeSamey/go_utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func LoginWithCCS(c fiber.Ctx) error {
	// Get decrypted data from middleware
	decryptedData, ok := c.Locals("data").(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user data",
		})
	}

	// Extract user details
	email, _ := decryptedData["email"].(string)
	name, _ := decryptedData["name"].(string)

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing email",
		})
	}

	result, exists, err := database.UserDB.GetExists(bson.M{"email": email})
	if err != nil {
		log.Println(utils.WithStack(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": utils.WithStack(err),
		})
	}
	var user database.User
	if !exists {
		user = database.User{
			Email:      email,
			Name:       name,
			ID: 	   primitive.NewObjectID(),
		}
		_, err = database.UserDB.InsertOne(c.Context(), user)
		if err != nil {
			log.Println(utils.WithStack(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": utils.WithStack(err),
			})
		}
	}else{
		user = result
	}

	// Generate token
	log.Println(user.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Hex(),
		"email": user.Email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Getenv("SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600 * 24,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged in",
		"user": fiber.Map{
			"id":         user.ID.Hex(),
			"email":      user.Email,
			"name":       user.Name,
		},
	})
}