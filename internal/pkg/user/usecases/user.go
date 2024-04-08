package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type UsersInfo interface {
	UserExists(ctx context.Context, email string) bool
	CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserBySession(sessionID string) (*models.User, error)
	GetLastID() uint

	EditUser(id uint, email, passwordHash string) (*models.User, error)

	SessionExists(sessionID string) bool
	AddSession(email uint) string
	RemoveSession(sessionID string) error
}
