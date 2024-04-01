package data

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
