package controllers

import (
	"path"

	"github.com/gin-gonic/gin"
	appContext "github.com/lstrgiang/cryptowatch-server/internal/app/rest/context"
)

func RegisterPrivateControllers(
	e *gin.RouterGroup,
	appContext appContext.Context,
	authMiddleware gin.HandlerFunc,
	rootPath string,
) {
	authHandler := NewAuthHandler()
	// POST /users
	e.POST("users", authMiddleware, authHandler.GetUser(appContext))

	orderHandler := NewOrderHandler()
	orderEndpoint := e.Group("/orders", authMiddleware)
	// POST /orders/buy
	orderEndpoint.POST("buy", orderHandler.NewBuyOrder(appContext))
	// POST /orders/sell
	orderEndpoint.POST("sell", orderHandler.NewSellOrder(appContext))
	// GET /orders
	orderEndpoint.GET("/", orderHandler.GetOrders(appContext))
	orderEndpoint.GET(":id", orderHandler.GetOrderByID(appContext))

	priceHandler := NewPriceHandler()
	// GET /price
	e.GET("price", priceHandler.GetPrice(appContext))
}

func RegisterPublicControllers(
	e *gin.RouterGroup,
	appContext appContext.Context,
	rootPath string,
) {
	authHandler := NewAuthHandler()
	e.POST(path.Join(rootPath, "login"), authHandler.Login(appContext))
}
