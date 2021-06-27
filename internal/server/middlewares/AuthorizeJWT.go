package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/morticius/accuracy/internal/server/services"
	"log"
	"net/http"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		jwtService := services.NewJWTAuthService()
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || token == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println(claims)
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
