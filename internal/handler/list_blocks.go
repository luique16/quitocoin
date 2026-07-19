package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

// @Summary      List blocks
// @Description  List all blocks with pagination
// @Tags         Blockchain
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query int false "Results per page (default 20)"
// @Param        offset query int false "Page offset (default 0)"
// @Success      200    {object}  object{blocks=array,total_count=int}
// @Failure      401    {object}  object{error=string}
// @Router       /blockchain/blocks [get]
func HandleListBlocks(uc *usecase.ListBlocksUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 20
		if l := c.Query("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}

		offset := 0
		if o := c.Query("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil && v >= 0 {
				offset = v
			}
		}

		result, err := uc.Execute(c.Request.Context(), usecase.ListBlocksInput{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
