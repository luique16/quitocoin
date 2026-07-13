package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleUpdatePassword(uc *usecase.UpdatePasswordUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.GetClaims(c)

		var input user.UpdatePasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if err := uc.Execute(c.Request.Context(), claims.UserID, input); err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password updated"})
	}
}
