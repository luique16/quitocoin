package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

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
