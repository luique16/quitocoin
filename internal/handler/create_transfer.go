package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleCreateTransfer(uc *usecase.CreateTransferUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input usecase.CreateTransferInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		claims := middleware.GetClaims(c)

		result, err := uc.Execute(c.Request.Context(), claims.PublicID, input)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
