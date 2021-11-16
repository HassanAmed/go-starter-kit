package middlewares

import (
	"log"
	"net/http"

	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(bearerSchema):]
		if token, err := service.NewJWTService().ValidateToken(tokenString); err != nil || !token.Valid {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
