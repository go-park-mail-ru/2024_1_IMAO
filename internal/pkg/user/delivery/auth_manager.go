package delivery

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	protobuf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	errWrongPassword = errors.New("passwords do not match")
)

type AuthManager struct {
	protobuf.UnimplementedAuthServer

	UserStorage userusecases.UsersStorageInterface
}

func NewAuthManager(storage userusecases.UsersStorageInterface) *AuthManager {
	return &AuthManager{
		UserStorage: storage,
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
	}
}

func (manager *AuthManager) Login(ctx context.Context, in *protobuf.ExistedUserData) (*protobuf.LoggedUser, error) {
	email := in.GetEmail()
	password := in.GetPassword()
	storage := manager.UserStorage

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

func (manager *AuthManager) Logout(ctx context.Context, in *protobuf.SessionData) (*emptypb.Empty, error) {
	sessionID := in.GetSessionID()
	storage := manager.UserStorage

	err := storage.RemoveSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (manager *AuthManager) Signup(ctx context.Context,
	in *protobuf.NewUserData) (*protobuf.LoggedUser, error) {
	email := in.GetEmail()
	password := in.GetPassword()
	passwordRepeat := in.GetPasswordRepeat()
	storage := manager.UserStorage

	user, err := storage.CreateUser(ctx, email, password, passwordRepeat)
	if err != nil {
		return nil, err
	}
	sessionID := storage.AddSession(user.ID)

	return &protobuf.LoggedUser{
		ID:           uint64(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		SessionID:    sessionID,
		IsAuth:       true,
	}, nil
}

func (manager *AuthManager) GetCurrentUser(ctx context.Context, in *protobuf.SessionData) (*protobuf.AuthUser, error) {
	sessionID := in.GetSessionID()
	storage := manager.UserStorage

	if !storage.SessionExists(sessionID) {
		return &protobuf.AuthUser{}, nil
	}

	user, _ := storage.GetUserBySession(ctx, sessionID)

	return &protobuf.AuthUser{
		ID:           uint64(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsAuth:       true,
	}, nil
}

func (manager *AuthManager) EditEmail(ctx context.Context, in *protobuf.EditEmailRequest) (*protobuf.User, error) {
	email := in.GetEmail()
	sessionID := in.GetSessionID()
	storage := manager.UserStorage

	user, _ := storage.GetUserBySession(ctx, sessionID)
	user, err := storage.EditUserEmail(ctx, user.ID, email)
	if err != nil {
		return nil, err
	}

	return newProtobufUser(user, "", true), nil
}
