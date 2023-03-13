package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	appContext "github.com/lstrgiang/cryptowatch-server/internal/app/rest/context"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/contract"
)

type (
	PriceHandler struct{}
)

func NewPriceHandler() PriceHandler {
	return PriceHandler{}
}

func (h PriceHandler) GetPrice(appContext appContext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cache := appContext.GetCache()
		latestPrice, err := cache.Get()
		if err != nil {
			log.Printf("Cannot get price wiht error %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, contract.Price{
			Price: latestPrice,
		})
	}
}
