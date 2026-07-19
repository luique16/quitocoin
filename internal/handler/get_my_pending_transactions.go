package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Get my pending transactions
// @Description  Get pending transactions for the authenticated user
// @Tags         Transfer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Max transactions (default 100)"
// @Success      200   {object}  object{transactions=array}
// @Failure      401   {object}  object{error=string}
// @Router       /transfer/pending/me [get]
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
