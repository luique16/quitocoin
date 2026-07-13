package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/provider"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func Auth(jwtProvider provider.JWTProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := jwtProvider.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(string(ClaimsKey), claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *provider.JWTClaims {
	claims, _ := c.Get(string(ClaimsKey))
	return claims.(*provider.JWTClaims)
}
