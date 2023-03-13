package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type (
	NewLoginRequestFn func() LoginRequest
	LoginRequest      interface {
		Bind(*gin.Context) error
		Validate() error
		GetToken() string
	}
)

type loginRequest struct {
	Token string `validate:"required"`
}

func NewLoginRequest() LoginRequest {
	return &loginRequest{}
}

func (r *loginRequest) Bind(c *gin.Context) error {
	return c.Copy().BindJSON(r)
}

func (r loginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r loginRequest) GetToken() string {
	return r.Token
}
