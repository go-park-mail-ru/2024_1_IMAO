//nolint:forcetypeassert
package responses

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
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

func sendResponse(writer http.ResponseWriter, serverResponse []byte) {
	_, err := writer.Write(serverResponse)
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}
}

func SendOkResponse(writer http.ResponseWriter, response *models.OkResponse) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(StatusOk)

	serverResponse, err := response.MarshalJSON()
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	sendResponse(writer, serverResponse)
}

func SendErrResponse(request *http.Request, writer http.ResponseWriter, response *models.ErrResponse) {
	code := request.Context().Value("code").(*int)
	*code = response.Code

	serverResponse, err := response.MarshalJSON()
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	sendResponse(writer, serverResponse)
}

func NewErrResponse(code int, status string) *models.ErrResponse {
	return &models.ErrResponse{
		Code:   code,
		Status: status,
	}
}

func NewOkResponse(items any) *models.OkResponse {
	return &models.OkResponse{
		Code:  StatusOk,
		Items: items,
	}
}
