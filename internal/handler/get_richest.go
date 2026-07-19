package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Get richest users
// @Description  Get ranking of users by balance
// @Tags         Ranking
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{richest=array}
// @Failure      401  {object}  object{error=string}
// @Router       /ranking [get]
func HandleGetRichest(uc *usecase.GetRichestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := uc.Execute(c.Request.Context())
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
