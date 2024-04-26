package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

//go:generate mockgen -source=user.go -destination=../mocks/user_mocks.go

type UsersStorageInterface interface {
	UserExists(ctx context.Context, email string) bool
	CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserBySession(ctx context.Context, sessionID string) (*models.User, error)
	GetLastID(ctx context.Context) uint

	EditUserEmail(ctx context.Context, id uint, email string) (*models.User, error)

	SessionExists(sessionID string) bool
	AddSession(id uint) string
	RemoveSession(ctx context.Context, sessionID string) error

	MAP_GetUserIDBySession(sessionID string) uint
}
