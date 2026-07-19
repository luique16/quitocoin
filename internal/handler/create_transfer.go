package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Create transfer
// @Description  Create a new coin transfer transaction
// @Tags         Transfer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body usecase.CreateTransferInput true "Transfer data"
// @Success      201  {object}  object{from=string,to=string,amount=number}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      422  {object}  object{error=string}
// @Router       /transfer [post]
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
