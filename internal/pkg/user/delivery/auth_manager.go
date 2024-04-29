package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	protobuf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	"log"
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

func (manager *AuthManager) UserExists(ctx context.Context, in *protobuf.EmailOnlyRequest) *protobuf.ExistsResponse {
	email := in.GetEmail()
	storage := manager.UserStorage.storage

	exists := storage.UserExists(ctx, email)

	return &protobuf.ExistsResponse{
		Exists: exists,
	}
}

func (manager *AuthManager) CreateUser(ctx context.Context, in *protobuf.UnauthorizedUser) (*protobuf.User, error) {
	email := in.GetEmail()
	passwordHash := in.GetPasswordHash()
	storage := manager.UserStorage.storage

	user, err := storage.CreateUser(ctx, email, passwordHash)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newProtobufUser(user), nil
}

func (manager *AuthManager) GetUserByEmail(ctx context.Context, in *protobuf.EmailOnlyRequest) (*protobuf.User, error) {
	storage := manager.UserStorage.storage
	email := in.GetEmail()

	user, err := storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newProtobufUser(user), nil
}

func (manager *AuthManager) GetUserBySession(ctx context.Context, in *protobuf.SessionData) (*protobuf.User, error) {
	storage := manager.UserStorage.storage
	sessionID := in.GetSessionID()

	user, err := storage.GetUserBySession(ctx, sessionID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newProtobufUser(user), nil
}

func (manager *AuthManager) EditUserEmail(ctx context.Context, in *protobuf.EditEmailRequest) (*protobuf.User, error) {
	storage := manager.UserStorage.storage
	email := in.GetEmail()
	id := in.GetID()

	user, err := storage.EditUserEmail(ctx, uint(id), email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newProtobufUser(user), nil
}

func (manager *AuthManager) SessionExists(ctx context.Context, in *protobuf.SessionData) *protobuf.ExistsResponse {
	storage := manager.UserStorage.storage
	sessionID := in.GetSessionID()

	exists := storage.SessionExists(sessionID)

	return &protobuf.ExistsResponse{
		Exists: exists,
	}
}

func (manager *AuthManager) AddSession(ctx context.Context, in *protobuf.IDOnlyRequest) *protobuf.SessionData {
	storage := manager.UserStorage.storage
	id := in.GetID()

	sessionID := storage.AddSession(uint(id))

	return &protobuf.SessionData{
		SessionID: sessionID,
	}
}

func (manager *AuthManager) RemoveSession(ctx context.Context, in *protobuf.SessionData) error {
	storage := manager.UserStorage.storage
	sessionID := in.GetSessionID()

	err := storage.RemoveSession(ctx, sessionID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
