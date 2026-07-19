package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleGetNextBlock(uc *usecase.GetNextBlockDataUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.GetClaims(c)

		result, err := uc.Execute(c.Request.Context(), claims.PublicID)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
