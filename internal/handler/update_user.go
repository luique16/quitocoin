package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Update current user
// @Description  Update authenticated user profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body user.UpdateUserInput true "Update user data"
// @Success      200  {object}  object{id=string,name=string,email=string,public_id=string,balance=number,created_at=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /me [put]
func HandleUpdateMe(uc *usecase.UpdateUserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.GetClaims(c)

		var input user.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		result, err := uc.Execute(c.Request.Context(), claims.UserID, input)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, toUserResponse(result, 0))
	}
}
