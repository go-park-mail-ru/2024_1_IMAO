package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (active *ActiveUsers) Login(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	password := request.FormValue("password")
	user, ok := active.getUserByEmail(email)

	if ok != nil {
		http.Error(writer, `Wrong username or password`, 404)
		return
	}

	if !CheckPassword(password, user.PasswordHash) {
		http.Error(writer, `Wrong username or password`, 404)
		return
	}

	sessionID := active.addSession(email)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(writer, cookie)
	writer.Write([]byte("You have been authorized with session ID: "))
	writer.Write([]byte(sessionID))
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
	email := request.FormValue("email")
	password := request.FormValue("password")
	passwordRepeat := request.FormValue("password_repeat")

	if password != passwordRepeat {
		http.Error(writer, "Passwords do not match", 401)

		return
	}

	_, err := active.createUser(email, HashPassword(password))

	if err != nil {
		http.Error(writer, "User already exists", 401)
	}

	sessionID := active.addSession(email)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(writer, cookie)
	writer.Write([]byte("You have been authorized with session ID: "))
	writer.Write([]byte(sessionID))
}

func (active *ActiveUsers) Root(writer http.ResponseWriter, request *http.Request) {
	authorized := false
	session, err := request.Cookie("session_id")

	if err == nil && session != nil {
		authorized = active.sessionExists(session.Value)
	}

	if authorized {
		user, err := active.GetUserBySession(session.Value)

		if err != nil {
			fmt.Println("No such session")
			return
		}

		writer.Write([]byte("You authorized as user "))
		writer.Write([]byte(user.Email))
	} else {
		writer.Write([]byte("You have no authorization"))
	}
}
