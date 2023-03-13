package claim

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type (
	NewClaimFn func() Claim
	Claim      interface {
		Bind(*gin.Context) error
		Valid() error
		GetJTI() string
		GetUserID() int64
		GetExpireTime() int64
		SignWithSecret(secret string, signingMethod *jwt.SigningMethodHMAC) (string, error)
	}
	claim struct {
		UserID     int64  `json:"user_id"`
		ExpireTime int64  `json:"exp"`
		JTI        string `json:"jti"`
	}
)

const (
	claimKey = "ctx_claim"
	//TODO: move this to global configuration
	defaultExpireTime = 3 * 60 * 24 * time.Hour
)

func Set(c *gin.Context, claim Claim) {
	c.Set(claimKey, claim)
}

func Get(c *gin.Context) Claim {
	val, ok := c.Get(claimKey)
	if val == nil || !ok {
		return nil
	}
	claim, ok := val.(*claim)
	if !ok {
		return nil
	}
	return claim
}

// create new claim with default expire time
// and auto generated jti with given userId
func NewDefaultExpireClaim(userId int64) Claim {
	return &claim{
		UserID:     userId,
		ExpireTime: time.Now().Add(defaultExpireTime).Unix(),
		JTI:        uuid.New().String(),
	}
}

func NewClaim() Claim {
	return &claim{}
}

func (c *claim) Bind(ctx *gin.Context) error {
	authCookie, err := ctx.Copy().Cookie("Authorization")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(authCookie), c)
	return err
}

func (c claim) Valid() error {
	if c.ExpireTime == 0 {
		return errors.New("token expired")
	}
	if c.UserID == 0 || c.JTI == "" {
		return errors.New("Invalid token")
	}
	return nil
}

func (c claim) GetJTI() string {
	return c.JTI
}

func (c claim) GetUserID() int64 {
	return c.UserID
}

func (c claim) GetExpireTime() int64 {
	return c.ExpireTime
}

func (c claim) SignWithSecret(secret string, signingMethod *jwt.SigningMethodHMAC) (string, error) {
	token := jwt.NewWithClaims(signingMethod, c)
	return token.SignedString(secret)
}
