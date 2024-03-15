package storage_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

var (
	errWrongID            = errors.New("wrong adverts ID")
	errWrongAdvertsAmount = errors.New("too many elements specified")
)

func TestGetAdvert(t *testing.T) {
	t.Parallel()

	ads := storage.NewAdvertsList()
	storage.FillAdvertsList(ads)

	advert, err := ads.GetAdvert(0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if advert.ID != 1 {
		t.Errorf("Expected ID 1, got %d", advert.ID)
	}

	_, err = ads.GetAdvert(100)

	if !reflect.DeepEqual(err, errWrongID) {
		t.Errorf("Expected error %v, got %v", errWrongID, err)
	}

}

func TestGetSeveralAdverts(t *testing.T) {
	t.Parallel()

	ads := storage.NewAdvertsList()
	storage.FillAdvertsList(ads)

	adsList, err := ads.GetSeveralAdverts(5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(adsList) != 5 {
		t.Errorf("Expected 5 adverts, got %d", len(adsList))
	}

	_, err = ads.GetSeveralAdverts(100)
	if !reflect.DeepEqual(err, errWrongAdvertsAmount) {
		t.Errorf("Expected error %v, got %v", errWrongAdvertsAmount, err)
	}
}

func TestGetLastID(t *testing.T) {
	t.Parallel()

	ads := storage.NewAdvertsList()

	if ads.GetLastID() != 1 {
		t.Errorf("Expected ID 1, got %d", ads.GetLastID())
	}

	lastID := ads.GetLastID()
	if ads.GetLastID() != lastID+1 {
		t.Errorf("Expected ID %d, got %d", lastID+1, ads.GetLastID())
	}
}
