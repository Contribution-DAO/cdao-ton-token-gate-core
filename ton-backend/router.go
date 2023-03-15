package main

import (
	"github.com/Contribution-DAO/cdao-ton-token-gate-core/controllers"
	"github.com/Contribution-DAO/cdao-ton-token-gate-core/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) {
	// Service handler
	serviceHandler := services.NewServiceHandler(db)

	// Controller handler
	controllerHandler := controllers.NewControllerHandler(serviceHandler, db)

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	router.Use(cors.New(config))

	// internal := router.Group("/internal")

	// Wallet binding
	router.POST("/ton-proof/generatePayload", controllerHandler.GenerateWalletSignPayload)
	router.POST("/ton-proof/checkProof", controllerHandler.ValidateWalletSignature)
	router.GET("/dapp/getAccountInfo", controllerHandler.GetTonAddressInfo)
	router.GET("/wallet/me", controllers.AuthorizationMiddleware, controllerHandler.MeHandler)

	// Social auth binding (Currently only twitter)
	router.GET("/social-auth/:provider/login", controllerHandler.HandleTwitterLogin)
	router.GET("/social-auth/:provider/callback", controllerHandler.HandleTwitterCallback)

	// Telegram binding
	router.GET("/telegram/callback", controllerHandler.HandleTelegramCallback)

	// Telegram group fetch
	router.GET("/telegram/groups", controllerHandler.ListTelegramGroups)
	router.GET("/telegram/groups/:id", controllerHandler.GetTelegramGroup)
	router.GET("/internal/:address/telegram/groups/:id", controllerHandler.GetTelegramGroupRoot)
	router.POST("/telegram/groups/link", controllerHandler.CreateTelegramGroup)

	// Twitter follow
	router.GET("/twitter/follow/:groupId", controllerHandler.VerifyTwitterFollow)

	// SBT
	router.POST("/sbt/link", controllerHandler.LinkSbt)
	router.POST("/sbt/scan/:approvalId", controllerHandler.ScanSbt)
	router.GET("/sbt/metadata/:approvalId", controllerHandler.GetNftMetadata)

	println("The graph backend server listen to port 8040")
	router.Run(":8040")
}
