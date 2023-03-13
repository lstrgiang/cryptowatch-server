package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/claim"
	appContext "github.com/lstrgiang/cryptowatch-server/internal/app/rest/context"
	"github.com/lstrgiang/cryptowatch-server/internal/services"
)

type (
	AuthHandler struct {
		NewAuthClaim    claim.NewClaimFn
		NewLoginRequest NewLoginRequestFn
		NewAuthService  services.NewAuthServiceFn
	}
)

func NewAuthHandler() AuthHandler {
	return AuthHandler{
		NewLoginRequest: NewLoginRequest,
		NewAuthClaim:    claim.NewClaim,
		NewAuthService:  services.NewAuthService,
	}
}

func (h AuthHandler) Login(appContext appContext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := h.NewLoginRequest()
		if err := request.Bind(ctx); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := request.Validate(); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		authService := h.NewAuthService(appContext.GetDB())
		userId, err := authService.Authenticate(ctx, request.GetToken())
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claim := claim.NewDefaultExpireClaim(userId)
		jwtString, err := claim.SignWithSecret(
			appContext.GetAuthSecretKey(),
			jwt.SigningMethodHS256,
		)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.SetCookie(
			"Authorization",
			jwtString,
			int(claim.GetExpireTime()),
			"/",
			appContext.GetDomain(),
			true,
			true,
		)
	}

}

func (h AuthHandler) GetUser(appContext appContext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := h.NewLoginRequest()
		if err := request.Bind(ctx); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := request.Validate(); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}
