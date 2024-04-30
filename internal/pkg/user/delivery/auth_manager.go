package delivery

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	protobuf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

var (
	errWrongPassword = errors.New("passwords do not match")
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

func newProtobufUser(user *models.User, avatar string, isAuth bool) *protobuf.User {
	if user == nil {
		user = &models.User{}
	}

	return &protobuf.User{
		ID:           uint64(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Avatar:       avatar,
		IsAuth:       isAuth,
	}
}

func (manager *AuthManager) Login(ctx context.Context, in *protobuf.ExistedUserData) (*protobuf.LoggedUser, error) {
	email := in.GetEmail()
	password := in.GetPassword()
	storage := manager.UserStorage.storage

	user, err := storage.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errWrongPassword
	}

	sessionID := storage.AddSession(user.ID)
	return &protobuf.LoggedUser{
		ID:           uint64(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsAuth:       true,
		SessionID:    sessionID,
	}, nil
}

func (manager *AuthManager) Logout(ctx context.Context, in *protobuf.SessionData) error {
	sessionID := in.GetSessionID()
	storage := manager.UserStorage.storage

	err := storage.RemoveSession(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (manager *AuthManager) CheckAuth(ctx context.Context, in *protobuf.SessionData) *protobuf.User {
	sessionID := in.GetSessionID()
	storage := manager.UserStorage.storage
	profileStorage := manager.UserStorage.profileStorage

	if !storage.SessionExists(sessionID) {
		return newProtobufUser(nil, "", false)
	}

	user, _ := storage.GetUserBySession(ctx, sessionID)
	profile, _ := profileStorage.GetProfileByUserID(ctx, user.ID)

	return newProtobufUser(user, profile.AvatarIMG, true)
}

func (manager *AuthManager) EditEmail(ctx context.Context, in *protobuf.EditEmailRequest) (*protobuf.User, error) {
	email := in.GetEmail()
	id := in.GetID()
	storage := manager.UserStorage.storage

	user, err := storage.EditUserEmail(ctx, uint(id), email)
	if err != nil {
		return nil, err
	}

	return newProtobufUser(user, "", true), nil
}
