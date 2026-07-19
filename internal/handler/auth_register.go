package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      Register
// @Description  Register a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body user.CreateUserInput true "User registration"
// @Success      201  {object}  object{token=string}
// @Failure      400  {object}  object{error=string}
// @Failure      409  {object}  object{error=string}
// @Router       /auth/register [post]
func HandleRegister(uc *usecase.RegisterUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input user.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		result, err := uc.Execute(c.Request.Context(), input)
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"token": result.Token,
		})
	}
}
