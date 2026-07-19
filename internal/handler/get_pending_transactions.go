package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Get pending transactions
// @Description  Get all pending transactions from the mempool
// @Tags         Transfer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{transactions=array}
// @Failure      401  {object}  object{error=string}
// @Router       /transfer/pending [get]
func HandleGetPendingTransactions(uc *usecase.GetPendingTransactionsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := uc.Execute(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
