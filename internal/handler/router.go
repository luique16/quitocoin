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
	getMe *usecase.GetUserUseCase,
	updateMe *usecase.UpdateUserUseCase,
	updatePassword *usecase.UpdatePasswordUseCase,
	deleteMe *usecase.DeleteUserUseCase,
	jwtProvider provider.JWTProvider,
) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", HandleRegister(register))
		auth.POST("/login", HandleLogin(login))
	}

	me := r.Group("/me")
	me.Use(middleware.Auth(jwtProvider))
	{
		me.GET("", HandleGetMe(getMe))
		me.PUT("", HandleUpdateMe(updateMe))
		me.PUT("/password", HandleUpdatePassword(updatePassword))
		me.DELETE("", HandleDeleteMe(deleteMe))
	}

	return r
}
