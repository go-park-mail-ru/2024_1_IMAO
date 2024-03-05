package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	StatusOk      = 200
	StatusCreated = 201

	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusNotFound     = 404
	StatusNotAllowed   = 405

	StatusInternalServerError = 500
)

const (
	ErrUserAlreadyExists  = "User with same email already exists"
	ErrWrongCredentials   = "Wrong credentials" //nolint:gosec
	ErrUnauthorized       = "User not authorized"
	ErrAuthorized         = "User already authorized"
	ErrWrongEmailFormat   = "Wrong email format"
	ErrDifferentPasswords = "Passwords do not match"

	ErrAdvertNotExist = "Advert does not exist"
	ErrTooManyAdverts = "Too many adverts specified"

	ErrInternalServer = "Server error"
	ErrBadRequest     = "Bad request"
	ErrNotAllowed     = "Method not allowed"
)

func sendResponse(writer http.ResponseWriter, response any) {
	serverResponse, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	_, err = writer.Write(serverResponse)
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}
}

func SendOkResponse(writer http.ResponseWriter, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(StatusOk)
	sendResponse(writer, response)
}

func SendErrResponse(writer http.ResponseWriter, response any) {
	writer.Header().Set("Content-Type", "application/json")
	sendResponse(writer, response)
}
