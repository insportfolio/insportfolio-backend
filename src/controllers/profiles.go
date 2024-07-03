package controllers

import (
	db "backend/src/db"
	_ "backend/src/docs"
	"backend/src/models"
	"backend/src/schemas"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	supabaseHeader := os.Getenv("SUPABASE_AUTH_HOOK_KEY")
	apiKey := c.Request.Header.Get("x-api-key")
	if supabaseHeader != apiKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth"})
		return
	}

	profileRequest := schemas.NewProfileRequest{}
	if err := c.ShouldBindJSON(&profileRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid data"})
		return
	}

	profile := models.Profile{
		Email:  profileRequest.Record.Email,
		UserId: profileRequest.Record.ID,
	}
	if err := db.Db.Create(&profile).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Insertion failed"})
		return
	}
}

// @Tags Users
// @Accept  json
// @Produce  json
// @Title Get Profile
// @Description Get a string by ID
// @Summary Return authenticated user profile
// @Success 200 {object} schemas.ProfileSchema
// @Router /users/profile [get]
// @Security ApiKeyAuth
func GetProfile(c *gin.Context) {
	if val, ok := c.Get("profile"); ok {
		c.JSON(200, gin.H{"profile": val})
	}
}
