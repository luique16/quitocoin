package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Get next block data
// @Description  Get data needed to mine the next block
// @Tags         Blockchain
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{mined=boolean,data=string}
// @Failure      401  {object}  object{error=string}
// @Router       /blockchain/next-block [get]
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
