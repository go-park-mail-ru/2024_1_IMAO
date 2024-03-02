package auth

import (
	"errors"
	"math/rand"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
)

type unauthorizedUser struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

type response struct {
	User      User   `json:"user"`
	SessionID string `json:"session_id"`
}

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type ActiveUsers struct {
	//Ключ - id сессии, значение - id пользователя
	Sessions   map[string]uint
	Users      map[string]*User
	UsersCount uint
}

type UsersInfo interface {
	userExists(email string) bool
	createUser(email, password string) (*User, error)
	getUserByEmail(email string) (*User, error)
	getUserBySession(sessionID string) (*User, error)
	getLastID() uint

	sessionExists(sessionID string) bool
	addSession(email string) string
	removeSession(sessionID string)
}

func randString(length int) string {
	var result string
	for i := 0; i < length; i++ {
		result += string(letters[rand.Intn(len(letters))])
	}

	return result
}

func (active *ActiveUsers) userExists(email string) bool {
	_, exists := active.Users[email]

	return exists
}

func (active *ActiveUsers) getLastID() uint {
	active.UsersCount++

	return active.UsersCount
}

func (active *ActiveUsers) createUser(email, passwordHash string) (*User, error) {
	if active.userExists(email) {
		return nil, errUserExists
	}

	active.Users[email] = &User{
		ID:           active.getLastID(),
		PasswordHash: passwordHash,
		Email:        email,
	}

	return active.Users[email], nil
}

func (active *ActiveUsers) getUserByEmail(email string) (*User, error) {
	if !active.userExists(email) {
		return nil, errUserNotExists
	}

	return active.Users[email], nil
}

func (active *ActiveUsers) GetUserBySession(sessionID string) (*User, error) {
	id := active.Sessions[sessionID]

	for _, val := range active.Users {
		if val.ID == id {
			return val, nil
		}
	}

	return nil, errUserNotExists
}

func (active *ActiveUsers) sessionExists(sessionID string) bool {
	_, exists := active.Sessions[sessionID]

	return exists
}

func (active *ActiveUsers) addSession(email string) string {
	sessionID := randString(32)
	user := active.Users[email]

	active.Sessions[sessionID] = user.ID

	return sessionID
}

func (active *ActiveUsers) removeSession(sessionID string) error {
	if !active.sessionExists(sessionID) {
		return errSessionNotExists
	}

	delete(active.Sessions, sessionID)

	return nil
}

func NewActiveUser() *ActiveUsers {
	return &ActiveUsers{
		Sessions: make(map[string]uint, 1),
		Users: map[string]*User{
			"example@mail.ru": {
				ID:           1,
				Email:        "example@mail.ru",
				PasswordHash: HashPassword("123456"),
			},
		},
		UsersCount: 1,
	}
}
