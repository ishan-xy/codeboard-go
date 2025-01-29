package main

import (
	"backend/controllers"
	"backend/initializers"
	"backend/middleware"
	"fmt"
	"log"

	utils "github.com/ItsMeSamey/go_utils"
	"github.com/gofiber/fiber/v2"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	app := fiber.New()

	// Add middleware for parsing JSON
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		return c.Next()
	})

	utils.SetErrorStackTrace(true)

	// Authentication routes
	app.Post("/login", middleware.ValidateAuth, controllers.LoginWithCCS)
	app.Get("/validate", middleware.RequireAuth, controllers.Validate)
	app.Post("/leetcode", middleware.RequireAuth, controllers.VerifyLeetcodeUsername)
	app.Post("/confirm-leetcode", middleware.RequireAuth, controllers.ConfirmLeetcode)
	app.Get("/profile", middleware.RequireAuth, controllers.GetUserProfile)

	// Start the server
	err := app.Listen(fmt.Sprintf(":%s", "8080"))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
