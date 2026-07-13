package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleDeleteUser(uc *usecase.DeleteUserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := uc.Execute(c.Request.Context(), id)
		if err != nil {
			code := mapError(err)

			c.JSON(code, gin.H{"error": err.Error()})

			return
		}

		c.Status(http.StatusNoContent)
	}
}
