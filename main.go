package main

import (
	"backend/controllers"
	"backend/initializers"
	"backend/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main(){
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	// r.POST("/signup", controllers.Signup)
	r.POST("/login",middleware.ValidateAuth, controllers.LoginWithCCS)
	r.GET("/validate", middleware.RequireAuth,controllers.Validate)
	r.Run(fmt.Sprintf(":%s", "8080"))
}