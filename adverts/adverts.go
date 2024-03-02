package adverts

import (
	"encoding/json"
	"net/http"
)

func (ads *AdvertsStorage) Root(writer http.ResponseWriter, request *http.Request) {
	adsList, err := ads.GetSeveralAdverts(50)

	if err != nil {
		http.Error(writer, `Wrong numbers of ads`, 404)

		return
	}

	serverResponse := response{
		Adverts: adsList,
	}

	data, _ := json.Marshal(serverResponse)
	writer.Write(data)
}
