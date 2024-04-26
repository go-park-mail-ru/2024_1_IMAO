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
	StatusForbidden    = 403
	StatusNotFound     = 404
	StatusNotAllowed   = 405

	StatusInternalServerError = 500
)

const (
	ErrUserAlreadyExists  = "User with same email already exists"
	ErrWrongCredentials   = "Wrong credentials" //nolint:gosec
	ErrUnauthorized       = "User not authorized"
	ErrAuthorized         = "User already authorized"
	ErrDifferentPasswords = "Passwords do not match"
	ErrNotValidData       = "User data is not valid"

	ErrAdvertNotExist = "Advert does not exist"

	ErrInternalServer = "Server error"
	ErrBadRequest     = "Bad request"
	ErrNotAllowed     = "Method not allowed"
	ErrForbidden      = "User have no access to this content"
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

type ErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewErrResponse(code int, status string) *ErrResponse {
	return &ErrResponse{
		Code:   code,
		Status: status,
	}
}
