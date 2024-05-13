package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	// set header
	header := http.Header{}
	header.Set("Cookie", "session_id=iWyVQtoKdnnuYXnIUalFYpVAjRsKbiUy")
	header.Set("Content-Type", "multipart/form-data; boundary=vegetaboundary")
	// get image
	file, err := os.Open("body.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.SetBoundary("vegetaboundary")
	writer.WriteField("category", "handmade")
	writer.WriteField("condition", "2")
	writer.WriteField("title", "Дз 3 по базам данных")
	writer.WriteField("description", "Данное объявление создано в рамках Дз 3 по базам данных")
	writer.WriteField("price", "777")
	writer.WriteField("phone", "7 777 777 77 77")
	writer.WriteField("userId", "7")
	writer.WriteField("city", "Москва")
	writer.WriteField("CSRFToken", "4af181c4e490836b73a663f274b749ce4ddd7a99906543b0d6d6a45867c1375c:1715704283")
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		panic(err)
	}
	io.Copy(part, file)
	writer.Close()

	rate := vegeta.Rate{Freq: 125, Per: time.Second}
	duration := 3000 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://www.vol-4-ok.ru:8080/api/adverts/create",
		//URL:    "http://localhost:8080/api/adverts/create",
		Header: header,
		Body:   body.Bytes(),
	})

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Println(metrics)

	reporter := vegeta.NewTextReporter(&metrics)
	reporter(os.Stdout)

}
