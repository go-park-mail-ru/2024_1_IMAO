package models

import "sync"

type City struct {
	ID          uint   `json:"id"`
	CityName    string `json:"name"`
	Translation string `json:"translation"`
}

type CityList struct {
	CityItems []*City
	Mux       sync.RWMutex
}

type Coordinates struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Language string  `json:"language"`
}

type GeocoderResponse struct {
	Suggestions []SuggestionInfo `json:"suggestions"`
}

type SuggestionInfo struct {
	Value             string         `json:"value"`
	UnrestrictedValue string         `json:"unrestricted_value"`
	Data              SuggestionData `json:"data"`
}

type SuggestionData struct {
	PostalCode      string `json:"postal_code"`
	Country         string `json:"country"`
	CountryISOCode  string `json:"country_iso_code"`
	FederalDistrict string `json:"federal_district"`
	City            string `json:"city"`
}

type CityResponse struct {
	CityName string `json:"cityName"`
}
