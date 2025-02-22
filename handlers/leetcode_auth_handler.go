package handlers

import (
	"backend/database"
	"backend/utility"
	"log"

	utils "github.com/ItsMeSamey/go_utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func VerifyLeetcodeUsername(c fiber.Ctx) error {

	var req database.Username
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	_, exists, err := database.PublicUsersDB.GetExists(bson.M{"leetcode_username": req.Username})
	
	if err != nil {
		log.Printf("Error checking username: %v", utils.WithStack(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify LeetCode username",
		})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "LeetCode username already exists",
		})
	}

	variables := map[string]interface{}{
		"username": req.Username,
	}

	profile, err := utility.SendQuery(utility.PROFILE_QUERY, variables)
	if err != nil {
		log.Printf("Error fetching profile: %v", utils.WithStack(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch LeetCode profile",
		})
	}

	numQuestions, err := utility.SendQuery(utility.NUM_QUESTION_QUERY, variables)
	if err != nil {
		log.Printf("Error fetching questions: %v", utils.WithStack(err))
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
	cacheData := database.LeetCodeCacheData{
		RealName:       profileData["realName"].(string),
		UserAvatar:     profileData["userAvatar"].(string),
		QuestionCount:  int(acSubmissions["count"].(float64)),
		Ranking:        int(profileData["ranking"].(float64)),
		TotalQuestions: int(allQuestions["count"].(float64)),
	}

	database.UserCache.PublicUsers.Set(req.Username, cacheData)

	return c.JSON(fiber.Map{
		"message":  "LeetCode profile verified",
		"username": req.Username,
	})
}

func ConfirmLeetcode(c fiber.Ctx) error {

	var req database.Username
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	cachedData, found := database.UserCache.PublicUsers.Get(req.Username)
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Verification data expired or not found",
		})
	}

	user, ok := c.Locals("user").(database.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Create public profile
	id := primitive.NewObjectID()
	publicProfile := database.UserPublicProfile{
		ID:               id,
		ProfilePic:       cachedData.UserAvatar,
		LeetcodeUsername: req.Username,
		QuestionCount:    cachedData.QuestionCount,
		Ranking:          cachedData.Ranking,
		TotalQuestions:   cachedData.TotalQuestions,
	}
	_, err := database.PublicUsersDB.InsertOne(c.Context(), publicProfile)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user: " + utils.WithStack(err).Error(),
		})
	}

	user.Name = cachedData.RealName
	user.PublicProfileID = id
	update := bson.M{"$set": bson.M{"name": user.Name, "public_profile": user.PublicProfileID}}
	_, err = database.UserDB.UpdateByID(c.Context(), user.ID, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user: " + utils.WithStack(err).Error(),
		})
	}

	// Clear cache entry
	database.UserCache.PublicUsers.Remove(req.Username)

	return c.JSON(fiber.Map{
		"message": "LeetCode profile connected successfully",
		"profile": publicProfile,
	})
}
