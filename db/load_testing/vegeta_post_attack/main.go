//nolint:errcheck
package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"net/http"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	duration = 3 * time.Second // change the duration here
	freq     = 1
)

func customTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"

		tgt.URL = "http://localhost:8080/api/adverts/create" // your url here

		var b bytes.Buffer

		w := multipart.NewWriter(&b)

		// Добавляем файлы или поля в multipart writer
		// Здесь мы добавляем поле CSRFToken
		part, err := w.CreateFormFile("CSRFToken", "562737c1ff567dbd5574c814")
		if err != nil {
			return err
		}

		_, err = part.Write([]byte(`562737c1ff567dbd5574c814`)) // значение salon_id
		if err != nil {
			return err
		}

		// Закрываем writer
		err = w.Close()
		if err != nil {
			return err
		}

		// Присваиваем созданное тело запроса
		tgt.Body = b.Bytes()

		header := http.Header{}
		header.Add("Accept", "application/json")
		header.Add("Content-Type", "multipart/form-data; boundary=--------------------------583322175473986533287422")
		header.Add("Cookie", "session_id=ySbsnTRmsAsAxlFRnCSPDiVQcnzqGDYt")
		tgt.Header = header

		return nil
	}
}

func main() {
	rate := vegeta.Rate{Freq: freq, Per: time.Second} // change the rate here

	targeter := customTargeter()
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Whatever name") {
		metrics.Add(res)
	}

	metrics.Close()

	fmt.Println(metrics)

	reporter := vegeta.NewTextReporter(&metrics)
	reporter(os.Stdout)
}
