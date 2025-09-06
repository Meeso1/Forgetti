package dto

type ErrorResponse struct {
	Message string `json:"message"`
	ErrorCode string `json:"error_code"`
	Data map[string]string `json:"data"`
}