package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleDeleteMe(uc *usecase.DeleteUserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.GetClaims(c)

		err := uc.Execute(c.Request.Context(), claims.UserID)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
