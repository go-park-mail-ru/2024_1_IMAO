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
