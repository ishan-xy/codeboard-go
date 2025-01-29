package controllers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"backend/initializers"
	"backend/models"
)

func LoginWithCCS(c *fiber.Ctx) error {
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
	profilePic, _ := decryptedData["profilePic"].(string)

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing email",
		})
	}

	// Find or create user
	var user models.User
	result := initializers.DB.First(&user, "email = ?", email)

	if result.Error != nil {
		// User not found, create new user
		user = models.User{
			Email:      email,
			Name:       name,
			ProfilePic: profilePic,
		}
		result = initializers.DB.Create(&user)

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	} else {
		// Update existing user if details changed
		user.Name = name
		user.ProfilePic = profilePic
		initializers.DB.Save(&user)
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
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
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"profilePic": user.ProfilePic,
		},
	})
}

func Validate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Validated",
	})
}

func GetUserProfile(c *fiber.Ctx) error {
	// Get the user from context (basic info)
	authUser := c.Locals("user").(models.User)
	
	// Create a fresh user object with preloaded relations
	var user models.User
	err := initializers.DB.
		Preload("PublicProfile").
		Where("id = ?", authUser.ID).
		First(&user).
		Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user profile",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}