package main

import (
	"backend/src/controllers"
	db "backend/src/db"
	_ "backend/src/docs"
	"backend/src/middlewares"
	"backend/src/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Insportfolio Docs
// @Version 0.0.1
// @description Insportfolio Api Docs
// @BasePath /api/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-auth-key
func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Cannot load .env")
	}

	db.Initialize()
	db.Db.AutoMigrate(
		&models.Profile{},
		&models.Portfolio{},
		&models.Technology{},
		&models.PortfolioReaction{},
		&models.Comment{},
		&models.CommentLike{},
	)

	r := gin.Default()

	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		api.POST("/auth/profile/new", controllers.CreateProfile)

		api.GET("/users/profile", middlewares.AuthMiddleware(), controllers.GetProfile)

		api.GET("/portfolios/:profile_id", middlewares.AuthOptionalMiddleware(), controllers.GetPortfolio)
		api.POST("/portfolios", middlewares.AuthMiddleware(), controllers.CreatePortfilio)
		api.PUT("/portfolios", middlewares.AuthMiddleware(), controllers.UpdatePortfolio)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))

	r.Run("127.0.0.1:8080")
}
