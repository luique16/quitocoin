package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/usecase"
)

func HandleGetUserBlocks(uc *usecase.GetUserBlocksUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		publicID := c.Param("public_id")
		if publicID == "" {
			claims := middleware.GetClaims(c)
			publicID = claims.PublicID
		}
		role := c.DefaultQuery("role", "")
		if role == "any" || role == "all" {
			role = ""
		}
		limit := 100

		if l := c.Query("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}

		result, err := uc.Execute(c.Request.Context(), usecase.GetUserBlocksInput{
			PublicID: publicID,
			Role:     role,
			Limit:    limit,
		})
		if err != nil {
			code := mapError(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
