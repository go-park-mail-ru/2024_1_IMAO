package storage_test

import (
	"errors"
	"reflect"
	"testing"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
)

func TestUserExists(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()

	if !active.UserExists("example@mail.ru") {
		t.Errorf("Expected user to exist")
	}

	if active.UserExists("nonexistent@mail.ru") {
		t.Errorf("Expected user to not exist")
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()

	user, err := active.CreateUser("newuser@mail.ru", "hashedpassword")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.Email != "newuser@mail.ru" {
		t.Errorf("Expected email to be newuser@mail.ru, got %s", user.Email)
	}

	_, err = active.CreateUser("example@mail.ru", "hashedpassword")
	if !reflect.DeepEqual(err, errUserExists) {
		t.Errorf("Expected error %v, got %v", errUserExists, err)
	}
}

func TestGetUserBySession(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()
	sessionID := active.AddSession("example@mail.ru")

	user, err := active.GetUserBySession(sessionID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Email != "example@mail.ru" {
		t.Errorf("Expected email to be example@mail.ru, got %s", user.Email)
	}

	_, err = active.GetUserBySession("nonexistentSession")
	if !reflect.DeepEqual(err, errUserNotExists) {
		t.Errorf("Expected error %v, got %v", errUserNotExists, err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()

	user, err := active.GetUserByEmail("example@mail.ru")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Email != "example@mail.ru" {
		t.Errorf("Expected email to be example@mail.ru, got %s", user.Email)
	}

	_, err = active.GetUserByEmail("nonexistent@mail.ru")
	if !reflect.DeepEqual(err, errUserNotExists) {
		t.Errorf("Expected error %v, got %v", errUserNotExists, err)
	}
}

func TestSessionExists(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()

	sessionID := active.AddSession("example@mail.ru")

	if !active.SessionExists(sessionID) {
		t.Errorf("Expected session to exist")
	}

	if active.SessionExists("nonexistentSession") {
		t.Errorf("Expected session to not exist")
	}
}

func TestAddSession(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()

	sessionID := active.AddSession("example@mail.ru")

	if !active.SessionExists(sessionID) {
		t.Errorf("Expected session to be added")
	}
}

func TestRemoveSession(t *testing.T) {
	t.Parallel()

	active := storage.NewActiveUser()

	sessionID := active.AddSession("example@mail.ru")

	err := active.RemoveSession(sessionID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if active.SessionExists(sessionID) {
		t.Errorf("Expected session to be removed")
	}

	err = active.RemoveSession("nonexistentSession")
	if !reflect.DeepEqual(err, errSessionNotExists) {
		t.Errorf("Expected error %v, got %v", errSessionNotExists, err)
	}
}
