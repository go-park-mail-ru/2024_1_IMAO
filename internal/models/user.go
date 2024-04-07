package models

import (
	"sync"
)

const (
	SessionIDLen = 32
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
	Mux        sync.RWMutex
}
