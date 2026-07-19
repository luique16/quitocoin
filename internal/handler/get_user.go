package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Get current user
// @Description  Get authenticated user profile and balance
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{id=string,name=string,email=string,public_id=string,balance=number,created_at=string}
// @Failure      401  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /me [get]
func HandleGetMe(uc *usecase.GetUserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.GetClaims(c)

		result, err := uc.Execute(c.Request.Context(), claims.UserID, claims.PublicID)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, toUserResponse(result.User, result.Balance))
	}
}
