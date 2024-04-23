package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type UsersInfo interface {
	UserExists(ctx context.Context, email string) bool
	CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserBySession(ctx context.Context, sessionID string) (*models.User, error)
	GetLastID(ctx context.Context) uint

	EditUser(id uint, email, passwordHash string) (*models.User, error)

	SessionExists(sessionID string) bool
	AddSession(email uint) string
	RemoveSession(ctx context.Context, sessionID string) error
}
