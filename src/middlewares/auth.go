package middlewares

import (
	db "backend/src/db"
	"backend/src/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.Request.Header.Get("x-auth-key")
		if authKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"auth_error": "Token has not been provided",
			})
			c.Abort()
			return
		}
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authKey, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SUPABASE_JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"auth_error": err.Error()})
			c.Abort()
			return
		}

		// TODO: use claims
		if _, ok := claims["email"].(string); ok {
			var profile models.Profile
			if err := db.Db.Where(&models.Profile{UserId: "0fc76ab9-27f6-402d-865f-280ffa15e068"}). // TODO: change
														First(&profile).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"auth_error": err.Error()})
				c.Abort()
				return
			}
			c.Set("profile", profile)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"auth_error": "Invalid email inside JWT"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthOptionalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.Request.Header.Get("x-auth-key")
		if authKey == "" {
			return
		}
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authKey, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SUPABASE_JWT_SECRET")), nil
		})

		if err != nil {
			return
		}

		// TODO: use claims
		if _, ok := claims["email"].(string); ok {
			var profile models.Profile
			if err := db.Db.Where(&models.Profile{UserId: "0fc76ab9-27f6-402d-865f-280ffa15e068"}). // TODO: change
														First(&profile).Error; err != nil {
				return
			}
			c.Set("profile", profile)
		} else {
			return
		}

		c.Next()
	}
}
