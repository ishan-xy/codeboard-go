package router

import (
	"backend/middleware"
	"backend/handlers"

	"github.com/gofiber/fiber/v3"
)

func addAuthRoutes(r fiber.Router){
	r.Post("/login", handlers.LoginWithCCS, middleware.ValidateAuth)
	r.Post("/leetcode", handlers.VerifyLeetcodeUsername, middleware.RequireAuth)
	r.Post("/confirm-leetcode", handlers.ConfirmLeetcode, middleware.RequireAuth)
	// r.Get("/profile", middleware.RequireAuth, handlers.GetUserProfile)
}