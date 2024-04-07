package usecases

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type UsersInfo interface {
	UserExists(email string) bool
	CreateUser(email, passwordHash string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserBySession(sessionID string) (*models.User, error)
	GetLastID() uint

	EditUser(id uint, email, passwordHash string) (*models.User, error)

	SessionExists(sessionID string) bool
	AddSession(email uint) string
	RemoveSession(sessionID string) error
}
