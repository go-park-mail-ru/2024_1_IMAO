package delivery

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
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

	UserStorage    userusecases.UsersStorageInterface
	ProfileStorage profusecases.ProfileStorageInterface
}

func NewAuthManager(storage userusecases.UsersStorageInterface,
	profileStorage profusecases.ProfileStorageInterface) *AuthManager {
	return &AuthManager{
		UserStorage:    storage,
		ProfileStorage: profileStorage,
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
	profileStorage := manager.ProfileStorage

	user, err := storage.CreateUser(ctx, email, password, passwordRepeat)
	if err != nil {
		return nil, err
	}
	profileStorage.CreateProfile(ctx, user.ID)
	sessionID := storage.AddSession(user.ID)

	return &protobuf.LoggedUser{
		ID:           uint64(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		SessionID:    sessionID,
		IsAuth:       true,
	}, nil
}

func (manager *AuthManager) CheckAuth(ctx context.Context, in *protobuf.SessionData) (*protobuf.User, error) {
	sessionID := in.GetSessionID()
	storage := manager.UserStorage
	profileStorage := manager.ProfileStorage

	if !storage.SessionExists(sessionID) {
		return newProtobufUser(nil, "", false), nil
	}

	user, _ := storage.GetUserBySession(ctx, sessionID)
	profile, _ := profileStorage.GetProfileByUserID(ctx, user.ID)

	return newProtobufUser(user, profile.AvatarIMG, true), nil
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
