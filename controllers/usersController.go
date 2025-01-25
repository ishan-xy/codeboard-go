package controllers

import (
	"backend/initializers"
	"backend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// func Signup(c *gin.Context){

// 	// bind request body to struct
// 	var body struct {
// 		Email string
// 		Password string
// 	}
// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid parameters",
// 		})
// 		return
// 	}

// 	// hash password
// 	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to hash password",
// 		})
// 		return
// 	}
// 	user := models.User{Email: body.Email, Password: string(hash)}
// 	result := initializers.DB.Create(&user)

// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to create user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User created",})
// }

// func Login(c *gin.Context){

// 	// get email and password from request
// 	var body struct {
// 		Email string
// 		Password string
// 	}
// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid parameters",
// 		})
// 		return
// 	}

// 	// find user by email
// 	var user models.User
// 	result := initializers.DB.First(&user, "email = ?", body.Email)

// 	if result.Error != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid email or password",
// 		})
// 		return
// 	}
	
// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid email or password",
// 		})
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": user.ID,
// 		"exp": time.Now().Add(time.Hour * 24).Unix(),
// 	})

// 	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to generate token",
// 		})
// 		return
// 	}
// 	c.SetSameSite(http.SameSiteLaxMode)
// 	c.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true)

// 	c.JSON(http.StatusOK, gin.H{"message": "Logged in",})
// }

func LoginWithCCS(c *gin.Context){

	var body struct {
		Email string
		Name string
		ProfilePic string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parameters",
		})
		return
	}

	// Signup user if not exists
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)

	if result.Error != nil {
		user := models.User{Email: body.Email, Name: body.Name, ProfilePic: body.ProfilePic}
		result := initializers.DB.Create(&user)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created",})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged in",})
}

func Validate(c *gin.Context){

	c.JSON(http.StatusOK, gin.H{"message": "Valid token",})
}