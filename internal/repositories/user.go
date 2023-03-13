package repositories

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lstrgiang/cryptowatch-server/internal/data"
)

func NewUserRepository() UserRepository {
	return userRepository{}
}

type (
	UserRepository interface {
		CreateIfNotExist(ctx context.Context, db *sqlx.DB, user *data.User) (*data.User, error)
		GetByEmail(ctx context.Context, db *sqlx.DB, email string, columns ...string) (*data.User, error)
	}
	userRepository struct {
	}
)

// all columns
var Columns = struct {
	ID        string
	Email     string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}{
	ID:        "id",
	Email:     "email",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// all user columns
var AllColumns = []string{
	Columns.ID,
	Columns.Email,
	Columns.CreatedAt,
	Columns.UpdatedAt,
	Columns.DeletedAt,
}

var allColumnsString = strings.Join(AllColumns, ",")

func New(db *sqlx.DB) UserRepository {
	return userRepository{}
}

// Create new user with email if no row with given email exist
func (u userRepository) CreateIfNotExist(ctx context.Context, db *sqlx.DB, user *data.User) (*data.User, error) {
	var createdUser *data.User
	if err := db.Select(createdUser, `
		INSERT INTO users (email)
		VALUES (:email)
		ON CONFLICT (email) DO NOTHING
		RETURNING *
	`, user); err != nil {
		return nil, err
	}
	return createdUser, nil
}

// Get user by email
func (u userRepository) GetByEmail(ctx context.Context, db *sqlx.DB, email string, columns ...string) (*data.User, error) {
	colList := allColumnsString
	if len(columns) > 0 {
		colList = strings.Join(columns, ",")
	}
	var user *data.User
	if err := db.Select(user, `
		SELECT (?) FROM users
		WHERE email = ?
	`, colList, email); err != nil {
		return nil, err
	}
	return user, nil
}
