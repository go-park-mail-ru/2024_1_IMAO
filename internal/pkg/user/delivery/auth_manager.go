package delivery

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	protobuf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

type AuthManager struct {
	protobuf.UnimplementedAuthServer

	UserStorage *AuthHandler
}

func NewAuthManager(storage userusecases.UsersStorageInterface,
	profileStorage profusecases.ProfileStorageInterface) *AuthManager {
	return &AuthManager{
		UserStorage: NewAuthHandler(storage, profileStorage),
	}
}

func newProtobufUser(user *models.User) *protobuf.User {
	return &protobuf.User{
		ID:           uint64(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}
