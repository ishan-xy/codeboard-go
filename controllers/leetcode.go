package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/utility"
	"log"

	"github.com/gofiber/fiber/v2"
)

func VerifyLeetcodeUsername(c *fiber.Ctx) error {
	type username struct {
		Username string `json:"username"`
	}

	var usernameData username
	if err := c.BodyParser(&usernameData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	variables := map[string]interface{}{
		"username": usernameData.Username,
	}

	profile, err := utility.SendQuery(utility.PROFILE_QUERY, variables)
	if err != nil {
		log.Printf("Error fetching profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch LeetCode profile",
		})
	}

	numQuestions, err := utility.SendQuery(utility.NUM_QUESTION_QUERY, variables)
	if err != nil {
		log.Printf("Error fetching questions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch LeetCode questions",
		})
	}

	if len(profile.Errors) > 0 || len(numQuestions.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid LeetCode username",
		})
	}

	profileData := profile.Data["matchedUser"].(map[string]interface{})["profile"].(map[string]interface{})
	submitStats := numQuestions.Data["matchedUser"].(map[string]interface{})["submitStats"].(map[string]interface{})
	acSubmissions := submitStats["acSubmissionNum"].([]interface{})[0].(map[string]interface{})
	allQuestions := numQuestions.Data["allQuestionsCount"].([]interface{})[0].(map[string]interface{})

	// Store both user info and stats in cache
	cacheData := models.LeetCodeCacheData{
		RealName:       profileData["realName"].(string),
		UserAvatar:     profileData["userAvatar"].(string),
		QuestionCount:  int(acSubmissions["count"].(float64)),
		Ranking:        int(profileData["ranking"].(float64)),
		TotalQuestions: int(allQuestions["count"].(float64)),
	}

	models.UserCache.PublicUsers.Set(usernameData.Username, cacheData)

	return c.JSON(fiber.Map{
		"message": "LeetCode profile verified",
		"username": usernameData.Username,
	})
}

func ConfirmLeetcode(c *fiber.Ctx) error {
	type username struct {
		Username string `json:"username"`
	}

	var usernameData username
	if err := c.BodyParser(&usernameData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	cachedData, found := models.UserCache.PublicUsers.Get(usernameData.Username)
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Verification data expired or not found",
		})
	}

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Update user fields
	user.LeetcodeUsername = usernameData.Username
	user.Name = cachedData.RealName
	user.ProfilePic = cachedData.UserAvatar

	// Save user updates
	if err := initializers.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Failed to update user: " + err.Error(),
		})
	}

	// Create public profile
	publicProfile := models.UserPublicProfile{
		LeetcodeUsername: user.LeetcodeUsername,
		QuestionCount:    cachedData.QuestionCount,
		Ranking:          cachedData.Ranking,
		TotalQuestions:   cachedData.TotalQuestions,
	}

	if err := initializers.DB.Create(&publicProfile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create public profile: " + err.Error(),
		})
	}

	// Clear cache entry
	models.UserCache.PublicUsers.Remove(usernameData.Username)

	return c.JSON(fiber.Map{
		"message": "LeetCode profile connected successfully",
		"profile": publicProfile,
	})
}