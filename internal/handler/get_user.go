package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleGetUser(uc *usecase.GetUserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := uc.Execute(c.Request.Context(), id)
		if err != nil {
			code := mapError(err)

			c.JSON(code, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusOK, toUserResponse(result))
	}
}
