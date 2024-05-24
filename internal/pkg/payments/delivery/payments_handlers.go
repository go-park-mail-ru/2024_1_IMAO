package delivery

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	paymentsusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/payments/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	yookassa_url = "https://api.yookassa.ru/v3/payments"
)

type PaymentsHandler struct {
	storage    paymentsusecases.PaymentsStorageInterface
	authClient authproto.AuthClient
}

func NewPaymentsHandler(storage paymentsusecases.PaymentsStorageInterface, authClient authproto.AuthClient) *PaymentsHandler {
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

	err := json.NewDecoder(request.Body).Decode(&frontendData)
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

	ownership := storage.CheckAdvertOwnership(ctx, frontendData.AdvertId, uint(user.ID))

	if !ownership {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	utils.LoadEnv()

	returnURL := os.Getenv("RETURN_URL")
	username := os.Getenv("YUKASSA_USERNAME")
	password := os.Getenv("YUKASSA_PASSWORD")
	idempotencyKey := uuid.New().String()

	var priceAndDescription *models.PriceAndDescription

	priceAndDescription, err = storage.GetPriceAndDescription(ctx, frontendData.AdvertId, frontendData.Rate)

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
			ReturnURL: returnURL + priceAndDescription.UrlEnding,
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

	req, err := http.NewRequest("POST", yookassa_url, bytes.NewBuffer(jsonData))
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

	err = json.NewDecoder(resp.Body).Decode(&payment)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	err = storage.CreatePayment(ctx, &payment, idempotencyKey, frontendData.AdvertId)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	confirmationURL := payment.Confirmation.ConfirmationURL

	// city, err := h.storage.GetPaymentForm(ctx)
	// if err != nil {
	// 	logging.LogHandlerError(logger, err, responses.StatusBadRequest)
	// 	log.Println(err, responses.StatusBadRequest)
	// 	responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
	// 		responses.ErrBadRequest))

	// 	return
	// }

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewPaymentFormOkResponse(confirmationURL))
}
