package services

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lstrgiang/cryptowatch-server/internal/clients/google"
	"github.com/lstrgiang/cryptowatch-server/internal/repositories"
)

type (
	NewAuthServiceFn func(db *sqlx.DB) AuthService
	AuthService      interface {
		Authenticate(ctx context.Context, googleToken string) (int64, error)
	}

	authService struct {
		db       *sqlx.DB
		userRepo repositories.UserRepository
	}
)

func NewAuthService(db *sqlx.DB) AuthService {
	return authService{
		db:       db,
		userRepo: repositories.NewUserRepository(),
	}
}

func (s authService) Authenticate(ctx context.Context, googleToken string) (int64, error) {
	info, err := google.VerifyToken(googleToken, &http.Client{})
	if err != nil {
		return 0, errors.New("invalid token")
	}
	user, err := s.userRepo.GetByEmail(ctx, s.db, info.Email)
	if err == sql.ErrNoRows {
		// user not found, insert and return new id
		//TODO: implement
	}
	return user.ID, nil
}
