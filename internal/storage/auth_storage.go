package storage

import (
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2024_1_IMAO/pkg"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
)

const (
	sessionIDLen = 32
)

type UnauthorizedUser struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
}

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UsersList struct {
	// Ключ - id сессии, значение - id пользователя
	Sessions   map[string]uint
	Users      map[uint]*User
	UsersCount uint
	mu         sync.RWMutex
}

type UsersInfo interface {
	UserExists(email string) bool
	CreateUser(email, passwordHash string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserBySession(sessionID string) (*User, error)
	getLastID() uint

	EditUser(email, passwordHash string) (*User, error)

	SessionExists(sessionID string) bool
	AddSession(email string) string
	RemoveSession(sessionID string) error

	UsersList
}

func (active *UsersList) UserExists(email string) bool {
	_, err := active.getIDByEmail(email)

	return err == nil
}

func (active *UsersList) getLastID() uint {
	active.UsersCount++

	return active.UsersCount
}

func (active *UsersList) getIDByEmail(email string) (uint, error) {
	active.mu.Lock()
	defer active.mu.Unlock()

	for _, val := range active.Users {
		if val.Email == email {
			return val.ID, nil
		}
	}

	return 0, errUserNotExists
}

func (active *UsersList) CreateUser(email, passwordHash string) (*User, error) {
	if active.UserExists(email) {
		return nil, errUserExists
	}

	active.mu.Lock()
	defer active.mu.Unlock()

	id := active.getLastID()

	active.Users[id] = &User{
		ID:           id,
		PasswordHash: passwordHash,
		Email:        email,
	}

	return active.Users[id], nil
}

func (active *UsersList) EditUser(id uint, email, passwordHash string) (*User, error) {
	active.mu.Lock()
	defer active.mu.Unlock()

	usr, ok := active.Users[id]

	if !ok {
		return nil, errUserNotExists
	}

	usr.PasswordHash = passwordHash
	usr.Email = email

	return usr, nil
}

func (active *UsersList) GetUserByEmail(email string) (*User, error) {
	active.mu.Lock()
	defer active.mu.Unlock()

	usr, err := active.getIDByEmail(email)

	if err != nil {

		return nil, errUserNotExists
	}

	return active.Users[usr], nil
}

func (active *UsersList) GetUserByID(userID uint) (*User, error) {
	active.mu.Lock()
	defer active.mu.Unlock()

	usr, ok := active.Users[userID]

	if ok {
		return usr, nil
	}

	return nil, errUserNotExists
}

func (active *UsersList) GetUserBySession(sessionID string) (*User, error) {
	active.mu.Lock()
	defer active.mu.Unlock()

	id := active.Sessions[sessionID]

	for _, val := range active.Users {
		if val.ID == id {
			return val, nil
		}
	}

	return nil, errUserNotExists
}

func (active *UsersList) SessionExists(sessionID string) bool {
	_, exists := active.Sessions[sessionID]

	return exists
}

func (active *UsersList) AddSession(id uint) string {
	sessionID := pkg.RandString(sessionIDLen)

	active.mu.Lock()
	defer active.mu.Unlock()

	user := active.Users[id]

	active.Sessions[sessionID] = user.ID

	return sessionID
}

func (active *UsersList) RemoveSession(sessionID string) error {
	if !active.SessionExists(sessionID) {
		return errSessionNotExists
	}

	active.mu.Lock()
	defer active.mu.Unlock()

	delete(active.Sessions, sessionID)

	return nil
}

func NewActiveUser() *UsersList {
	return &UsersList{
		Sessions: make(map[string]uint, 1),
		Users: map[uint]*User{
			1: {
				ID:           1,
				Email:        "example@mail.ru",
				PasswordHash: pkg.HashPassword("123456"),
			},
		},
		UsersCount: 1,
		mu:         sync.RWMutex{},
	}
}
