package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/claim"
	appContext "github.com/lstrgiang/cryptowatch-server/internal/app/rest/context"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/contract"
	"github.com/lstrgiang/cryptowatch-server/internal/data"
	"github.com/lstrgiang/cryptowatch-server/internal/services"
)

func NewOrderHandler() OrderHandler {
	return OrderHandler{
		NewOrderService: services.NewOrderService,
	}
}

type OrderHandler struct {
	NewOrderService services.NewOrderServiceFn
}

func (h OrderHandler) NewBuyOrder(appContext appContext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := claim.Get(ctx)
		order := &data.Order{
			UserID: auth.GetUserID(),
			Type:   data.BuyOrder,
		}
		orderService := h.NewOrderService(appContext.GetDB())
		createdOrder, err := orderService.NewOrder(ctx, order)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusCreated, contract.NewOrderContract(*createdOrder))
	}
}

func (h OrderHandler) NewSellOrder(appContext appContext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := claim.Get(ctx)
		order := &data.Order{
			UserID: auth.GetUserID(),
			Type:   data.SellOrder,
		}
		orderService := h.NewOrderService(appContext.GetDB())
		createdOrder, err := orderService.NewOrder(ctx, order)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusCreated, contract.NewOrderContract(*createdOrder))
	}
}

func (h OrderHandler) GetOrders(appContext appContext.Context) gin.HandlerFunc {
	// should support pagination
	return func(ctx *gin.Context) {
		auth := claim.Get(ctx)
		orderService := h.NewOrderService(appContext.GetDB())
		orders, err := orderService.ListOrder(ctx, auth.GetUserID())
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, contract.NewOrdersContract(orders))
	}
}
func (h OrderHandler) GetOrderByID(appContext appContext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusNotImplemented)
	}
}
