package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (active *ActiveUsers) Login(writer http.ResponseWriter, request *http.Request) {
	var user unauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	email := user.Email
	password := user.Password
	expectedUser, ok := active.getUserByEmail(email)

	if ok != nil {
		http.Error(writer, `Wrong username or password`, 404)
		return
	}

	if !CheckPassword(password, expectedUser.PasswordHash) {
		http.Error(writer, `Wrong username or password`, 404)
		return
	}

	sessionID := active.addSession(email)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
	}

	serverResponse := response{
		User:      *expectedUser,
		SessionID: sessionID,
		IsAuth:    true,
	}

	userData, _ := json.Marshal(serverResponse)
	writer.Write(userData)

	http.SetCookie(writer, cookie)
	fmt.Println("You have been authorized with session ID: ")
	fmt.Println(sessionID)
}

func (active *ActiveUsers) Logout(writer http.ResponseWriter, request *http.Request) {
	session, err := request.Cookie("session_id")

	if errors.Is(err, http.ErrNoCookie) {
		fmt.Println("No cookie")
		http.Error(writer, `You have no authorization`, 401)
		return
	}

	err = active.removeSession(session.Value)

	if err != nil {
		fmt.Println(err)
		http.Error(writer, `You have no authorization`, 401)

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)
	writer.Write([]byte("You have been logged out"))
}

func (active *ActiveUsers) Signup(writer http.ResponseWriter, request *http.Request) {
	var newUser unauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&newUser)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	if password != passwordRepeat {
		http.Error(writer, "Passwords do not match", 401)

		return
	}

	user, err := active.createUser(email, HashPassword(password))

	if err != nil {
		http.Error(writer, "User already exists", 401)
	}

	sessionID := active.addSession(email)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
	}

	serverResponse := response{
		User:      *user,
		SessionID: sessionID,
		IsAuth:    true,
	}

	userData, _ := json.Marshal(serverResponse)
	writer.Write(userData)

	http.SetCookie(writer, cookie)
	fmt.Println("You have been authorized with session ID: ")
	fmt.Println(sessionID)
}

func (active *ActiveUsers) CheckAuth(writer http.ResponseWriter, request *http.Request) {
	session, err := request.Cookie("session_id")

	if errors.Is(err, http.ErrNoCookie) || !active.sessionExists(session.Value) {
		serverResponse := response{
			IsAuth: false,
		}

		userData, _ := json.Marshal(serverResponse)
		writer.Write(userData)

		return
	}

	user, ok := active.GetUserBySession(session.Value)

	if ok != nil {
		http.Error(writer, `Wrong username or password`, 404)

		return
	}

	serverResponse := response{
		User:      *user,
		SessionID: session.Value,
		IsAuth:    true,
	}

	userData, _ := json.Marshal(serverResponse)
	writer.Write(userData)
}
