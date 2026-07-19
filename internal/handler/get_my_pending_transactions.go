package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleGetMyPendingTransactions(uc *usecase.GetMyPendingTransactionsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.GetClaims(c)

		limit := 100
		if l := c.Query("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}

		result := uc.Execute(usecase.GetMyPendingTransactionsInput{
			PublicID: claims.PublicID,
			Limit:    limit,
		})

		c.JSON(http.StatusOK, result)
	}
}
