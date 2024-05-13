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

type SessionList struct {
	Sessions map[string]uint
	Mux      sync.RWMutex
}

type CSRFToken struct {
	TokenBody string
}

type Session struct {
	UserID uint32
	Value  string
}

type DBInsertionUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
