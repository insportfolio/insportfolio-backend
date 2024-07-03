package controllers

import (
	db "backend/src/db"
	"backend/src/models"
	"backend/src/schemas"
	"backend/src/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags Portfolios
// @Produce  json
// @Title Get Portfolio
// @Summary Get Portfolio by Profile id
// @Description Get portfolio by Profile id, the views counter is incremented only if the user is logged and is not the owner
// @Param profile_id path string true "profile_id"
// @Router /portfolios/{profile_id} [get]
// @Success 200 {object} schemas.GetPortfolioResponse
// @Security ApiKeyAuth
func GetPortfolio(c *gin.Context) {
	profileIdStr := c.Params.ByName("profile_id")
	profileId, err := strconv.Atoi(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile_id must be an integer"})
		return
	}

	var portfolio models.Portfolio
	if err := db.Db.Preload("Technologies").Where(&models.Portfolio{ProfileId: profileId, Active: true}).
		First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	profile, ok := c.Get("profile")
	if ok {
		parsedProfile := profile.(models.Profile)
		if uint(portfolio.ProfileId) != parsedProfile.ID {
			db.Db.Model(&models.Portfolio{}).Where("id = ?", portfolio.ID).Update("views", portfolio.Views+1)
		}
	}

	portfolio.SetFullPreview()

	var reactionsCount int64
	var commentsCount int64

	if err := db.Db.Model(&models.PortfolioReaction{}).
		Where("portfolio_id = ?", portfolio.ID).Count(&reactionsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching reactions count"})
		return
	}
	if err := db.Db.Model(&models.Comment{}).
		Where("portfolio_id = ?", portfolio.ID).Count(&commentsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching comments count"})
		return
	}
	c.JSON(http.StatusOK, schemas.GetPortfolioResponse{
		DefaultSchema: schemas.DefaultSchema{
			ID:        portfolio.ID,
			CreatedAt: portfolio.CreatedAt,
			UpdatedAt: portfolio.UpdatedAt,
			DeletedAt: portfolio.DeletedAt.Time,
		},
		Url:            portfolio.Url,
		Preview:        portfolio.Preview,
		ProfileId:      portfolio.ProfileId,
		Active:         portfolio.Active,
		Views:          portfolio.Views,
		ReactionsCount: reactionsCount,
		CommentsCount:  commentsCount,
		Technologies:   models.ConvertTechnologiesToStructArray(portfolio.Technologies),
	})
}

// @Tags Portfolios
// @Accept  mpfd
// @Produce  json
// @Title Create Portfolio
// @Summary create portfolio and upload file to AWS
// @Description The portfolio is added to the list of portfolios of the profile and the others are setted as active: false
// @Param url formData string true "Portfolio URL"
// @Param preview formData file true "Portfolio preview image"
// @Param technologies formData array true "Technologies used"
// @Router /portfolios/ [post]
// @Success 200 {object} schemas.Portfolio
// @Security ApiKeyAuth
func CreatePortfilio(c *gin.Context) {
	profile, _ := c.Get("profile")
	parsedProfile := profile.(models.Profile)

	var portfolio models.Portfolio
	var technologies []models.Technology

	portfolioUrl := c.PostForm("url")
	portfolioTechnologies := c.PostForm("technologies")
	portfolioTechnologiesParsed, err := utils.StringToArray(portfolioTechnologies)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid technologies array",
		})
		return
	}

	// Retrieve the file
	file, err := c.FormFile("preview")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	name, err := utils.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	portfolio.ProfileId = int(parsedProfile.ID)
	portfolio.Url = portfolioUrl
	portfolio.Preview = name

	// Deactivate all the user portfolios
	if err := db.Db.Model(&models.Portfolio{}).Where("active = ?", true).Update("active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrive al the requested technologies
	if err := db.Db.Where("id IN ?", portfolioTechnologiesParsed).Find(&technologies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Append the fetched Technology records to the Portfolio
	if err := db.Db.Model(&portfolio).Association("Technologies").Append(&technologies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	portfolio.SetFullPreview()

	c.JSON(http.StatusOK, portfolio)
}

// @Tags Portfolios
// @Accept  mpfd
// @Produce  json
// @Title Udate Portfolio
// @Summary Update profile activated portfolio
// @Param url formData string true "Portfolio URL"
// @Param preview formData file true "Portfolio preview image"
// @Param technologies formData array true "Technologies used"
// @Router /portfolios/ [put]
// @Success 200 {object} schemas.Portfolio
// @Security ApiKeyAuth
func UpdatePortfolio(c *gin.Context) {
	profile, _ := c.Get("profile")
	parsedProfile := profile.(models.Profile)

	var portfolio models.Portfolio
	var technologies []models.Technology

	portfolioUrl := c.PostForm("url")
	portfolioTechnologies := c.PostForm("technologies")
	portfolioTechnologiesParsed, err := utils.StringToArray(portfolioTechnologies)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid technologies array",
		})
	}

	if err := db.Db.Where(models.Portfolio{
		Active:    true,
		ProfileId: int(parsedProfile.ID),
	}).First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	// Retrieve the file
	file, err := c.FormFile("preview")
	if err == nil {
		name, err := utils.UploadFile(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		portfolio.Preview = name
	}

	portfolio.Url = portfolioUrl
	// Retrive al the requested technologies
	if err := db.Db.Where("id IN ?", portfolioTechnologiesParsed).Find(&technologies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Append the fetched Technology records to the Portfolio
	if err := db.Db.Model(&portfolio).Association("Technologies").Replace(&technologies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	portfolio.SetFullPreview()

	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}
