package storage

import (
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
)

type UsersListWrapper struct {
	UsersList *models.UsersList
}

func (active *UsersListWrapper) UserExists(email string) bool {
	_, err := active.getIDByEmail(email)

	return err == nil
}

func (active *UsersListWrapper) GetLastID() uint {
	active.UsersList.UsersCount++

	return active.UsersList.UsersCount
}

func (active *UsersListWrapper) getIDByEmail(email string) (uint, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	for _, val := range active.UsersList.Users {
		if val.Email == email {
			return val.ID, nil
		}
	}

	return 0, errUserNotExists
}

func (active *UsersListWrapper) CreateUser(email, passwordHash string) (*models.User, error) {
	if active.UserExists(email) {
		return nil, errUserExists
	}

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	id := active.GetLastID()

	active.UsersList.Users[id] = &models.User{
		ID:           id,
		PasswordHash: passwordHash,
		Email:        email,
	}

	return active.UsersList.Users[id], nil
}

func (active *UsersListWrapper) EditUser(id uint, email, passwordHash string) (*models.User, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	usr, ok := active.UsersList.Users[id]

	if !ok {
		return nil, errUserNotExists
	}

	usr.PasswordHash = passwordHash
	usr.Email = email

	return usr, nil
}

func (active *UsersListWrapper) GetUserByEmail(email string) (*models.User, error) {
	usr, err := active.getIDByEmail(email)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	if err != nil {

		return nil, errUserNotExists
	}

	return active.UsersList.Users[usr], nil
}

func (active *UsersListWrapper) GetUserByID(userID uint) (*models.User, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	usr, ok := active.UsersList.Users[userID]

	if ok {
		return usr, nil
	}

	return nil, errUserNotExists
}

func (active *UsersListWrapper) GetUserBySession(sessionID string) (*models.User, error) {
	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	id := active.UsersList.Sessions[sessionID]

	for _, val := range active.UsersList.Users {
		if val.ID == id {
			return val, nil
		}
	}

	return nil, errUserNotExists
}

func (active *UsersListWrapper) SessionExists(sessionID string) bool {
	_, exists := active.UsersList.Sessions[sessionID]

	return exists
}

func (active *UsersListWrapper) AddSession(id uint) string {
	sessionID := utils.RandString(models.SessionIDLen)

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	user := active.UsersList.Users[id]

	active.UsersList.Sessions[sessionID] = user.ID

	return sessionID
}

func (active *UsersListWrapper) RemoveSession(sessionID string) error {
	if !active.SessionExists(sessionID) {
		return errSessionNotExists
	}

	active.UsersList.Mux.Lock()
	defer active.UsersList.Mux.Unlock()

	delete(active.UsersList.Sessions, sessionID)

	return nil
}

func NewActiveUser() *UsersListWrapper {
	return &UsersListWrapper{
		UsersList: &models.UsersList{
			Sessions: make(map[string]uint, 1),
			Users: map[uint]*models.User{
				1: {
					ID:           1,
					Email:        "example@mail.ru",
					PasswordHash: utils.HashPassword("123456"),
				},
			},
			UsersCount: 1,
			Mux:        sync.RWMutex{}},
	}
}
