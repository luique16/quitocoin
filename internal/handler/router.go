package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/provider"
	"github.com/luique16/quitocoin/internal/usecase"
)

func NewRouter(
	register *usecase.RegisterUseCase,
	login *usecase.LoginUseCase,
	get *usecase.GetUserUseCase,
	list *usecase.ListUsersUseCase,
	update *usecase.UpdateUserUseCase,
	del *usecase.DeleteUserUseCase,
	jwtProvider provider.JWTProvider,
) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", HandleRegister(register))
		auth.POST("/login", HandleLogin(login))
	}

	api := r.Group("/users")
	api.Use(middleware.Auth(jwtProvider))
	{
		api.GET("/:id", HandleGetUser(get))
		api.GET("", HandleListUsers(list))
		api.PUT("/:id", HandleUpdateUser(update))
		api.DELETE("/:id", HandleDeleteUser(del))
	}

	return r
}
