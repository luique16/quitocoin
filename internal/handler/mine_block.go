package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleMineBlock(uc *usecase.MineBlockUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input usecase.MineBlockInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		claims := middleware.GetClaims(c)

		b, err := uc.Execute(c.Request.Context(), claims.UserID, input)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"hash":          b.Hash,
			"index":         b.Index,
			"previous_hash": b.PreviousHash,
			"nonce":         b.Nonce,
			"miner":         b.Miner,
			"reward":        b.Reward,
			"transactions":  b.Transactions,
		})
	}
}
