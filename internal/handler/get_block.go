package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Get block by index
// @Description  Get a specific block by its index
// @Tags         Blockchain
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        index path int true "Block index"
// @Success      200    {object}  object{block=object{index=int,hash=string,previous_hash=string,nonce=int,miner=string,reward=number,created_at=string,tx_count=int,transactions=array}}
// @Failure      400   {object}  object{error=string}
// @Failure      401   {object}  object{error=string}
// @Failure      404   {object}  object{error=string}
// @Router       /blockchain/blocks/{index} [get]
func HandleGetBlock(uc *usecase.GetBlockUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid block index"})
			return
		}

		result, err := uc.Execute(c.Request.Context(), usecase.GetBlockInput{Index: index})
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
