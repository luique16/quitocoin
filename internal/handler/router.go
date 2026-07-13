package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/usecase"
)

func NewRouter(
	create *usecase.CreateUserUseCase,
	get *usecase.GetUserUseCase,
	list *usecase.ListUsersUseCase,
	update *usecase.UpdateUserUseCase,
	del *usecase.DeleteUserUseCase,
) *gin.Engine {
	r := gin.Default()

	r.POST("/users", HandleCreateUser(create))
	r.GET("/users/:id", HandleGetUser(get))
	r.GET("/users", HandleListUsers(list))
	r.PUT("/users/:id", HandleUpdateUser(update))
	r.DELETE("/users/:id", HandleDeleteUser(del))

	return r
}
