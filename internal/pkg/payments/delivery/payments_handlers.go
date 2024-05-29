//nolint:noctx
package delivery

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	paymentsusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/payments/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	yookassaURL = "https://api.yookassa.ru/v3/payments"
)

type PaymentsHandler struct {
	storage    paymentsusecases.PaymentsStorageInterface
	authClient authproto.AuthClient
}

func NewPaymentsHandler(storage paymentsusecases.PaymentsStorageInterface,
	authClient authproto.AuthClient) *PaymentsHandler {
	return &PaymentsHandler{
		storage:    storage,
		authClient: authClient,
	}
}

func (h *PaymentsHandler) GetPaymentForm(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := h.storage
	authClient := h.authClient

	var frontendData models.ReceivedPaymentFormItem

	data, _ := io.ReadAll(request.Body)

	err := frontendData.UnmarshalJSON(data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	ownership := storage.CheckAdvertOwnership(ctx, frontendData.AdvertID, uint(user.ID))

	if !ownership {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	returnURL := os.Getenv("RETURN_URL")
	username := os.Getenv("YUKASSA_USERNAME")
	password := os.Getenv("YUKASSA_PASSWORD")
	idempotencyKey := uuid.New().String()

	var priceAndDescription *models.PriceAndDescription

	priceAndDescription, err = storage.GetPriceAndDescription(ctx, frontendData.AdvertID, frontendData.Rate)

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	paymentData := &models.PaymentInitData{
		Amount: models.PaymentInitAmount{
			Value:    priceAndDescription.Price,
			Currency: "RUB",
		},
		PaymentMethodData: models.PaymentInitPaymentMethodData{
			Type: "bank_card",
		},
		Confirmation: models.PaymentInitConfirmation{
			Type:      "redirect",
			ReturnURL: returnURL + priceAndDescription.URLEnding,
		},
		Description: priceAndDescription.Description,
	}

	client := &http.Client{}

	jsonData, err := json.Marshal(paymentData)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	req, err := http.NewRequest("POST", yookassaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Idempotence-Key", idempotencyKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}
	defer resp.Body.Close()

	var payment models.Payment

	respData, _ := io.ReadAll(resp.Body)

	err = payment.UnmarshalJSON(respData)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	err = storage.CreatePayment(ctx, &payment, idempotencyKey, frontendData.AdvertID, priceAndDescription.Duration)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	confirmationURL := payment.Confirmation.ConfirmationURL

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer,
		responses.NewOkResponse(models.PaymentFormResponse{PaymentFormURL: confirmationURL}))
}
