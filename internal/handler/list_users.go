package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleListUsers(uc *usecase.ListUsersUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := uc.Execute(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		items := make([]userResponse, 0, len(result))

		for _, u := range result {
			items = append(items, toUserResponse(u))
		}

		c.JSON(http.StatusOK, items)
	}
}
