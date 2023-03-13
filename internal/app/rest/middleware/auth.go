package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/claim"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim := claim.NewClaim()
		if err := claim.Bind(c); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err := claim.Valid(); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("auth", claim)
		c.Next()
	}
}
