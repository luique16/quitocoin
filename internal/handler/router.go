package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/luique16/quitocoin/internal/middleware"
	"github.com/luique16/quitocoin/internal/provider"
	"github.com/luique16/quitocoin/internal/usecase"

	_ "github.com/luique16/quitocoin/docs"
)

// cors permite que o frontend (Next.js, em outra origem) consuma a API.
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
			c.Header("Access-Control-Expose-Headers", "Content-Length")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func NewRouter(
	register *usecase.RegisterUseCase,
	login *usecase.LoginUseCase,
	getMe *usecase.GetUserUseCase,
	updateMe *usecase.UpdateUserUseCase,
	updatePassword *usecase.UpdatePasswordUseCase,
	deleteMe *usecase.DeleteUserUseCase,
	mineBlock *usecase.MineBlockUseCase,
	getNextBlock *usecase.GetNextBlockDataUseCase,
	createTransfer *usecase.CreateTransferUseCase,
	getPendingTransactions *usecase.GetPendingTransactionsUseCase,
	getRichest *usecase.GetRichestUseCase,
	getUserBlocks *usecase.GetUserBlocksUseCase,
	getMyPendingTxs *usecase.GetMyPendingTransactionsUseCase,
	listBlocks *usecase.ListBlocksUseCase,
	getBlock *usecase.GetBlockUseCase,
	jwtProvider provider.JWTProvider,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors())

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

	blockchain := r.Group("/blockchain")
	blockchain.Use(middleware.Auth(jwtProvider))
	{
		blockchain.POST("/mine", HandleMineBlock(mineBlock))
		blockchain.GET("/next-block", HandleGetNextBlock(getNextBlock))
		blockchain.GET("/history", HandleGetUserBlocks(getUserBlocks))
		blockchain.GET("/history/:public_id", HandleGetUserBlocks(getUserBlocks))
		blockchain.GET("/blocks", HandleListBlocks(listBlocks))
		blockchain.GET("/blocks/:index", HandleGetBlock(getBlock))
	}

	transfer := r.Group("/transfer")
	transfer.Use(middleware.Auth(jwtProvider))
	{
		transfer.POST("", HandleCreateTransfer(createTransfer))
		transfer.GET("/pending", HandleGetPendingTransactions(getPendingTransactions))
		transfer.GET("/pending/me", HandleGetMyPendingTransactions(getMyPendingTxs))
	}

	utxo := r.Group("/")
	utxo.Use(middleware.Auth(jwtProvider))
	{
		utxo.GET("/ranking", HandleGetRichest(getRichest))
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
