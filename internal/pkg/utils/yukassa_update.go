package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

func YuKassaUpdates() (*models.PaymentList, error) {
	LoadEnv()
	username := os.Getenv("YUKASSA_USERNAME")
	password := os.Getenv("YUKASSA_PASSWORD")

	client := &http.Client{}

	url := "https://api.yookassa.ru/v3/payments?limit=10&status=waiting_for_capture"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	paymentList := models.PaymentList{}
	err = json.NewDecoder(resp.Body).Decode(&paymentList)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return &paymentList, nil
}
