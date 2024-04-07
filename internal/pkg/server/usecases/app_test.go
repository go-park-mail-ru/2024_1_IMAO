package usecases_test

import (
	"net/http"
	"testing"
	"time"
)

func TestServerRun_Configuration(t *testing.T) {
	t.Parallel()

	srv := &app.Server{}

	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("Сервер не смог запуститься: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	err := srv.Shutdown()
	if err != nil {
		t.Fatalf("Сервер не смог остановиться: %v", err)
	}
}
