package models

type ErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type OkResponse struct {
	Code  int `json:"code"`
	Items any `json:"items"`
}
